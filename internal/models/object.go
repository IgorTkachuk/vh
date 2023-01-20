package models

import "time"

type StorageObject struct {
	SourceName string    `json:"source_name,omitempty"`
	SrcDate    time.Time `json:"src_date"`
	CustomerPN string    `json:"customer_pn,omitempty"`
	User       string    `json:"user,omitempty"`
	Addition   string    `json:" ,omitempty"`
}
