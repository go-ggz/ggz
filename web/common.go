package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func errorJSON(c *gin.Context, code int, err InnError) {
	logrus.Errorln(err.Error())

	c.AbortWithStatusJSON(
		code,
		gin.H{
			"code":  err.Code,
			"error": err.Error(),
		},
	)
}

func errorHTML(c *gin.Context, code int, file string, err InnError) {
	logrus.Errorln(err.Error())

	c.HTML(
		code,
		file,
		gin.H{
			"code":  err.Code,
			"error": err.Error(),
		},
	)
}
