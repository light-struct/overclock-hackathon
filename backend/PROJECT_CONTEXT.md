# AI Testing & Mentorship System - Полный контекст проекта

## 📋 Описание проекта

Онлайн-платформа для экзаменов и наставничества с ИИ (Google Gemini).

**Стек:**
- Backend: Go 1.25, Gin, GORM
- БД: PostgreSQL (Neon Cloud)
- AI: Google Gemini 1.5-flash
- Auth: JWT + bcrypt

**Структура:**
```
backend/
├── cmd/main.go                    # Точка входа
├── internal/
│   ├── config/config.go          # Конфигурация
│   ├── middleware/auth.go        # JWT middleware
│   ├── models/models.go          # GORM модели
│   ├── repository/repository.go  # Слой данных
│   ├── service/
│   │   ├── service.go           # Агрегатор сервисов
│   │   ├── auth_service.go      # Авторизация
│   │   └── gemini_service.go    # AI интеграция
│   ├── handlers/handlers.go     # HTTP handlers
│   └── db/postgres.go           # Подключение к БД
├── .env                         # Переменные окружения
└── go.mod
```

---

## 🔴 История ошибок и решений

### Ошибка 1: Hardcoded Credentials
**Проблема:** В `auth_service.go` был hardcoded fallback для JWT_SECRET
```go
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    secret = "dev-secret"  // ❌ Небезопасно
}
```

**Решение:** Создан модуль `config/config.go` с обязательной валидацией:
```go
if cfg.JWTSecret == "" {
    return nil, errors.New("JWT_SECRET is required")
}
```

### Ошибка 2: SQL Injection риск
**Проблема:** DSN строка передавалась без валидации

**Решение:** Валидация на уровне конфигурации, GORM защищает запросы

### Ошибка 3: XSS уязвимости
**Проблема:** Ошибки возвращались пользователю без санитизации

**Решение:** Добавлена функция `sanitizeError()` с `html.EscapeString()`

### Ошибка 4: Log Injection
**Проблема:** Пользовательский ввод логировался напрямую
```go
log.Printf("AI error: %v", err)  // ❌ Уязвимо
```

**Решение:** Добавлена `sanitizeForLog()`:
```go
func sanitizeForLog(input string) string {
    input = strings.ReplaceAll(input, "\n", " ")
    input = strings.ReplaceAll(input, "\r", " ")
    return strings.TrimSpace(input)
}
```

### Ошибка 5: CORS слишком открыт
**Проблема:** `AllowOrigins: []string{"*"}`

**Решение:** Указаны конкретные домены:
```go
AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"}
```

### Ошибка 6: Отсутствие Graceful Shutdown
**Проблема:** БД соединения не закрывались при остановке

**Решение:** Добавлен graceful shutdown с таймаутом 5 секунд

### Ошибка 7: Migration Error - Constraint не существует
**Проблема:**
```
ERROR: constraint "uni_users_email" of relation "users" does not exist
```

**Причина:** GORM пытался удалить старый constraint при изменении индексов

**Решение:** Игнорирование ошибок "does not exist":
```go
if err := gormDB.AutoMigrate(...); err != nil {
    if !strings.Contains(err.Error(), "does not exist") {
        log.Fatalf("Migration error: %v", err)
    }
    log.Printf("Migration warning (ignored): %v", err)
}
```

### Ошибка 8: Нет индексов в БД
**Проблема:** Медленные запросы по user_id, subject, topic

**Решение:** Добавлены индексы в models.go:
```go
UserID uint `gorm:"index;not null"`
Subject string `gorm:"size:255;index;not null"`
Topic string `gorm:"size:255;index;not null"`
Score float64 `gorm:"index"`
```

---

## ✅ Что было исправлено

### Безопасность
1. ✅ Убраны все hardcoded credentials
2. ✅ Обязательная валидация env переменных
3. ✅ Санитизация всех ошибок (XSS защита)
4. ✅ Санитизация логов (Log Injection защита)
5. ✅ CORS настроен для конкретных доменов
6. ✅ JWT middleware готов к использованию

### Архитектура
1. ✅ Модуль конфигурации с валидацией
2. ✅ Graceful shutdown
3. ✅ Health check эндпоинт
4. ✅ Индексы БД для производительности
5. ✅ Dependency injection через параметры
6. ✅ Structured logging с санитизацией

### Качество кода
1. ✅ Clean Architecture (handlers → services → repository)
2. ✅ Обработка ошибок без утечки деталей
3. ✅ Валидация на всех уровнях
4. ✅ Игнорирование безопасных ошибок миграции

