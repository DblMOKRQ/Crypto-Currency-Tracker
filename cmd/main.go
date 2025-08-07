package main

import (
	"awesomeProject/api"
	"awesomeProject/internal/config"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/router"
	"awesomeProject/internal/router/handler"
	"awesomeProject/internal/scheduler"
	"awesomeProject/internal/service"
	logger "awesomeProject/pkg"
	"go.uber.org/zap"
)

func main() {
	// TODO: Написать документацию

	log, err := logger.NewLogger("debug")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Debug("Initializing config")
	cfg := config.MustLoad()
	log.Debug("Config initialized")

	storage, err := repository.NewStorage(cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SslMode, log)
	if err != nil {
		log.Fatal("Error initializing storage", zap.Error(err))
	}
	defer storage.Close()

	repo := storage.NewRepository()
	coinService := service.NewCoinService(repo)
	coinHandler := handler.NewHandler(coinService)
	rout := router.NewRouter(coinHandler, log)

	geckoApi := api.NewCoinGeckoApi(cfg.ApiKey, cfg.VsCurrency, log)
	if geckoApi.Init() != nil {
		log.Fatal("Error initializing coingecko api", zap.Error(err))
	}
	pricePoller := scheduler.NewPricePoller(coinService, cfg.PriceUpdates, geckoApi, log)
	go pricePoller.Start(cfg.MaxConcurrent)
	if err := rout.RunRouter(cfg.Address); err != nil {
		log.Fatal("Error initializing router", zap.Error(err))
	}
}
