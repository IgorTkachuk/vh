package models

import "time"

type StorageObjectMeta struct {
	Id          int       `json:"id,omitempty"`
	StorageName string    `json:"storage_name,omitempty"`
	OrigName    string    `json:"orig_name,omitempty"`
	OrigDate    time.Time `json:"orig_date"`
	AddDate     time.Time `json:"add_date"`
	BillingPn   string    `json:"billing_pn,omitempty"`
	UserName    string    `json:"user_name,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}
