package main

import (
	"Effective/config"
	"Effective/internal/repository"
	"Effective/internal/service"
	"Effective/internal/transport/http/handler"
	"Effective/internal/transport/server"
	"Effective/pkg/db"
	"Effective/pkg/migrations"

	"Effective/pkg/enricher"
	"Effective/pkg/logger"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	if err := run(); err != nil {
		log.Fatalf("Failed app start:%v", err)
	}

}

func run() error {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("create config ", zap.Error(err))
	}

	ctx := context.Background()
	conn, err := db.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer conn.Close()

	if err := migrations.Run(ctx, conn, logger); err != nil {
		logger.Fatal("Migrations failed", zap.Error(err))
	}

	enrich := enricher.New(logger)
	if enrich == nil {
		logger.Fatal("Enricher initialization failed", zap.Error(err))
	}
	repo := repository.NewPersonRepository(conn)
	service := service.NewPersonService(repo, enrich, logger)
	h := handler.NewPersonHandler(service, logger)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.POST("/person", h.CreatPerson)
	router.DELETE("/person/:id", h.DeletePerson)
	router.PATCH("/person/:id", h.UpdatePerson)
	router.GET("/person", h.GetPerson)

	srv := server.NewServer(cfg, logger, router)
	srv.Run()

	return nil
}
