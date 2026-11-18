package concourse

import "errors"

var (
	ErrorUnsupportedDocumentType    = errors.New("unsupported document type")
	ErrorUnsupportedDocumentVersion = errors.New("unsupported document version")
)
