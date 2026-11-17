package concourse

import (
	"context"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PeerManagement struct {
	peerStates map[peer.ID]*PeerState
	host       host.Host
}

type PeerState struct {
	id       peer.ID
	lastSeen time.Time
}

func NewPeerManagement(host host.Host) *PeerManagement {
	return &PeerManagement{
		peerStates: make(map[peer.ID]*PeerState),
		host:       host,
	}
}

// OnNewPeer reacts to the peer discovery task finding a new cluster peer. The
// task runs periodically. Each cycle the task provides us with *every* peer that
// is has found, so we ensure not to connect more than once to the same peer.
func (p *PeerManagement) OnNewPeer(ctx context.Context, addr *peer.AddrInfo) error {
	var peerState *PeerState
	var ok bool
	if peerState, ok = p.peerStates[addr.ID]; ok {
		peerState.lastSeen = time.Now()
		slog.DebugContext(ctx,
			"updated peer last-seen timestamp",
			"peerID", addr.ID,
			"timestamp", peerState.lastSeen)
		return nil
	}
	if err := p.host.Connect(ctx, *addr); err != nil {
		slog.ErrorContext(ctx,
			"failed to connect to peer",
			"peerID", addr.ID,
			"addr", addr,
			"error", err)
		return err
	}
	peerState = &PeerState{
		id:       addr.ID,
		lastSeen: time.Now(),
	}
	p.peerStates[addr.ID] = peerState
	slog.InfoContext(ctx,
		"succesfully connected to peer",
		"peerID", addr.ID)
	return nil
}
