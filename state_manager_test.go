package concourse

import (
	"testing"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func TestStateVizString(t *testing.T) {
	tests := []struct {
		name     string
		sv       StateViz
		expected string
	}{
		{"Private", StateVizPrivate, "StateVizPrivate"},
		{"Public", StateVizPublic, "StateVizPublic"},
		{"PublicPrivate", StateVizPublicPrivate, "StateVizPublicPrivate"},
		{"Illegal", StateViz(99), "Illegal StateViz value 99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sv.String(); got != tt.expected {
				t.Errorf("StateViz.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewKademliaStateManager(t *testing.T) {
	var nilIpfs *dht.IpfsDHT
	kms := NewKademliaStateManager(nilIpfs)
	if kms == nil {
		t.Fatal("NewKademliaStateManager returned nil")
	}
	if kms.Ipfs != nilIpfs {
		t.Errorf("KademliaStateManager.Ipfs = %v, want %v", kms.Ipfs, nilIpfs)
	}
}
