package documents

import "github.com/ChrisEineke/concourse/types"

type ResourceRequirementsV_1_0 struct {
	Cpu    types.Quantity `json:"CPU"`
	Gpu    types.Quantity `json:"GPU"`
	Memory types.Quantity `json:"Memory"`
}
