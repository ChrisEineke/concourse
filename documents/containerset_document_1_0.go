package documents

import (
	"fmt"
)

type RestartPolicy = int

const (
	RestartPolicyAlways = iota
	RestartPolicyOnFailure
	RestartPolicyNever
)

func RestartPolicyIntToString(restartPolicy int) (string, error) {
	if restartPolicy < RestartPolicyAlways || restartPolicy > RestartPolicyNever {
		return "", fmt.Errorf("invalid restart policy: %d", restartPolicy)
	}
	return [...]string{
		"RestartPolicyAlways",
		"RestartPolicyOnFailure",
		"RestartPolicyNever",
	}[restartPolicy], nil
}

type ContainerSetDocumentV_1_0 struct {
	DocumentHeaderV1_0
	ContainerSet ContainerSetV_1_0 `json:"containerSet"`
}
