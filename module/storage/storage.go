package storage

import (
	"github.com/go-ggz/ggz/module/storage/disk"
	"github.com/go-ggz/ggz/module/storage/minio"
	"github.com/go-ggz/ggz/pkg/config"
)

// Storage for s3 and disk
type Storage interface {
	// CreateBucket for create new folder
	CreateBucket(string, string) error
	// UploadFile for upload single file
	UploadFile(string, string, []byte) error
	// DeleteFile for delete single file
	DeleteFile(string, string) error
	// FilePath for store path + file name
	FilePath(string, string) string
	// GetFile for storage host + bucket + filename
	GetFile(string, string) string
}

// S3 for storage interface
var S3 Storage

// NewEngine return storage interface
func NewEngine() (Storage, error) {
	switch config.Storage.Driver {
	case "s3":
		return minio.NewEngine(
			config.Minio.EndPoint,
			config.Minio.AccessID,
			config.Minio.SecretKey,
			config.Minio.SSL,
		)
	case "disk":
		return disk.NewEngine(
			config.Server.Host,
			config.Storage.Path,
		), nil
	}

	return nil, nil
}
