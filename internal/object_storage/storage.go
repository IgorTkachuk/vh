package object_storage

import (
	"context"
	"vh/internal/models"
)

type ImageStorage interface {
	Connect() error
	UploadFile(ctx context.Context, obj models.StorageObjectUnit, name string) (string, error)
	DownloadFile(ctx context.Context, imgId string) (*models.StorageObjectUnit, error)
}
