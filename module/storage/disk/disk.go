package disk

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

// Disk client
type Disk struct {
	Host   string
	Path   string
	Bucket string
}

// NewEngine struct
func NewEngine(host, path, bucket string) *Disk {
	return &Disk{
		Host:   host,
		Path:   path,
		Bucket: bucket,
	}
}

// UploadFile to s3 server
func (d *Disk) UploadFile(_, _, filePath string, content []byte, _ string) error {
	return ioutil.WriteFile(filePath, content, os.FileMode(0644))
}

// CreateBucket create bucket
func (d *Disk) CreateBucket(bucketName, region string) error {
	storage := path.Join(d.Path, bucketName)
	if err := os.MkdirAll(storage, os.ModePerm); err != nil {
		return nil
	}
	log.Info().Msgf("Successfully created disk path: %s", storage)

	return nil
}

// FilePath for store path + file name
func (d *Disk) FilePath(fileName string) string {
	return path.Join(
		d.Path,
		d.Bucket,
		fileName,
	)
}

// DeleteFile delete file
func (d *Disk) DeleteFile(bucketName, fileName string) error {
	filePath := d.FilePath(fileName)
	return os.Remove(filePath)
}

// GetFile for storage host + bucket + filename
func (d *Disk) GetFile(bucketName, fileName string) string {
	return d.Host + "/" + d.Path + "/" + bucketName + "/" + fileName
}
