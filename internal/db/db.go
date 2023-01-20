package db

import "vh/internal/models"

type DB interface {
	AddObject(obj models.StorageObjectMeta) error
	GetObjectByBillingPN(customerPn string) (obj []models.StorageObjectMeta, err error)
}
