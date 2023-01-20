package pgdb

import (
	"github.com/jmoiron/sqlx"
	"log"

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

func (p *PGS) AddObject(obj models.StorageObject) error {
	q := `
		INSERT INTO object (source_name, src_date, customer_pn, "user", addition)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	_, err := p.db.Exec(q, obj.SourceName, obj.SrcDate, obj.CustomerPN, obj.User, obj.Addition)
	if err != nil {
		logrus.Errorf("Can`t create object record in db: %v\n", err)
	}

	return err
}

func (p *PGS) GetObjectByCustomerPN(customerPn string) (obj []models.StorageObject, err error) {
	objList := make([]models.StorageObject, 0)
	q := `
		SELECT source_name, src_date, customer_pn, "user", addition
		FROM object
		WHERE customer_pn = $1
	`
	rows, err := p.db.Query(q, customerPn)
	if err != nil {
		logrus.Errorf("Can`t get objects records in db: %v\n", err)
		return objList, err
	}
	defer rows.Close()

	for rows.Next() {
		o := models.StorageObject{}
		err := rows.Scan(&o.SourceName, &o.SrcDate, &o.CustomerPN, &o.User, &o.Addition)

		if err != nil {
			logrus.Errorf("Can`t get object record from rows: %v\n", err)
			continue
		}

		objList = append(objList, o)
	}

	return objList, nil
}
