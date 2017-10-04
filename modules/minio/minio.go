package minio

import (
	"errors"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

// Minio client
type Minio struct {
	client *minio.Client
}

// S3 Struct
var S3 = &Minio{}

// NewEngine struct
func NewEngine(endpoint, accessID, secretKey string, ssl bool) error {

	if endpoint == "" || accessID == "" || secretKey == "" {
		return errors.New("endpoint, accessID and secretKey can't be empty")
	}

	// Initialize minio client object.
	client, err := minio.New(endpoint, accessID, secretKey, ssl)
	if err != nil {
		return err
	}

	S3.client = client

	return nil
}

// Upload file to s3
func (m *Minio) Upload(bucketName, objectName, filePath, contentType string) error {
	// Upload the zip file with FPutObject
	_, err := m.client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	return nil
}

// MakeBucket create bucket
func (m *Minio) MakeBucket(bucketName, region string) error {
	exists, err := m.client.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if exists {
		logrus.Infof("We already own %s bucket", bucketName)
		return nil
	}

	if err := m.client.MakeBucket(bucketName, region); err != nil {
		return err
	}
	logrus.Infof("Successfully created s3 bucket: %s", bucketName)

	return nil
}
