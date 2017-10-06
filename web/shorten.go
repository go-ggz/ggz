package web

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/models"
	"github.com/go-ggz/ggz/modules/minio"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
)

// ShortenedIndex index page.
func ShortenedIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, config.Server.Host)
}

// FormURL URL Struct
type FormURL struct {
	URL string `json:"url" binding:"required,url"`
}

// CreateShortenURL create shorten url
func CreateShortenURL(c *gin.Context) {
	var data FormURL
	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		errorJSON(c, http.StatusBadRequest, errBadRequest)
		return
	}

	_, err := models.GetShortenFromURL(data.URL)

	if models.IsErrURLExist(err) {
		errorJSON(c, http.StatusBadRequest, errURLExist)
		return
	}

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	row, err := models.NewShortenURL(data.URL)

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	// upload QRCode image.
	go func(slug string) {
		if err := QRCodeGenerator(slug); err != nil {
			logrus.Errorln(err)
		}
	}(row.Slug)

	c.JSON(
		http.StatusBadRequest,
		row,
	)
}

// FetchShortenedURL show URL content
func FetchShortenedURL(c *gin.Context) {
	r, err := regexp.Compile(`^[a-zA-Z0-9]+$`)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":  http.StatusBadRequest,
				"error": "regexp not correct",
			},
		)
		return
	}

	slug := c.Param("slug")

	if !r.MatchString(slug) {
		errorJSON(c, http.StatusBadRequest, errSlugNotMatch)
		return
	}

	row := &models.Shorten{}

	has, err := row.GetFromSlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if !has {
		errorJSON(c, http.StatusNotFound, errSlugNotFound)
		return
	}

	c.JSON(
		http.StatusOK,
		row,
	)
}

// ShortenedURL redirect origin URL.
func ShortenedURL(c *gin.Context) {
	r, err := regexp.Compile(`^[a-zA-Z0-9]+$`)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":  http.StatusBadRequest,
				"error": "regexp not correct",
			},
		)
		return
	}

	slug := c.Param("slug")

	if !r.MatchString(slug) {
		errorJSON(c, http.StatusBadRequest, errSlugNotMatch)
		return
	}

	row := &models.Shorten{}

	has, err := row.GetFromSlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if !has {
		errorJSON(c, http.StatusNotFound, errSlugNotFound)
		return
	}

	err = row.UpdateHits(slug)
	if err != nil {
		errorJSON(c, http.StatusNotFound, errInternalServer)
		return
	}

	c.Redirect(http.StatusMovedPermanently, row.URL)
}

// QRCodeGenerator create QRCode
func QRCodeGenerator(slug string) error {
	objectName := fmt.Sprintf("%s.png", slug)
	filePath := ""
	host := strings.TrimRight(config.Server.ShortenHost, "/")

	switch config.Storage.Driver {
	case "disk":
		filePath = path.Join(
			config.Storage.Path,
			config.QRCode.Bucket,
			objectName,
		)
	case "s3":
		filePath = fmt.Sprintf("%s/%s", os.TempDir(), objectName)
	}

	if err := qrcode.WriteFile(
		host+"/"+slug,
		qrcode.Medium, 256, filePath); err != nil {
		return err
	}

	if config.Storage.Driver == "s3" {
		contentType := "image/png"
		if err := minio.S3.Upload(
			config.Minio.Bucket,
			objectName,
			filePath,
			contentType,
		); err != nil {
			return err
		}
	}

	return nil
}
