## AI Testing & Mentorship System — Техническое ТЗ (MVP)

Этот документ описывает текущую реализацию бэкенда и задаёт полное ТЗ, по которому систему можно переписать с нуля.

---

## 1. Цель системы

- **Назначение**: онлайн‑платформа для экзаменов и наставничества с ИИ.
- **Основной пользовательский поток**:
  1. Студент регистрируется / логинится.
  2. Выбирает предмет, тему, сложность, язык → ИИ (Gemini) генерирует JSON‑тест.
  3. Студент отвечает на вопросы → сервер считает балл, отправляет ответы в ИИ, получает анализ ошибок и рекомендации.
  4. Результат сохраняется в БД (для истории и аналитики).
  5. Преподаватель смотрит агрегированную аналитику по группе (средний балл, слабые темы).

---

## 2. Технологический стек

- **Язык**: Go.
- **Web‑фреймворк**: `github.com/gin-gonic/gin`.
- **БД**: PostgreSQL (Neon / Supabase) через:
  - ORM: `gorm.io/gorm` + `gorm.io/driver/postgres`.
  - Низкоуровневая проверка подключения: `github.com/jackc/pgx/v5`.
- **ИИ**: Google Gemini (`gemini-1.5-flash`) через `github.com/google/generative-ai-go/genai`.
- **Авторизация**: JWT (`github.com/golang-jwt/jwt/v5`) + bcrypt.
- **Конфигурация**: `.env` + `github.com/joho/godotenv`.
- **CORS**: `github.com/gin-contrib/cors`.

---

## 3. Структура проекта (Clean Architecture)

Все файлы Go находятся в папке `backend/` с модулем `module exam-system`.

- `cmd/`
  - `main.go` — входная точка, сборка DI‑графа и запуск Gin.
- `internal/`
  - `db/`
    - `postgres.go` — создание GORM‑соединения + проверка через pgx.
  - `models/`
    - `models.go` — GORM‑модели (User, Profile, TestAttempt).
  - `repository/`
    - `repository.go` — интерфейсы и реализации репозиториев (Users, Profiles, TestAttempts, аналитика по попыткам).
  - `service/`
    - `service.go` — агрегатор сервисов (`Services`), бизнес‑логика профилей/попыток/аналитики, создание GeminiService и AuthService.
    - `auth_service.go` — регистрация/логин пользователей, генерация JWT.
    - `gemini_service.go` — интеграция с Gemini: генерация теста и анализ попытки.
  - `handlers/`
    - `handlers.go` — Gin‑хендлеры и маршруты (`/api/...`).

Архитектурные слои:

- **Handlers (Transport)** → принимают HTTP‑запросы, валидируют JSON, вызывают сервисы.
- **Service (UseCases)** → бизнес‑логика, координация репозиториев и внешних сервисов (ИИ).
- **Repository (Data)** → обращение к PostgreSQL через GORM.
- **Models (Domain)** → структуры, отражающие таблицы БД.

---

## 4. Доменные модели и БД

Таблицы предполагаются такими (через `AutoMigrate`):

- **users**
  - `id` (PK, uint)
  - `email` (уникальный, not null)
  - `username` (not null)
  - `password_hash` (not null)
  - `role` (string, not null) — `"student"`, `"teacher"`, `"admin"` и т.п.
  - `created_at`, `updated_at` (timestamps)

- **profiles**
  - `id` (PK)
  - `full_name`
  - `role` (роль по бизнес‑смыслу, может дублировать логическую роль пользователя)
  - `preferred_lang` (например, `"en"`, `"ru"`, `"kk"`)
  - `created_at`, `updated_at`

- **test_attempts**
  - `id` (PK)
  - `user_id` (FK → users.id или связь логическая)
  - `subject` (например, `"Math"`)
  - `topic` (например, `"Algebra"`)
  - `score` (float, 0–100)
  - `language` (код языка теста/ответов)
  - `ai_feedback` (text) — текстовый анализ/советы от ИИ.
  - `created_at`, `updated_at`

---

## 5. Сервисы

### 5.1 AuthService

- Интерфейс:

  - `Signup(ctx, email, password, username, role) (*User, error)`
  - `Login(ctx, email, password) (token string, user *User, err error)`

