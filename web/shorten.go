package web

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/storage"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

// ShortenedIndex index page.
func ShortenedIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, config.Server.Host)
}

// FormURL URL Struct
type FormURL struct {
	URL string `json:"url,omitempty" binding:"required,url"`
}

// CreateShortenURL create shorten url
func CreateShortenURL(c *gin.Context) {
	var data FormURL
	if err := c.ShouldBindJSON(&data); err != nil {
		errorJSON(c, http.StatusBadRequest, errBadRequest)
		return
	}

	row, err := model.GetShortenFromURL(data.URL)

	if model.IsErrURLExist(err) {
		c.JSON(
			http.StatusOK,
			row,
		)
		return
	}

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	user := helper.GetUserDataFromModel(c.Request.Context())
	row, err = model.NewShortenURL(data.URL, config.Server.ShortenSize, user)

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
		http.StatusOK,
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

	row, err := model.GetShortenBySlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if model.IsErrShortenNotExist(err) {
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

	row, err := model.GetShortenBySlug(slug)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	if model.IsErrShortenNotExist(err) {
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
	host := strings.TrimRight(config.Server.ShortenHost, "/")
	filePath := storage.S3.FilePath(objectName)
	png, err := qrcode.Encode(host+"/"+slug, qrcode.Medium, 256)
	if err != nil {
		return nil
	}

	return storage.S3.UploadFile(
		config.Minio.Bucket,
		objectName,
		filePath,
		png,
		"image/png",
	)
}
