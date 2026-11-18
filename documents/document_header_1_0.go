package documents

import (
	"encoding/json"

	"github.com/Masterminds/semver/v3"
)

const (
	ContainerSetDocumentType = "ContainerSet"
)

var (
	ContainerSetV_1_0_Version = semver.MustParse("1.0")
)

type DocumentHeaderType = string
type DocumentHeaderVersion struct{ semver.Version }
type DocumentHeaderHandlerImage = string

type DocumentHeaderV1_0 struct {
	Type         DocumentHeaderType         `json:"Type"`
	Version      DocumentHeaderVersion      `json:"Version"`
	HandlerImage DocumentHeaderHandlerImage `json:"HandlerImage"`
}

func (u *DocumentHeaderVersion) UnmarshalRJSON(data []byte) error {
	var versionStr string
	if err := json.Unmarshal(data, &versionStr); err != nil {
		return err
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return err
	}
	u.Version = *version
	return nil
}

func (u *DocumentHeaderV1_0) MarshalRJSON() ([]byte, error) {
	return ([]byte)(u.Version.String()), nil
}
