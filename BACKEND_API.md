# Backend API Structure

## Необходимые эндпоинты:

### 1. POST /auth/register
**Request:**
```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "string",
    "name": "string",
    "email": "string"
  }
}
```

### 2. POST /auth/login
**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "string",
    "name": "string",
    "email": "string"
  }
}
```

### 3. GET /auth/me
**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Response:**
```json
{
  "id": "string",
  "name": "string",
  "email": "string"
}
```

## Переменные окружения (.env):
```
DATABASE_URL=your_database_url
JWT_SECRET=your_jwt_secret
PORT=8000
```

## Примечания:
- JWT токены должны содержать user_id и email
- Токены должны иметь срок действия (например, 7 дней)
- Пароли должны хешироваться (bcrypt)
