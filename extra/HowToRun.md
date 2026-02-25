# How to Run the Overclock Hackathon Project

## 1. Requirements

- **Backend**: Go 1.22+
- **Frontend**: Node.js 18+ (или новее), npm / pnpm / yarn
- **Database**: PostgreSQL (Neon, Supabase или локальный экземпляр)

## 2. Configure Environment Variables

### 2.1. Backend (`backend/.env`)

Файл уже существует и выглядит так:

```env
PORT=8080
DATABASE_URL=postgresql://<user>:<password>@<host>/<db>?sslmode=require
GEMINI_API_KEY=your_gemini_api_key_here
JWT_SECRET=super-anss-hack-2026
```

- **DATABASE_URL**: укажи свой URL PostgreSQL (например, Neon/Supabase).
  - В базе должны быть таблицы `users`, `profiles`, `test_attempts` из `notesSQL.txt`.
- **GEMINI_API_KEY**: ключ Gemini API (см. инструкцию в `geminiAPIfastintegration.txt`).
- **JWT_SECRET**: любой достаточно длинный секрет для подписи токенов.

> При старте backend автоматически подхватывает `.env` и пробрасывает `GEMINI_API_KEY` в `GOOGLE_API_KEY` для клиента `google.golang.org/genai`.

### 2.2. Frontend (`frontend/.env.local`)

Создай или отредактируй файл:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
GEMINI_API_KEY=your_gemini_api_key_here
```

- `NEXT_PUBLIC_API_URL` указывает на Go‑бэкенд.
- `GEMINI_API_KEY` нужен только, если какие‑то части фронта обращаются к Gemini напрямую.

## 3. Prepare Database

1. Подключись к своей базе PostgreSQL (через psql, DBeaver, Supabase SQL Editor и т.п.).
2. Выполни SQL‑скрипт из `notesSQL.txt` (копия лежит в `extra/notesSQL.txt`).
   - Он создаёт таблицы:
     - `users`
     - `profiles`
     - `test_attempts`
   - И добавляет пару тестовых пользователей (`admin@example.com`, `student@example.com`, пароль `"password"`).

## 4. Run the Backend (Go API)

Из корня репозитория:

```bash
cd backend
go run ./...
```

Бэкенд поднимается на порту из `PORT` (по умолчанию `:8080`), структура роутов:

- `POST /auth/register` — регистрация  
  Тело: `{ "name": string, "email": string, "password": string }`  
  Ответ: `{ "token": string, "user": { ... } }`

- `POST /auth/login` — логин  
  Тело: `{ "email": string, "password": string }`  
  Ответ: `{ "token": string, "user": { ... } }`

- `POST /auth/logout` — всегда 200, клиент просто удаляет токен.

- **Защищённые эндпоинты (с JWT в `Authorization: Bearer <token>`):**
  - `POST /exam/quiz/generate` — генерация квиза через Gemini  
    Тело: `{ "topic": string, "numQuestions": number, "difficulty": "easy" | "medium" | "hard" }`  
    Ответ: `{ "questions": Question[] }`, формат вопросов совпадает с фронтом.
  - `GET /exam/attempts` — список попыток текущего пользователя.
  - `POST /exam/attempts` — сохранение попытки теста (история/аналитика).

## 5. Run the Frontend (Next.js)

Из корня репозитория:

```bash
cd frontend
npm install        # или pnpm / yarn
npm run dev
```

Фронтенд поднимется на `http://localhost:3000`.

## 6. Full Flow Check

1. Открой `http://localhost:3000`.
2. **Регистрация**:
   - Перейди на `/register`, введи имя, email и пароль.
   - Фронт отправит `POST /auth/register` на Go‑бэкенд.
   - В ответ придёт `token`, который сохраняется в контекст приложения.
3. **Логин**:
   - Через `/login` можно войти существующим пользователем.
4. **Запуск квиза**:
   - Перейди на `/quiz`, настрой тему, сложность и количество вопросов.
   - Фронт отправит запрос на `POST /exam/quiz/generate` с `Authorization: Bearer <token>`.
   - Бэкенд вызовет Gemini и вернёт массив вопросов.
5. **Прохождение теста**:
   - Отвечай на вопросы, фронт может дополнительно использовать `/api/quiz/evaluate` (Next API) для детального анализа, а также может сохранить итоговую попытку через `POST /exam/attempts`.
6. **Проверка истории**:
   - Через `GET /exam/attempts` (можно дернуть из Postman/Insomnia) убедись, что попытки сохраняются в `test_attempts`.

## 7. Files in `extra/`

- `HowToRun.md` — этот файл с инструкциями по запуску.
- `geminiAPIfastintegration.txt` — заметки и пример кода по интеграции Gemini API на Go.
- `notesSQL.txt` — SQL‑скрипт создания таблиц `users`, `profiles`, `test_attempts` + демо‑данные.
- `ideaOfProjectAIAGENT.txt` — текст с основной идеей AI‑экзамен‑системы и AI‑агента (генерация теста, проверка, разбор ошибок, рекомендации).

Оригиналы этих файлов также лежат в корне репозитория; папка `extra` — это удобное место, где всё собранно в одном месте как документация.

## 8. About `authorize.go`

Файл `authorize.go` в корне — это пример реализации авторизации на базе **PocketBase**, он **не используется** текущим Go‑бэкендом:

- Модуль `go.mod` лежит в папке `backend`, поэтому `authorize.go` не попадает в сборку.
- В бэкенде реализована своя авторизация: Gin + JWT + PostgreSQL (`internal/auth`, `internal/service/auth_service.go` и т.д.).

Его можно оставить как справочный пример или перенести в `extra/` как архивный код, но на работу текущего сервера он не влияет.

