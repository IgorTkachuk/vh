package object_storage

import (
	"context"
	"vh/internal/models"
)

type ImageStorage interface {
	Connect() error
	UploadFile(ctx context.Context, obj models.StorageObjectUnit, name string) (string, error)
	DownloadFile(ctx context.Context, imgId string) (*models.StorageObjectUnit, error)
	RemoveFile(ctx context.Context, objName string) error
	GetPresignedUrl(ctx context.Context, objName string) (string, error)
}
