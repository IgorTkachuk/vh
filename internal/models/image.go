package models

import "io"

type ImageUnit struct {
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
