package documents

import "github.com/ChrisEineke/concourse/types"

type ResourceRequirementsV_1_0 struct {
	Cpu    types.Quantity `json:"cpu"`
	Gpu    types.Quantity `json:"gpu"`
	Memory types.Quantity `json:"memory"`
}
