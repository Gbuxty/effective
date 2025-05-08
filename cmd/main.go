// @title Effective API
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"Effective/config"
	"Effective/internal/repository"
	"Effective/internal/service"
	"Effective/internal/transport/http/handler"
	"Effective/internal/transport/server"
	"Effective/pkg/db"
	"Effective/pkg/logger"
	"Effective/pkg/migrations"
	"context"
	"log"
	"time"

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
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("failed to sync logger: %v", zap.Error(err))
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("create config ", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := db.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer conn.Close()

	if err := migrations.Run(ctx, conn, logger); err != nil {
		logger.Fatal("Migrations failed", zap.Error(err))
	}

	enrich := service.NewEnricher(logger, cfg)
	repo := repository.NewPersonRepository(conn)
	service := service.NewPersonService(repo, logger, enrich)
	h := handler.NewPersonHandler(service, logger)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	v1 := router.Group("/api/v1")
	{
		v1.POST("/person", h.CreatePerson)
		v1.DELETE("/person/:id", h.DeletePerson)
		v1.PATCH("/person/:id", h.UpdatePerson)
		v1.GET("/persons", h.GetPersons)
	}

	srv := server.NewServer(cfg, logger, router)
	srv.Run()

	return nil
}
