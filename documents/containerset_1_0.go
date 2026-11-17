package documents

type ContainerSetV_1_0 struct {
	Name          string                     `json:"name"`
	Containers    []ContainerDefinitionV_1_0 `json:"containers"`
	RestartPolicy RestartPolicy              `json:"restartPolicy"`
	ID            string
}
