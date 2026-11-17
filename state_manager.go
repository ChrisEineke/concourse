package concourse

import (
	"context"
	"fmt"
	"strings"

	ds "github.com/ipfs/go-datastore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/rogpeppe/rjson"
)

type StateViz int

const (
	// StateVizPrivate stores the key-value pair only locally on the node.
	StateVizPrivate StateViz = iota
	// StateVizPublic stores the key-value on the node and in the swarm.
	StateVizPublic
)

const (
	DefaultNamespace = "concourse"
	KeySeparator     = "/"
)

func (sa StateViz) String() string {
	switch sa {
	case StateVizPrivate:
		return "StateVizPrivate"
	case StateVizPublic:
		return "StateVizPublic"
	default:
		return fmt.Sprintf("Illegal StateViz value %d", sa)
	}
}

type KademliaStateManager struct {
	Namespace string
	Ipfs      *dht.IpfsDHT
	Datastore ds.Batching
	IsOffline bool
}

func NewKademliaStateManager(ipfs *dht.IpfsDHT, ds ds.Batching) *KademliaStateManager {
	return &KademliaStateManager{
		Namespace: DefaultNamespace,
		Ipfs:      ipfs,
		Datastore: ds,
		IsOffline: true,
	}
}

func (k *KademliaStateManager) OnPeersFound() error {
	if !k.IsOffline {
		return nil
	}
	// transition to online state
	k.IsOffline = false
	fmt.Println("switched to online mode")
	return nil
}

func (k *KademliaStateManager) Get[T any](ctx context.Context, into *T, key string) error {
	options := []routing.Option{}
	if k.IsOffline {
		options = append(options, routing.Offline)
	}

	key = fmt.Sprintf("/%s/%s", k.Namespace, key)

	res, err := k.Ipfs.GetValue(ctx, key, options...)
	if err != nil {
		return err
	}

	return rjson.Unmarshal(res, into)
}

func (k *KademliaStateManager) GetPath[T any](ctx context.Context, into *T, keyParts ...string) error {
	options := []routing.Option{}
	if k.IsOffline {
		options = append(options, routing.Offline)
	}

	key := fmt.Sprintf("/%s/%s", k.Namespace, strings.Join(keyParts, KeySeparator))

	res, err := k.Ipfs.GetValue(ctx, key, options...)
	if err != nil {
		return err
	}

	return rjson.Unmarshal(res, into)
}

func (k *KademliaStateManager) Put[T any](ctx context.Context, viz StateViz, value T, key string) error {
	options := []routing.Option{}
	if k.IsOffline || viz == StateVizPrivate {
		options = append(options, routing.Offline)
	}

	key = fmt.Sprintf("/%s/%s", k.Namespace, key)

	res, err := rjson.Marshal(value)
	if err != nil {
		return err
	}

	return k.Ipfs.PutValue(ctx, key, res, options...)
}

func (k *KademliaStateManager) PutPath[T any](ctx context.Context, viz StateViz, value T, keyParts ...string) error {
	options := []routing.Option{}
	if k.IsOffline || viz == StateVizPrivate {
		options = append(options, routing.Offline)
	}

	key := fmt.Sprintf("/%s/%s", k.Namespace, strings.Join(keyParts, KeySeparator))

	res, err := rjson.Marshal(value)
	if err != nil {
		return err
	}

	return k.Ipfs.PutValue(ctx, key, res, options...)
}
