package models

import "io"

type StorageObjectUnit struct {
	Payload     io.ReadSeeker
	PayloadName string
	PayloadSize int64
}
