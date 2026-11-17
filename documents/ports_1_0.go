package documents

import "github.com/moby/moby/api/types/network"

type Port = uint16
type Protocol = network.IPProtocol

type Ports_1_0 struct {
	HostPort      Port     `json:"hostPort"`
	ContainerPort Port     `json:"containerPort"`
	Protocol      Protocol `json:"protocol"`
}
