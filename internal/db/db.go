package db

import "vh/internal/models"

type DB interface {
	AddObject(obj models.StorageObject) error
	GetObjectByCustomerPN(customerPn string) (obj []models.StorageObject, err error)
}
