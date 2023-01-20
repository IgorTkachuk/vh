package minio_provider

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"vh/internal/image_storage"
	"vh/internal/models"
)

var _ image_storage.ImageStorage = &MinioProvider{}

type MinioProvider struct {
	minioAuthData
	client *minio.Client
}

func NewMinioProvider(minioUrl, minioUser, minioPassword string, ssl bool) (image_storage.ImageStorage, error) {
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
	m.client, err = minio.New(m.url, &minio.Options{Creds: credentials.NewStaticV4(m.user, m.password, ""), Secure: m.ssl})
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (m *MinioProvider) UploadFile(ctx context.Context, img models.ImageUnit) (string, error) {
	uploadInfo, err := m.client.PutObject(
		ctx,
		"videohosting",
		img.PayloadName,
		img.Payload,
		img.PayloadSize,
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	if err != nil {
		log.Fatal(err)
	}

	return uploadInfo.Location, err
}

func (m *MinioProvider) DownloadFile(ctx context.Context, imgId string) (*models.ImageUnit, error) {
	reader, err := m.client.GetObject(
		ctx,
		"video_hosting",
		imgId,
		minio.GetObjectOptions{},
	)

	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	return &models.ImageUnit{
		PayloadName: imgId,
		Payload:     reader,
	}, err
}
