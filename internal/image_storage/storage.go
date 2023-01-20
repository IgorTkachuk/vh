package image_storage

import (
	"context"
	"vh/internal/models"
)

type ImageStorage interface {
	Connect() error
	UploadFile(ctx context.Context, img models.ImageUnit) (string, error)
	DownloadFile(ctx context.Context, imgId string) (*models.ImageUnit, error)
}
