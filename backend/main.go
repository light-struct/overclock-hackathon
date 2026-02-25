package main

import (
	"context"
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/server"
	"backend/internal/service"
)

func main() {
	cfg := config.Load()

	pool := db.NewPool(cfg.DatabaseURL)
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	testAttemptRepo := repository.NewTestAttemptRepository(pool)

	// Services
	authSvc := service.NewAuthService(cfg, userRepo)

	ctx := context.Background()
	aiSvc, err := service.NewAIService(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init AI service: %v", err)
	}
	examSvc := service.NewExamService(aiSvc, testAttemptRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	examHandler := handler.NewExamHandler(examSvc)

	// Router
	r := server.NewRouter(cfg, authHandler, examHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

