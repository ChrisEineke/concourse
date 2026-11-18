package documents

import "github.com/moby/moby/api/types/network"

type Port = uint16
type Protocol = network.IPProtocol

type Ports_1_0 struct {
	HostPort      Port     `json:"HostPort"`
	ContainerPort Port     `json:"ContainerPort"`
	Protocol      Protocol `json:"Protocol"`
	PortName      string   `json:"PortName"`
}
