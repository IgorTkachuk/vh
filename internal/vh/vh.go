package vh

import (
	"context"
	"vh/internal/db"
	"vh/internal/models"
	"vh/internal/object_storage"
)

type Vh struct {
	database db.DB
	storage  object_storage.ImageStorage
}

func NewVh(database db.DB, storage object_storage.ImageStorage) *Vh {
	return &Vh{database: database, storage: storage}
}

func (v *Vh) UploadObject(ctx context.Context, obj models.StorageObjectUnit, objMeta models.StorageObjectMeta) error {
	objMeta.StorageName = GenerateObjectName(objMeta.BillingPn, objMeta.OrigName)
	_, err := v.storage.UploadFile(ctx, obj, objMeta.StorageName)

	err = v.database.AddObject(objMeta)

	return err
}

func (v *Vh) GetObjectByBillingPn(billingPn string) ([]models.StorageObjectMeta, error) {
	metaList, err := v.database.GetObjectByBillingPN(billingPn)
	if err != nil {
		return []models.StorageObjectMeta{}, err
	}

	return metaList, nil
}

func (v *Vh) GetContent(ctx context.Context, name string) (*models.StorageObjectUnit, error) {
	return v.storage.DownloadFile(ctx, name)
}

func (v *Vh) GetPresignedUrl(ctx context.Context, name string) (string, error) {
	return v.storage.GetPresignedUrl(ctx, name)
}

func (v *Vh) RemoveObject(ctx context.Context, id int) error {
	objName, err := v.database.GetStorageNameById(id)
	if err != nil {
		return err
	}

	err = v.storage.RemoveFile(ctx, objName)
	if err != nil {
		return err
	}

	err = v.database.RemoveObject(id)
	if err != nil {
		return err
	}

	return nil
}
