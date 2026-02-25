# QuizAgent - AI Testing Platform

## Изменения

### ✅ Реализовано:

1. **Цветовая схема**: Основной цвет изменен на #d50032
2. **Удалено поле API ключа**: API ключ Gemini теперь хранится в `.env.local` на бэкенде
3. **Локализация**: Добавлена поддержка 3 языков (Казахский, Русский, Английский)
4. **JWT Авторизация**: Регистрация и вход через бэкенд с JWT токенами
5. **Logout**: Добавлена кнопка выхода в header

### 🔧 Настройка Frontend

1. Создайте файл `.env.local` в папке `frontend/`:
```env
GEMINI_API_KEY=your_gemini_api_key_here
NEXT_PUBLIC_API_URL=http://localhost:8000
```

2. Установите зависимости:
```bash
cd frontend
npm install
```

3. Запустите проект:
```bash
npm run dev
```

### 🔧 Настройка Backend

Бэкенд должен предоставлять следующие эндпоинты:

- `POST /auth/register` - Регистрация пользователя
- `POST /auth/login` - Вход пользователя
- `GET /auth/me` - Получение данных текущего пользователя

Подробности в файле `BACKEND_API.md`

### 📝 Структура изменений:

**Новые файлы:**
- `lib/translations.ts` - Переводы на 3 языка
- `lib/app-context.tsx` - Контекст для языка и авторизации
- `app/login/page.tsx` - Страница входа
- `app/register/page.tsx` - Страница регистрации
- `middleware.ts` - Защита маршрутов
- `.env.local` - Переменные окружения
- `.env.example` - Пример переменных окружения

**Измененные файлы:**
- `app/globals.css` - Цвета изменены на #d50032
- `app/layout.tsx` - Добавлен AppProvider
- `components/header.tsx` - Добавлены переключатель языка и кнопка выхода
- `components/quiz-setup.tsx` - Удалено поле API ключа, добавлена локализация
- `app/quiz/page.tsx` - Использование JWT токена
- `app/api/quiz/generate/route.ts` - API ключ из .env
- `lib/quiz-types.ts` - Удален apiKey из QuizConfig

### 🌐 Переключение языка

Переключатель языка находится в header (иконка глобуса). Доступные языки:
- 🇬🇧 English
- 🇷🇺 Русский
- 🇰🇿 Қазақша

### 🔐 Авторизация

- Пользователи должны войти/зарегистрироваться для доступа к квизу
- JWT токен сохраняется в localStorage
- Кнопка выхода в header (иконка выхода)
- Защищенные маршруты: `/quiz/*`