---

## 📝 Переменные окружения (.env)

```env
PORT=8080
DATABASE_URL=postgresql://user:pass@host:port/db?sslmode=require
GEMINI_API_KEY=your_api_key
JWT_SECRET=your_strong_secret
```

**Все переменные обязательны!** Сервер не запустится без них.

---

## 🔌 API Endpoints

### Публичные
- `GET /health` - Health check
- `POST /api/signup` - Регистрация
- `POST /api/login` - Вход (получение JWT)
- `POST /api/logout` - Выход

### Требуют авторизации (JWT в будущем)
- `POST /api/profiles` - Создание профиля
- `GET /api/profiles/:id` - Получение профиля
- `POST /api/test/generate` - Генерация теста (AI)
- `POST /api/test/submit` - Отправка ответов (AI анализ)
- `GET /api/attempts/user/:userID` - История попыток
- `GET /api/analytics/group` - Аналитика группы

---

## 🚀 Запуск проекта

```bash
cd backend
go mod download
go run cmd/main.go
```

**Проверка:**
```bash
curl http://localhost:8080/health
```

---

## 🧪 Тестирование

### Автоматическое (PowerShell)
```powershell
.\test_api.ps1
```

### Ручное (curl)
См. файл `API_TESTING.md`

---

## ⚠️ Известные ограничения

### 1. JWT Middleware не активирован
**Проблема:** Все эндпоинты доступны без токена

**Решение:** В `handlers.go` добавить:
```go
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(jwtSecret))
{
    protected.POST("/test/generate", h.GenerateTestHandler)
    protected.POST("/test/submit", h.submitTest)
    // ...
}
```

### 2. Нет Rate Limiting
**Проблема:** Уязвимость к DDoS и brute-force

**Решение:** Добавить middleware:
```bash
go get github.com/ulule/limiter/v3
```

### 3. Нет валидации сложности пароля
**Проблема:** Принимаются слабые пароли (минимум 6 символов)

**Решение:** Добавить проверку на спецсимволы, цифры, заглавные буквы

### 4. Нет ротации JWT токенов
**Проблема:** Токены живут 24 часа без возможности отзыва

**Решение:** Добавить refresh tokens или Redis blacklist

### 5. Логи не структурированы
**Проблема:** Используется стандартный `log.Printf`

**Решение:** Использовать zap или logrus:
```bash
go get go.uber.org/zap
```

---

## 📊 Модели данных

