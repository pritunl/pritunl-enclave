package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-enclave/constants"
	"github.com/pritunl/pritunl-enclave/handlers"
	"github.com/pritunl/tools/logger"
)

func main() {
	logger.Init(
		logger.SetTimeFormat(""),
	)

	logger.AddHandler(func(record *logger.Record) {
		fmt.Print(record.StringColor())
	})

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	handlers.Register(router)

	server := &http.Server{
		Addr:              "127.0.0.1:" + constants.BindPort,
		Handler:           router,
		ReadTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       1 * time.Minute,
	}

	logger.WithFields(logger.Fields{
		"port": constants.BindPort,
	}).Info("main: Starting HTTP server")

	err := server.ListenAndServe()
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Error("main: Server error")
		panic(err)
	}
}
