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
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

var shortenPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

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
			log.Error().Err(err).Msg("QRCode Generator fail")
		}
	}(row.Slug)

	c.JSON(
		http.StatusOK,
		row,
	)
}

// FetchShortenedURL show URL content
func FetchShortenedURL(c *gin.Context) {
	slug := c.Param("slug")

	if !shortenPattern.MatchString(slug) {
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

// RedirectURL redirect shorten text to origin URL.
func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	if !shortenPattern.MatchString(slug) {
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
