package minio_provider

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"time"
	"vh/internal/models"
	"vh/internal/object_storage"
)

var _ object_storage.ImageStorage = &MinioProvider{}

type MinioProvider struct {
	minioAuthData
	client *minio.Client
}

func NewMinioProvider(minioUrl, minioUser, minioPassword string, ssl bool) (object_storage.ImageStorage, error) {
	return &MinioProvider{
		minioAuthData: minioAuthData{
			url:      minioUrl,
			user:     minioUser,
			password: minioPassword,
			ssl:      ssl,
		},
	}, nil
}

type minioAuthData struct {
	url      string
	user     string
	password string
	token    string
	ssl      bool
}

func (m *MinioProvider) Connect() error {
	var err error
	m.client, err = minio.New(fmt.Sprintf("%s:9000", m.url), &minio.Options{Creds: credentials.NewStaticV4(m.user, m.password, ""), Secure: m.ssl})
	if err != nil {
		log.Fatal(err)
	}

	if m.client.IsOnline() {
		fmt.Println("Mninio is online")
	} else {
		fmt.Println("Mninio is offline")
	}

	return err
}

func (m *MinioProvider) UploadFile(ctx context.Context, img models.StorageObjectUnit, name string) (string, error) {
	uploadInfo, err := m.client.PutObject(
		ctx,
		"videohosting",
		name,
		img.Payload,
		img.PayloadSize,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		log.Fatal(err)
	}

	return uploadInfo.Location, err
}

func (m *MinioProvider) DownloadFile(ctx context.Context, objId string) (*models.StorageObjectUnit, error) {
	reader, err := m.client.GetObject(
		ctx,
		"videohosting",
		objId,
		minio.GetObjectOptions{},
	)

	if err != nil {
		log.Fatal(err)
	}

	//defer reader.Close()

	return &models.StorageObjectUnit{
		PayloadName: objId,
		Payload:     reader,
	}, err
}

func (m *MinioProvider) GetPresignedUrl(ctx context.Context, objId string) (string, error) {
	// Use next and PresignedGetObject reqParam parameter for declare browser to save object as file.
	// It will use Content-Desposition HTTP header (inline(default) or attachment) for this functionality.
	// Refernce: https://developer.mozilla.org/ru/docs/Web/HTTP/Headers/Content-Disposition
	// reqParams := make(url.Values)
	// reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objId))

	presignedURL, err := m.client.PresignedGetObject(ctx, "videohosting", objId, time.Minute*5, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *MinioProvider) RemoveFile(ctx context.Context, objName string) error {
	return m.client.RemoveObject(ctx, "videohosting", objName, minio.RemoveObjectOptions{})
}
