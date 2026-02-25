package server

import (
	"time"

	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	examHandler *handler.ExamHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/")

	// Auth (public)
	authGroup := api.Group("/auth")
	authHandler.RegisterRoutes(authGroup)

	// Protected routes
	protected := api.Group("/")
	protected.Use(auth.JWTMiddleware(cfg))

	examGroup := protected.Group("/exam")
	examHandler.RegisterRoutes(examGroup)

	return r
}

