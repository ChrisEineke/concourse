package documents

type ContainerName = string
type ContainerImage = string

type ContainerDefinitionV_1_0 struct {
	Name      ContainerName  `json:"name"`
	Image     ContainerImage `json:"image"`
	Ports     []Ports_1_0    `json:"ports"`
	Resources struct {
		Minimum ResourceRequirementsV_1_0 `json:"minimum"`
		Maximum ResourceRequirementsV_1_0 `json:"maximum"`
	} `json:"resources"`
}