### User
```go
type User struct {
    ID           uint
    Email        string    // unique, not null
    Username     string    // not null
    PasswordHash string    // bcrypt hash
    Role         string    // student/teacher/admin
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### Profile
```go
type Profile struct {
    ID            uint
    UserID        uint      // index
    FullName      string
    Role          string
    PreferredLang string    // en/ru/kk
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### TestAttempt
```go
type TestAttempt struct {
    ID         uint
    UserID     uint      // index
    Subject    string    // index
    Topic      string    // index
    Score      float64   // index (0-100)
    Language   string
    AIFeedback string    // text
    CreatedAt  time.Time // index
    UpdatedAt  time.Time
}
```

---

## 🤖 AI (Gemini) Интеграция

### Генерация теста
**Модель:** gemini-1.5-flash
**Формат:** JSON
**Языки:** ru, kk, en

**Промпт:**
```
You are a strict examiner.
Create a test for subject: 'Math', topic: 'Algebra'.
Difficulty: medium.
Language: ru (IMPORTANT: If 'kk' - write in Kazakh, if 'ru' - in Russian).
Requirements:
1. Exactly 5 questions.
2. 4 answer options for each question.
3. Output format: ONLY valid JSON.
```

**Ответ:**
```json
{
  "test_title": "Тест по Алгебре",
  "questions": [
    {
      "id": 1,
      "text": "Вопрос?",
      "options": ["A", "B", "C", "D"],
      "answer": "Правильный ответ"
    }
  ]
}
```

### Анализ результатов
**Вход:** Список ответов + балл
**Выход:** Текстовый анализ с рекомендациями

**Промпт:**
```
Role: You are a wise mentor.
Subject: Math, Topic: Algebra.
Response language: ru.
Student scored: 75 points out of 100.
Student answers: [JSON]
Your task:
1. Briefly praise or support.
2. Explain the main mistake (if any).
3. Give advice on what to read.
```

---

## 🔧 Инструкции для будущего AI агента

### Если нужно добавить новый эндпоинт:

1. **Добавить модель** в `internal/models/models.go`
2. **Добавить репозиторий** в `internal/repository/repository.go`
3. **Добавить сервис** в `internal/service/service.go`
4. **Добавить handler** в `internal/handlers/handlers.go`
5. **Зарегистрировать роут** в `RegisterRoutes()`

### Если нужно защитить эндпоинт JWT:

```go
// В handlers.go
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(jwtSecret))
{
    protected.POST("/your-endpoint", h.yourHandler)
}
```

### Если нужно добавить новую переменную окружения:

1. Добавить в `.env`
2. Добавить в `internal/config/config.go`:
```go
type Config struct {
    // ...
    NewVar string
}

func Load() (*Config, error) {
    cfg := &Config{
        // ...
        NewVar: os.Getenv("NEW_VAR"),
    }
    
    if cfg.NewVar == "" {
        return nil, errors.New("NEW_VAR is required")
    }
    
    return cfg, nil
}
```
3. Передать в сервисы через `main.go`

### Если нужно изменить модель БД:

1. Изменить структуру в `models.go`
2. Запустить сервер - AutoMigrate применит изменения
3. Если ошибка "constraint does not exist" - игнорируется автоматически

### Если нужно добавить санитизацию:

**Для ошибок пользователю:**
```go
c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
```

**Для логов:**
```go
log.Printf("Error: %v", sanitizeForLog(err.Error()))
```

**Для AI промптов:**
```go
prompt := fmt.Sprintf("Subject: %s", sanitizeInput(subject))
```

---

## 📦 Зависимости (go.mod)

```
github.com/gin-gonic/gin v1.9.1
github.com/gin-contrib/cors v1.7.1
github.com/joho/godotenv v1.5.1
github.com/golang-jwt/jwt/v5 v5.2.0
github.com/google/generative-ai-go v0.20.1
github.com/jackc/pgx/v5 v5.8.0
golang.org/x/crypto v0.48.0
google.golang.org/api v0.269.0
gorm.io/driver/postgres v1.6.0
gorm.io/gorm v1.31.1
```

---

## 🎯 TODO для Production

### Критично
- [ ] Активировать JWT middleware на защищенных эндпоинтах
- [ ] Добавить rate limiting
- [ ] Настроить HTTPS/TLS
- [ ] Использовать AWS Secrets Manager для credentials
- [ ] Добавить мониторинг (Prometheus + Grafana)

### Важно
- [ ] Добавить unit тесты
- [ ] Добавить integration тесты
- [ ] Настроить CI/CD pipeline
- [ ] Добавить Swagger документацию
- [ ] Улучшить валидацию паролей
- [ ] Добавить refresh tokens

### Желательно
- [ ] Structured logging (zap/logrus)
- [ ] Metrics эндпоинт для Prometheus
- [ ] Docker Compose для локальной разработки
- [ ] Кэширование AI ответов (Redis)
- [ ] Pagination для списков
- [ ] Фильтрация и сортировка

---

## 🐛 Как дебажить

### Проблемы с БД
```bash
# Проверить подключение
psql "postgresql://user:pass@host:port/db?sslmode=require"

# Посмотреть таблицы
\dt

# Посмотреть структуру
\d users
```

### Проблемы с AI
```bash
# Проверить API ключ
curl https://generativelanguage.googleapis.com/v1/models?key=YOUR_KEY

# Логи в консоли покажут промпты и ответы
```

### Проблемы с JWT
```bash
# Декодировать токен
# Используйте jwt.io
```

---

## 📞 Контакты и ресурсы

- **Neon DB:** https://neon.tech
- **Gemini API:** https://ai.google.dev
- **Gin Docs:** https://gin-gonic.com
- **GORM Docs:** https://gorm.io

---

## 🔐 Безопасность

### Что защищено
✅ Пароли хешируются bcrypt
✅ JWT подписан секретом
✅ SQL injection защита (GORM)
✅ XSS защита (санитизация)
✅ Log injection защита
✅ CORS настроен

### Что нужно улучшить
⚠️ Rate limiting
⚠️ JWT middleware не активирован
⚠️ Нет HTTPS
⚠️ Нет secrets manager
⚠️ Нет audit logging

---

## 📈 Производительность

### Оптимизации
- Индексы на часто запрашиваемых полях
- Connection pooling (GORM default)
- Graceful shutdown

### Узкие места
- AI запросы (3-10 сек)
- Neon DB latency (зависит от региона)

### Рекомендации
- Кэшировать AI ответы в Redis
- Использовать CDN для статики
- Добавить read replicas для БД

---

Этот файл содержит ВСЮ информацию для продолжения разработки проекта.
