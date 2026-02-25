package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection creates a new GORM DB connection for PostgreSQL (Supabase).
// Дополнительно выполняет прямое подключение через pgx и запрос SELECT version()
// чтобы убедиться, что соединение с БД реально работает.
func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Основное подключение через GORM
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("connected to Postgres via GORM")

	// Дополнительная проверка через pgx: SELECT version()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		return nil, err
	}

	log.Println("Postgres version (pgx check):", version)

	return db, nil
}

