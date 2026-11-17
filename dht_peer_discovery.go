package concourse

import (
	"context"
	"log/slog"
	"time"

	evt "codeberg.org/ChrisEineke/go-events/pkg/events"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

type DhtPeerDiscovery struct {
	host              host.Host
	bootstrapNodes    []multiaddr.Multiaddr
	routingDiscovery  *routing.RoutingDiscovery
	runningState      int
	PeerDiscoveredEvt evt.E2[context.Context, *peer.AddrInfo]
}

const (
	DhtDisconnected = iota
	DhtBootstrapPeers
	DhtFindingPeers
)

type PeerDiscovered struct {
	ID peer.ID
}

func NewDhtPeerDiscovery(
	host host.Host,
	routingDiscovery *routing.RoutingDiscovery,
	bootstrapNodes []multiaddr.Multiaddr) *DhtPeerDiscovery {
	return &DhtPeerDiscovery{
		host:             host,
		routingDiscovery: routingDiscovery,
		bootstrapNodes:   bootstrapNodes,
	}
}

func (d *DhtPeerDiscovery) Start(ctx context.Context) error {
	go d.runLoop(ctx)
	return nil
}

func (d *DhtPeerDiscovery) Stop(ctx context.Context) error {
	return nil
}

func (d *DhtPeerDiscovery) runLoop(ctx context.Context) {
	for {
		switch d.runningState {
		case DhtDisconnected:
			slog.Debug("Starting DHT peer discovery...")
			d.runningState = DhtBootstrapPeers
		case DhtBootstrapPeers:
			slog.Debug("Connecting to bootstrap peers...")
			for _, peerAddr := range d.bootstrapNodes {
				peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
				if err != nil {
					slog.ErrorContext(ctx, "skipping bootstrap peer", "peerInfo", peerInfo, "error", err)
					continue
				}
				d.PeerDiscoveredEvt.Fire2(ctx, peerInfo)
			}
			d.runningState = DhtFindingPeers
		case DhtFindingPeers:
			dutil.Advertise(ctx, d.routingDiscovery, "concourse")
			for {
				slog.Debug("Searching for peers...")
				peerCh, err := d.routingDiscovery.FindPeers(ctx, "concourse")
				if err != nil {
					slog.ErrorContext(ctx, "failed to find dht peers: ", "error", err)
					continue
				}
				for addrInfo := range peerCh {
					if addrInfo.ID == d.host.ID() {
						continue
					}
					d.PeerDiscoveredEvt.Fire2(ctx, &addrInfo)
				}
				time.Sleep(10 * time.Second)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
