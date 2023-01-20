package models

import "io"

type StorageObjectUnit struct {
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
