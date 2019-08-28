package helper

import (
	"fmt"
	"strings"

	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/module/storage"

	"github.com/skip2/go-qrcode"
)

// QRCodeGenerator create QRCode
func QRCodeGenerator(slug string) error {
	objectName := fmt.Sprintf("%s.png", slug)
	host := strings.TrimRight(config.Server.ShortenHost, "/")
	png, err := qrcode.Encode(host+"/"+slug, qrcode.Medium, 256)
	if err != nil {
		return nil
	}

	return storage.S3.UploadFile(
		config.Minio.Bucket,
		objectName,
		png,
	)
}
