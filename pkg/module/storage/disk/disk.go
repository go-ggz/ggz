package disk

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

// Disk client
type Disk struct {
	Host string
	Path string
}

// NewEngine struct
func NewEngine(host, path string) *Disk {
	return &Disk{
		Host: host,
		Path: path,
	}
}

// UploadFile to s3 server
func (d *Disk) UploadFile(bucketName, fileName string, content []byte) error {
	return ioutil.WriteFile(d.FilePath(bucketName, fileName), content, os.FileMode(0644))
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
func (d *Disk) FilePath(bucketName, fileName string) string {
	return path.Join(
		d.Path,
		bucketName,
		fileName,
	)
}

// DeleteFile delete file
func (d *Disk) DeleteFile(bucketName, fileName string) error {
	return os.Remove(d.FilePath(bucketName, fileName))
}

// GetFile for storage host + bucket + filename
func (d *Disk) GetFile(bucketName, fileName string) string {
	if d.Host != "" {
		if u, err := url.Parse(d.Host); err == nil {
			u.Path = path.Join(u.Path, d.Path, bucketName, fileName)
			return u.String()
		}
	}
	return path.Join(d.Path, bucketName, fileName)
}
