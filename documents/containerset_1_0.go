package documents

import (
	"github.com/ChrisEineke/concourse/utils"
	"github.com/kaptinlin/jsonschema"
)

var (
	ContainerSetV_1_0Schema *jsonschema.Schema = utils.Must(jsonschema.NewCompiler().Compile([]byte(`{
	  "$schema": "http://json-schema.org/draft-07/schema#",
	  "type": "object",
	  "properties": {
		"type": {
		  "type": "string",
		  "const": "containerSet"
		},
		"version": {
		  "type": "string",
		  "pattern": "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$",
		  "minLength": 6,
		  "maxLength": 256
		},
		"handlerImage": {
		  "type": "string"
		},
		"containerSet": {
		  "type": "object",
		  "properties": {}
		}
	  },
	  "required": ["type", "version", "handlerImage", "containerSet"]
	}`)))
)

type ContainerSetV_1_0 struct {
	Name          string                     `json:"Name"`
	Containers    []ContainerDefinitionV_1_0 `json:"Containers"`
	RestartPolicy RestartPolicy              `json:"RestartPolicy"`
	ID            string
}
