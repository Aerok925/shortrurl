package main

import (
	"github.com/Aerok925/shortrurl/internal/api"
	"github.com/Aerok925/shortrurl/internal/app"
	"github.com/Aerok925/shortrurl/internal/configs"
	"github.com/Aerok925/shortrurl/internal/inmemory"
	"github.com/Aerok925/shortrurl/internal/reducing"
	"go.uber.org/zap"
	"log"
)

func main() {
	cfg := configs.New()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	cache := inmemory.New()
	r := reducing.New()
	service := app.New(cache, r, logger, cfg.Result.BaseAddress)
	a := api.New(service, cfg.Server.Address, logger)
	a.Rout()
	if err := a.Start(); err != nil {
		logger.Error("Start server error: " + err.Error())
	}
}