- Логика:
  - `Signup`:
    - Если `role` пустая → `"student"`.
    - Хеширует пароль `bcrypt`.
    - Создаёт `User` через `UserRepository.Create`.
  - `Login`:
    - Ищет пользователя по email.
    - Проверяет пароль через `bcrypt`.
    - Генерирует JWT с полями: `sub`, `email`, `username`, `role`, `exp (+24h)`.

### 5.2 GeminiService

- Интерфейс:

  - `GenerateTest(subject, topic, lang, difficulty string) (*GeneratedTest, error)`
  - `AnalyzeTest(subject, topic, lang string, answers []AnswerSummary, score float64) (string, error)`

- `GenerateTest`:
  - Вызывает модель `gemini-1.5-flash` c `ResponseMIMEType = "application/json"`.
  - Промпт строго требует JSON следующего вида (фактические ключи могут отличаться, но на переписывании стоит унифицировать):

    ```json
    {
      "title": "string",
      "questions": [
        {
          "id": 1,
          "text": "string",
          "options": ["string"],
          "answer": "string"
        }
      ]
    }
    ```

  - Поддерживает генерацию на разных языках, включая русский и казахский.

- `AnalyzeTest`:
  - На вход получает список ответов (`AnswerSummary`: вопрос, правильный ответ, ответ пользователя, флаг `is_correct`) и итоговый балл.
  - Сериализует ответы в JSON и передаёт в Gemini с инструкцией:
    - Объяснить ошибки.
    - Оценить «Topic Mastery».
    - Дать конкретные рекомендации по теме.
    - Ответить на заданном языке (`lang`).
  - Возвращает обычный текст (не JSON).

### 5.3 TestAttemptService

- Методы:
  - `CreateAttempt(ctx, *TestAttempt) error`
  - `ListAttemptsByUser(ctx, userID uint) ([]TestAttempt, error)`
  - `GetGroupAnalytics(ctx) (avgScore float64, weakTopics []string, error)`

- `GetGroupAnalytics`:
  - Использует агрегирующие запросы:
    - Средний балл по всем попыткам: `AVG(score)`.
    - «Слабые темы»: `GROUP BY subject, topic HAVING AVG(score) < 50`.
  - Формирует список строк `"Subject/Topic"` для слабых мест.

### 5.4 ProfileService

- Простые операции:
  - `CreateProfile(ctx, *Profile) error`
  - `GetProfile(ctx, id uint) (*Profile, error)`

---

## 6. HTTP API (эндпоинты)

Базовый префикс: `/api`.

### 6.1 Авторизация

- `POST /api/signup`
  - Тело: `{ "email", "password", "username", "role?" }`.
  - Действие: создаёт пользователя, по умолчанию роль `student`.
  - Ответ: `{ id, email, username, role }`.

- `POST /api/login`
  - Тело: `{ "email", "password" }`.
  - Действие: проверка пароля, генерация JWT.
  - Ответ: `{ "token", "username", "role" }`.

- `POST /api/logout`
  - Stateless JWT: просто возвращает сообщение «удалите токен на клиенте».

*(На данный момент нет middleware, которое проверяет JWT и ограничивает доступ по ролям — это надо добавить при переписывании.)*

### 6.2 Профили

- `POST /api/profiles`
  - Тело: `Profile` (full_name, role, preferred_lang).
  - Создаёт профиль в БД.

- `GET /api/profiles/:id`
  - Возвращает профиль по ID.

### 6.3 Генерация теста (ИИ)

- `POST /api/test/generate`
  - Тело: `{ "subject", "topic", "difficulty", "lang" }`.
  - Вызывает `GeminiService.GenerateTest`.
  - Ответ: JSON‑тест, который будет отображаться фронтом.

### 6.4 Отправка ответов (сдача теста)

- `POST /api/test/submit`
  - Тело:

    ```json
    {
      "user_id": 1,
      "subject": "Math",
      "topic": "Algebra",
      "language": "ru",
      "questions": [
        {
          "id": 1,
          "text": "2 + 2 = ?",
          "correct_answer": "4",
          "user_answer": "4"
        }
      ]
    }
    ```

  - Логика:
    - Считает количество правильных ответов (сравнение по строке, без учёта регистра и пробелов).
    - Вычисляет `score` (0–100%).
    - Формирует список `AnswerSummary` и передаёт в `GeminiService.AnalyzeTest`.
    - Создаёт запись `TestAttempt` с полями: user_id, subject, topic, score, language, ai_feedback.
  - Ответ:

    ```json
    {
      "attempt": { ... TestAttempt ... },
      "topic_mastery": "текстовый анализ от ИИ"
    }
    ```

