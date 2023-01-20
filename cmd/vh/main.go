package main

import (
	"fmt"
	"log"
	"os"
	"vh/internal/controllers"
	pgdb "vh/internal/db/postgresql"
	"vh/internal/image_storage/minio_provider"
	"vh/internal/vh"
)

func main() {
	fileStorage, err := minio_provider.NewMinioProvider(
		os.Getenv("MINIOHOST"),
		os.Getenv("MINIOUSER"),
		os.Getenv("MINIOPASS"),
		false,
	)

	if err != nil {
		log.Fatal("Init image storage error.", err)
	}

	err = fileStorage.Connect()
	if err != nil {
		log.Fatal("Connect image storage error.", err)
	}

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, dbClose, err := pgdb.NewDatabase(dbinfo)
	if err != nil {
		log.Fatal("Init DB error: ", err)
	}
	defer dbClose()

	core := vh.NewVh(db, fileStorage)

	srv := controllers.NewServer(core)

	err = srv.Run(os.Getenv("SERVERPORT"))

	if err != nil {
		log.Fatal("Failed to run service.", err)
	}
}
