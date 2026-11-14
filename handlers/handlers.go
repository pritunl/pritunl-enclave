package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/tools/errors"
	"github.com/pritunl/tools/logger"
)

type Token struct {
	Id  string `json:"id"`
	Ttl int64  `json:"ttl"`
}

func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50000)
}

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.WithFields(logger.Fields{
				"error": errors.New(fmt.Sprintf("%s", r)),
			}).Error("handlers: Handler panic")

			c.AbortWithStatus(http.StatusNotExtended)
		}
	}()

	c.Next()
}

func Errors(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Error("handlers: Handler error")
	}
}

func Cache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Next()
}

func Features(c *gin.Context) {
	if c.Request.ProtoMajor == 1 && c.Request.ProtoMinor == 0 {
		c.String(http.StatusHTTPVersionNotSupported, "HTTP/1.0 not supported")
		c.Abort()
		return
	}
	c.Next()
}

func Register(engine *gin.Engine) {
	engine.Use(Features)
	engine.Use(Limiter)
	engine.Use(Cache)
	engine.Use(Recovery)
	engine.Use(Errors)

	openAuth := engine.Group("")

	openAuth.GET("/check", checkGet)

	openAuth.GET("/robots.txt", robotsGet)
}
