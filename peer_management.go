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

func (p *PeerManagement) ConnectToPeer(ctx context.Context, addr *peer.AddrInfo) error {
	var peerState *PeerState
	var ok bool
	if peerState, ok = p.peerStates[addr.ID]; !ok {
		if err := p.host.Connect(ctx, *addr); err != nil {
			slog.ErrorContext(ctx, "failed to connect to peer", "addrInfo", addr, "error", err)
			return err
		}
		slog.Info("connected to peer", "peerID", addr.ID)
		peerState = &PeerState{
			id:       addr.ID,
			lastSeen: time.Now(),
		}
		p.peerStates[addr.ID] = peerState
	}
	return nil
}
