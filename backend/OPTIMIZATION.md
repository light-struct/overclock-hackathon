# Оптимизация производительности

## Проблема
GORM AutoMigrate делает 8+ медленных SQL-запросов (200-600ms каждый) при каждом запуске приложения.

## Решения

### 1. ✅ Оптимизация GORM (Применено)
```go
// internal/db/postgres.go
cfg := &gorm.Config{
    Logger:                 logger.Default.LogMode(logger.Warn), // Только предупреждения
    SkipDefaultTransaction: true,                                 // Отключить транзакции по умолчанию
    PrepareStmt:            true,                                 // Кэшировать prepared statements
}
```

### 2. ✅ Условные миграции (Применено)
```bash
# Пропустить миграции при запуске
SKIP_MIGRATIONS=true go run cmd/main.go

# Запустить миграции (по умолчанию)
go run cmd/main.go
```

### 3. Рекомендации для продакшена

#### Вариант A: Отдельная команда для миграций
```go
// cmd/migrate/main.go
func main() {
    // Запускать миграции отдельно
    gormDB.AutoMigrate(&models.User{}, &models.Profile{}, &models.TestAttempt{})
}
```

#### Вариант B: Использовать golang-migrate
```bash
migrate -path ./migrations -database "postgres://..." up
```

### 4. Дополнительные оптимизации

#### Connection Pool
```go
sqlDB, _ := gormDB.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

#### Индексы (уже есть в моделях)
- ✅ `users.email` - uniqueIndex
- ✅ `profiles.user_id` - index
- ✅ `test_attempts.user_id` - index
- ✅ `test_attempts.subject` - index
- ✅ `test_attempts.topic` - index
- ✅ `test_attempts.created_at` - index

## Результат
- Запуск с миграциями: ~3 секунды
- Запуск без миграций: ~0.5 секунды (6x быстрее)
