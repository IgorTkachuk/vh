package pgdb

import (
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"vh/internal/db"
	"vh/internal/models"
)

var _ db.DB = &PGS{}

type PGS struct {
	db *sqlx.DB
}

func NewDatabase(dsn string) (db.DB, func() error, error) {
	dbc, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &PGS{db: dbc}, dbc.Close, err
}

func (p *PGS) AddObject(meta models.StorageObjectMeta) error {
	q := `
		INSERT INTO object (storage_name, orig_name, orig_date, add_date, billing_pn, user_name, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	_, err := p.db.Exec(q,
		meta.StorageName,
		meta.OrigName,
		meta.OrigDate,
		time.Now(),
		meta.BillingPn,
		meta.UserName,
		meta.Notes,
	)
	if err != nil {
		logrus.Errorf("Can`t create object record in db: %v\n", err)
	}

	return err
}

func (p *PGS) GetObjectByBillingPN(billingPn string) (obj []models.StorageObjectMeta, err error) {
	objList := make([]models.StorageObjectMeta, 0)
	q := `
		SELECT id, storage_name, orig_name, orig_date, add_date, billing_pn, user_name, notes
		FROM object
		WHERE billing_pn = $1
	`
	rows, err := p.db.Query(q, billingPn)
	if err != nil {
		logrus.Errorf("Can`t get objects records in db: %v\n", err)
		return objList, err
	}
	defer rows.Close()

	for rows.Next() {
		o := models.StorageObjectMeta{}
		err := rows.Scan(
			&o.Id, &o.StorageName, &o.OrigName,
			&o.OrigDate, &o.AddDate, &o.BillingPn,
			&o.UserName, &o.Notes,
		)

		if err != nil {
			logrus.Errorf("Can`t get object record from rows: %v\n", err)
			continue
		}

		objList = append(objList, o)
	}

	return objList, nil
}
