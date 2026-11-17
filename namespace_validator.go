package concourse

import (
	"encoding/json"
	"fmt"
	"strings"

	record "github.com/libp2p/go-libp2p-record"
)

type NamespaceValidator struct {
	record.NamespacedValidator
}

func NewNamespaceValidator() *NamespaceValidator {
	return &NamespaceValidator{}
}

func (ns *NamespaceValidator) Validate(key string, value []byte) error {
	if !strings.HasPrefix(key, "/concourse/") {
		return fmt.Errorf("the key %s is invalid", key)
	}
	var val any
	if err := json.Unmarshal(value, &val); err != nil {
		return fmt.Errorf("the value to be stored at key %s is not RJSON-compatible", key)
	}
	return nil
}

func (ns *NamespaceValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}

var _ = (record.Validator)((*NamespaceValidator)(nil))
