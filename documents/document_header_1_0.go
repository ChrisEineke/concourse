package documents

import (
	"encoding/json"

	"github.com/Masterminds/semver/v3"
)

type DocumentHeaderType = string
type DocumentHeaderVersion struct{ semver.Version }
type DocumentHeaderHandlerImage = string

type DocumentHeaderV1_0 struct {
	Type         DocumentHeaderType         `json:"type"`
	Version      DocumentHeaderVersion      `json:"version"`
	HandlerImage DocumentHeaderHandlerImage `json:"handlerImage"`
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
