package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"exam-system/internal/db"
	"exam-system/internal/handlers"
	"exam-system/internal/repository"
	"exam-system/internal/service"
)

func main() {
	// 1. Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используются системные переменные")
	}

	// 2. Получаем строку подключения к БД
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Ошибка: Переменная DATABASE_URL не установлена в .env")
	}

	// 3. Подключаемся к базе (Neon/Supabase)
	gormDB, err := db.NewPostgresConnection(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// 4. Инициализируем слои приложения
	// Репозитории (работа с данными)
	repos := repository.NewRepositories(gormDB)
	
	// Сервисы (бизнес-логика + ИИ)
	// ВАЖНО: NewServices сам создаст GeminiService внутри себя, как мы и договорились
	services := service.NewServices(repos) 

	// Хендлеры (обработка HTTP запросов)
	handler := handlers.NewHandler(services)

	// 5. Настройка роутера и запуск
	r := gin.Default()

	// CORS для фронтенда (разрешаем все источники)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: false,
	}))

	// Регистрация всех путей (API Endpoints)
	handler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}