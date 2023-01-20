package vh

import (
	"context"
	"vh/internal/db"
	"vh/internal/image_storage"
	"vh/internal/models"
)

type Vh struct {
	database db.DB
	storage  image_storage.ImageStorage
}

func NewVh(database db.DB, storage image_storage.ImageStorage) *Vh {
	return &Vh{database: database, storage: storage}
}

func (v *Vh) UploadVideo(ctx context.Context, img models.ImageUnit, obj models.StorageObject) error {
	_, err := v.storage.UploadFile(ctx, img)

	err = v.database.AddObject(obj)

	return err
}
