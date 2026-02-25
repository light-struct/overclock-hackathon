# AI Testing & Mentorship System - Backend

Улучшенная версия бэкенда с исправлением всех проблем безопасности и качества кода.

## Основные улучшения

### Безопасность
- ✅ Убраны hardcoded credentials
- ✅ Обязательная валидация всех переменных окружения
- ✅ Санитизация входных данных (защита от XSS и Log Injection)
- ✅ JWT middleware для защиты эндпоинтов (готов к использованию)
- ✅ Улучшенная обработка ошибок без утечки информации
- ✅ CORS настроен для конкретных доменов

### Архитектура
- ✅ Модуль конфигурации с валидацией
- ✅ Graceful shutdown
- ✅ Health check эндпоинт
- ✅ Индексы в БД для производительности
- ✅ Structured logging
- ✅ Dependency injection через параметры

## Структура проекта

```
backend/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/
│   │   └── config.go        # Конфигурация
│   ├── db/
│   │   └── postgres.go      # Подключение к БД
│   ├── middleware/
│   │   └── auth.go          # JWT middleware
│   ├── models/
│   │   └── models.go        # GORM модели
│   ├── repository/
│   │   └── repository.go    # Слой данных
│   ├── service/
│   │   ├── service.go       # Бизнес-логика
│   │   ├── auth_service.go  # Авторизация
│   │   └── gemini_service.go # AI интеграция
│   └── handlers/
│       └── handlers.go      # HTTP handlers
├── .env                     # Переменные окружения
├── go.mod
└── README.md
```

## Переменные окружения

Создайте файл `.env`:

```env
PORT=8080
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=require
GEMINI_API_KEY=your_gemini_api_key
JWT_SECRET=your_strong_secret_key
```

**Важно**: Все переменные обязательны. Сервер не запустится без них.

## Запуск

```bash
cd backend
go mod download
go run cmd/main.go
```

## API Endpoints

### Публичные
- `POST /api/signup` - Регистрация
- `POST /api/login` - Вход
- `GET /health` - Health check

### Защищенные (требуют JWT)
Для использования добавьте middleware в `handlers.go`:

```go
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(jwtSecret))
{
    protected.POST("/test/generate", h.GenerateTestHandler)
    protected.POST("/test/submit", h.submitTest)
    // ...
}
```

## Миграции

Автоматические миграции выполняются при запуске:
- `users` - пользователи
- `profiles` - профили
- `test_attempts` - попытки тестов

## Безопасность

### JWT Authentication
```go
// В запросах используйте заголовок:
Authorization: Bearer <your_jwt_token>
```

### Rate Limiting
Рекомендуется добавить middleware для ограничения запросов:
```bash
go get github.com/ulule/limiter/v3
```

## Логирование

Все ошибки логируются с санитизацией:
- AI ошибки
- Ошибки БД
- Ошибки валидации

## Production Checklist

- [ ] Настроить CORS для production доменов
- [ ] Добавить rate limiting
- [ ] Настроить structured logging (zap/logrus)
- [ ] Добавить мониторинг (Prometheus)
- [ ] Настроить CI/CD
- [ ] Добавить unit тесты
- [ ] Использовать secrets manager (AWS Secrets Manager)
- [ ] Настроить SSL/TLS

## Тестирование

```bash
# Запуск тестов (когда будут добавлены)
go test ./...

# Health check
curl http://localhost:8080/health
```

## Troubleshooting

### Ошибка подключения к БД
Проверьте `DATABASE_URL` в `.env` и доступность Neon БД.

### AI не работает
Проверьте `GEMINI_API_KEY` и квоты API.

### JWT ошибки
Убедитесь, что `JWT_SECRET` установлен и одинаков на всех инстансах.
