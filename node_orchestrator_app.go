package concourse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chigopher/pathlib"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/pnet"
	"github.com/libp2p/go-libp2p/core/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/multiformats/go-multiaddr"
)

type NodeOrchestratorApp struct{}

func (app *NodeOrchestratorApp) Run() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Must instantiate a custom logger as the first thing so that all logging
	// conforms to our requirements.
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	)

	var config struct {
		dbPath           pathlib.Path          `mapstructure:"db.path"`
		swarmKeyFileName string                `mapstructure:"swarm.key.filename"`
		bootstrapNodes   []multiaddr.Multiaddr `mapstructure:"bootstrap.nodes"`
	}
	vi := ViperSetup()

	vi.SetDefault("db.path", "./data")
	if err := vi.UnmarshalKey("db.path", &config.dbPath); err != nil {
		panic(err)
	}
	fmt.Println("db path:", config.dbPath.String())

	vi.SetDefault("swarm.key.filename", "swarm.key")
	if err := vi.UnmarshalKey("swarm.key.filename", &config.swarmKeyFileName); err != nil {
		panic(err)
	}
	fmt.Println("swarm key filename:", config.swarmKeyFileName)

	vi.SetDefault("bootstrap.nodes", []string{})
	if err := vi.UnmarshalKey("bootstrap.nodes", &config.bootstrapNodes); err != nil {
		panic(err)
	}
	fmt.Println("bootstrap nodes:", config.bootstrapNodes)

	// swarm key
	rawSwarmKey, err := os.ReadFile(config.swarmKeyFileName)
	if err != nil {
		panic(err)
	}
	psk, err := pnet.DecodeV1PSK(bytes.NewReader(rawSwarmKey))
	if err != nil {
		panic(err)
	}

	// libp2p initialization
	host, err := libp2p.New(libp2p.PrivateNetwork(psk))
	if err != nil {
		panic(err)
	}
	defer host.Close()

	peerInfo := &peerstore.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(peerInfo)
	fmt.Println("libp2p node address:", addrs[0])

	pubSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	dhtDatastore, err := leveldb.NewDatastore(config.dbPath.String(), nil)
	if err != nil {
		panic(err)
	}
	defer dhtDatastore.Close()

	namespaceValidator := NewNamespaceValidator()

	kademliaDht, err := dht.New(
		host,
		dht.Datastore(dhtDatastore),
		dht.Mode(dht.ModeServer),
		dht.ProtocolPrefix("/"+DefaultNamespace),
		dht.BootstrapPeers(),
		dht.NamespacedValidator(DefaultNamespace, namespaceValidator))
	if err != nil {
		panic(err)
	}
	if err := kademliaDht.Bootstrap(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to start kademlia dht bootstrap", "error", err.Error())
		panic(err)
	}

	// peer discovery

	nodeDiscoveryTopic, err := pubSub.Join("/concourse/node_discovery_dht")
	if err != nil {
		panic(err)
	}
	go streamConsoleTo(ctx, nodeDiscoveryTopic)

	subscription, err := nodeDiscoveryTopic.Subscribe()
	if err != nil {
		panic(err)
	}
	go printMessagesFrom(ctx, subscription)

	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDht)

	peerDiscovery := NewDhtPeerDiscovery(
		host,
		routingDiscovery,
		config.bootstrapNodes)
	if err := peerDiscovery.Start(ctx); err != nil {
		panic(err)
	}

	stateManager := NewKademliaStateManager(kademliaDht, dhtDatastore)

	peerManagement := NewPeerManagement(host)
	peerDiscovery.PeerDiscoveredEvt.On(peerManagement.ConnectToPeer)
	peerDiscovery.PeerDiscoveredEvt.On(func(ctx context.Context, addr *peer.AddrInfo) error {
		stateManager.OnPeersFound()
		return nil
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var value1, value2 int = rand.IntN(1024 * 1024), 0

				if err = stateManager.PutPath(ctx, StateVizPublic, value1, "test"); err != nil && !strings.Contains(err.Error(), "failed to find any peer in table") {
					fmt.Println("failed to put kv", "re-trying in 1 second", err)
					continue
				}
				if err = stateManager.GetPath(ctx, &value2, "test"); err != nil && !strings.Contains(err.Error(), "failed to find any peer in table") {
					_ = routing.ErrNotFound
					fmt.Println("failed to get kv", "re-trying in 1 second", err)
					continue
				}
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down")
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
}
