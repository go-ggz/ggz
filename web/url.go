package web

import (
	"net/http"

	"github.com/go-ggz/ggz/modules/meta"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// URLMeta for fetch metadata from URL
func URLMeta(c *gin.Context) {
	var data FormURL
	if err := c.ShouldBindJSON(&data); err != nil {
		errorJSON(c, http.StatusBadRequest, errBadRequest)
		return
	}

	metaData, err := meta.FetchData(data.URL)
	logrus.Debugf("%#v", metaData)

	if err != nil {
		errorJSON(c, http.StatusInternalServerError, errInternalServer)
		return
	}

	c.JSON(
		http.StatusOK,
		metaData,
	)
}
