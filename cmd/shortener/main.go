package main

import (
	"github.com/Aerok925/shortrurl/internal/api"
	"github.com/Aerok925/shortrurl/internal/app"
	"github.com/Aerok925/shortrurl/internal/inmemory"
	"github.com/Aerok925/shortrurl/internal/reducing"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	cache := inmemory.New()
	r := reducing.New()
	service := app.New(cache, r, logger, "localhost:8080")
	a := api.New(service, "localhost", ":8080")
	a.Rout()
	logger.Debug("Server is starting")
	if err := a.Start(); err != nil {
		logger.Error("Start server error: " + err.Error())
	}
}
