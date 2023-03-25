package db

import "vh/internal/models"

type DB interface {
	AddObject(obj models.StorageObjectMeta) error
	RemoveObject(id int) error
	GetStorageNameById(id int) (storageName string, err error)
	GetObjectByBillingPN(customerPn string) (obj []models.StorageObjectMeta, err error)
}