### 6.5 История попыток и аналитика

- `POST /api/attempts`
  - Ручное создание попытки (вспомогательный/технический эндпоинт).

- `GET /api/attempts/user/:userID`
  - Возвращает список попыток конкретного пользователя.

- `GET /api/analytics/group`
  - Возвращает:

    ```json
    {
      "average_score": 75.0,
      "weak_topics": ["Math/Algebra", "History/Dates"]
    }
    ```

---

## 7. Конфигурация и окружение

Все переменные задаются в `backend/.env` (загружается в `cmd/main.go`):

- `PORT` — порт HTTP‑сервера (по умолчанию `8080`).
- `DATABASE_URL` — строка подключения к PostgreSQL, формат:

  ```text
  postgresql://user:password@host:port/dbname?sslmode=require
  ```

- `GEMINI_API_KEY` — API‑ключ для Gemini.
- `JWT_SECRET` — секрет для подписи JWT.

При запуске:

1. Загружается `.env`.
2. Создаётся соединение с БД через GORM и дополнительно проверяется через pgx (`SELECT version()`).
3. Выполняются `AutoMigrate` для `User`, `Profile`, `TestAttempt`.
4. Настраивается Gin + CORS (AllowOrigins `*`, стандартные методы/заголовки).
5. Регистрируются маршруты и запускается HTTP‑сервер.

---

## 8. Нефункциональные требования и ограничения

- **Формат API**: JSON, REST‑подобные маршруты.
- **CORS**: открыт для всех origins (для удобства фронтенда в MVP).
- **Локализация ИИ**: промпты рассчитаны на поддержку русского и казахского языков (а также английского).
- **Логи**:
  - Подключение к БД и версия Postgres логируются при старте.
  - Ошибки Gemini логируются, но пользователю возвращается аккуратное сообщение.
- **Безопасность (MVP)**:
  - Пароли хранятся только в виде bcrypt‑хеша.
  - JWT подписан, но нет ротации ключей и чёрных списков.
  - Нет полноценного middleware для проверки токена и ролей — это необходимо добавить при переработке.

---

## 9. Известные проблемы и точки для улучшения при переписывании

1. **Авторизация**:
   - Нет middleware для авторизации по JWT и разграничения ролей (student/teacher/admin).
   - Эндпоинты аналитики и генерации тестов пока доступны без проверки токена.
2. **Генерация тестов**:
   - Формат JSON от Gemini может отличаться от ожидаемого, требуется более строгая схема/валидатор.
   - Нет сохранения сгенерированного теста (только результат попытки).
3. **Структура БД**:
   - Нет явных внешних ключей и индексов, всё создаётся через `AutoMigrate`.
   - Связь между `users` и `profiles` пока не формализована.
4. **Обработка ошибок**:
   - Сообщения об ошибках частично на русском, частично на английском; нет единого стиля.
   - Не везде есть подробный лог ошибок.
5. **Тестирование**:
   - Нет unit‑ и integration‑тестов.
6. **Инфраструктура**:
   - Нет Docker/Docker Compose для локального запуска.
   - Нет миграций (например, через `golang-migrate`) — только GORM AutoMigrate.

---

## 10. Минимальный план переписывания

1. Создать новый модуль Go с аналогичной структурой (`cmd`, `internal/{db,models,repository,service,handlers}`).
2. Определить доменные модели и SQL‑схему (users, profiles, test_attempts) явно через миграции.
3. Реализовать:
   - AuthService + JWT middleware (доступ к аналитике только для `teacher`/`admin`).
   - GeminiService с жёстко заданной JSON‑схемой и валидацией.
   - TestAttemptService с расчётом баллов и аналитикой.
4. Реализовать набор HTTP‑эндпоинтов, перечисленных в разделе 6, с единым стилем ошибок и ответов.
5. Настроить конфигурацию через `.env` и добавить Docker‑окружение (Postgres + backend).
6. Написать базовые тесты (авторизация, генерация теста (mock Gemini), отправка попытки, аналитика).

