# API Testing Guide

## Базовый URL
```
http://localhost:8080
```

## 1. Health Check
```bash
curl http://localhost:8080/health
```

## 2. Регистрация (Signup)
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@test.com",
    "password": "password123",
    "username": "TestStudent",
    "role": "student"
  }'
```

Ответ:
```json
{
  "id": 1,
  "email": "student@test.com",
  "username": "TestStudent",
  "role": "student"
}
```

## 3. Вход (Login)
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@test.com",
    "password": "password123"
  }'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "TestStudent",
  "role": "student"
}
```

**Сохраните токен для следующих запросов!**

## 4. Создание профиля
```bash
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "full_name": "Иван Иванов",
    "role": "student",
    "preferred_lang": "ru"
  }'
```

## 5. Получение профиля
```bash
curl http://localhost:8080/api/profiles/1
```

## 6. Генерация теста (AI)
```bash
curl -X POST http://localhost:8080/api/test/generate \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Математика",
    "topic": "Алгебра",
    "difficulty": "medium",
    "lang": "ru"
  }'
```

Ответ:
```json
{
  "test_title": "Тест по Алгебре",
  "questions": [
    {
      "id": 1,
      "text": "Чему равно 2x + 3 = 7?",
      "options": ["x=1", "x=2", "x=3", "x=4"],
      "answer": "x=2"
    }
  ]
}
```

## 7. Отправка ответов (Submit Test)
```bash
curl -X POST http://localhost:8080/api/test/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "subject": "Математика",
    "topic": "Алгебра",
    "language": "ru",
    "questions": [
      {
        "id": 1,
        "text": "Чему равно 2x + 3 = 7?",
        "correct_answer": "x=2",
        "user_answer": "x=2"
      },
      {
        "id": 2,
        "text": "Решите уравнение x^2 = 9",
        "correct_answer": "x=3",
        "user_answer": "x=4"
      }
    ]
  }'
```

Ответ:
```json
{
  "attempt": {
    "id": 1,
    "user_id": 1,
    "subject": "Математика",
    "topic": "Алгебра",
    "score": 50.0,
    "language": "ru",
    "ai_feedback": "Хорошая попытка! Вы правильно решили первое уравнение..."
  },
  "topic_mastery": "Хорошая попытка! Вы правильно решили первое уравнение..."
}
```

## 8. История попыток пользователя
```bash
curl http://localhost:8080/api/attempts/user/1
```

## 9. Аналитика группы
```bash
curl http://localhost:8080/api/analytics/group
```

Ответ:
```json
{
  "average_score": 75.5,
  "weak_topics": ["Математика/Алгебра", "Физика/Механика"]
}
```

## 10. Выход (Logout)
```bash
curl -X POST http://localhost:8080/api/logout
```

---

## Полный сценарий тестирования

### Шаг 1: Проверка здоровья
```bash
curl http://localhost:8080/health
```

### Шаг 2: Регистрация студента
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"student1@test.com","password":"pass123","username":"Student1","role":"student"}'
```

### Шаг 3: Вход
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student1@test.com","password":"pass123"}' | jq -r '.token')

echo "Token: $TOKEN"
```

### Шаг 4: Генерация теста на русском
```bash
curl -X POST http://localhost:8080/api/test/generate \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "История",
    "topic": "Древний Рим",
    "difficulty": "easy",
    "lang": "ru"
  }' | jq
```

### Шаг 5: Генерация теста на казахском
```bash
curl -X POST http://localhost:8080/api/test/generate \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Математика",
    "topic": "Геометрия",
    "difficulty": "medium",
    "lang": "kk"
  }' | jq
```

### Шаг 6: Отправка результатов
```bash
curl -X POST http://localhost:8080/api/test/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "subject": "История",
    "topic": "Древний Рим",
    "language": "ru",
    "questions": [
      {"id":1,"text":"Вопрос 1","correct_answer":"A","user_answer":"A"},
      {"id":2,"text":"Вопрос 2","correct_answer":"B","user_answer":"C"},
      {"id":3,"text":"Вопрос 3","correct_answer":"C","user_answer":"C"}
    ]
  }' | jq
```

### Шаг 7: Просмотр истории
```bash
curl http://localhost:8080/api/attempts/user/1 | jq
```

### Шаг 8: Аналитика
```bash
curl http://localhost:8080/api/analytics/group | jq
```

---

## Тестирование с Postman

Импортируйте эту коллекцию в Postman:

1. Создайте новую коллекцию "AI Testing System"
2. Добавьте переменную окружения:
   - `base_url`: `http://localhost:8080`
   - `token`: (будет заполнен после login)

3. Создайте запросы из примеров выше

---

## Проверка ошибок

### Неверный email при регистрации
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid","password":"pass123","username":"Test"}'
```

### Короткий пароль
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123","username":"Test"}'
```

### Неверные credentials при входе
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"wrong@test.com","password":"wrongpass"}'
```

### Отсутствующие поля
```bash
curl -X POST http://localhost:8080/api/test/generate \
  -H "Content-Type: application/json" \
  -d '{"subject":"Math"}'
```

---

## Мониторинг логов

Во время тестирования следите за логами сервера:
```bash
go run cmd/main.go
```

Вы увидите:
- Подключение к БД
- Версию PostgreSQL
- Запросы к AI (Gemini)
- Ошибки валидации
- HTTP запросы

---

## Автоматическое тестирование (Bash скрипт)

Создайте файл `test_api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Testing Health Check ==="
curl -s $BASE_URL/health | jq

echo -e "\n=== Testing Signup ==="
curl -s -X POST $BASE_URL/api/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"auto@test.com","password":"test123","username":"AutoTest"}' | jq

echo -e "\n=== Testing Login ==="
TOKEN=$(curl -s -X POST $BASE_URL/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"auto@test.com","password":"test123"}' | jq -r '.token')

echo "Token: $TOKEN"

echo -e "\n=== Testing Test Generation ==="
curl -s -X POST $BASE_URL/api/test/generate \
  -H "Content-Type: application/json" \
  -d '{"subject":"Math","topic":"Algebra","difficulty":"easy","lang":"en"}' | jq

echo -e "\n=== All tests completed ==="
```

Запуск:
```bash
chmod +x test_api.sh
./test_api.sh
```
