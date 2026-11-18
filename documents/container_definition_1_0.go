package documents

type ContainerName = string
type ContainerImage = string

type ContainerDefinitionV_1_0 struct {
	Name      ContainerName  `json:"Name"`
	Image     ContainerImage `json:"Image"`
	Ports     []Ports_1_0    `json:"Ports"`
	Resources struct {
		Minimum ResourceRequirementsV_1_0 `json:"Minimum"`
		Maximum ResourceRequirementsV_1_0 `json:"Maximum"`
	} `json:"Resources"`
}
