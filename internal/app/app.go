package app

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upbit/internal/config"
	v1 "upbit/internal/http/v1"
	"upbit/internal/metrics"
	"upbit/internal/server"
	"upbit/pkg/log"
	"upbit/pkg/rabbitmq"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to load config: %v", err))
	}
	log.Logger.Info(fmt.Sprintf("Config FILE --> ", cfg))

	rabbitConnect := rabbitmq.NewConnectWithRetries(cfg)
	rabbitProducer, _ := rabbitmq.NewProducer(cfg, rabbitConnect)

	handler := v1.NewHandler(cfg, rabbitProducer)
	srv := server.NewServer(cfg, handler.Routes())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Error("error occurred while running http server: %s\n", zap.Error(err))
		}
	}()

	go metrics.UpdateResourceUsageMetrics(3)

	log.Logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Logger.Error("failed to stop server: %v", zap.Error(err))
	}
}
