# Сводка по API проекта Делимед

## Анализ проекта

Проект представляет собой REST API для расчета доставки через различные транспортные компании (СДЭК, Деловые Линии).

### Структура проекта:
- **cmd/api/main.go** - точка входа приложения
- **internal/transport/handler/** - обработчики HTTP запросов
- **internal/transport/dto/** - структуры запросов и ответов
- **internal/service/** - бизнес-логика
- **internal/repository/** - работа с базой данных
- **internal/domain/** - доменные модели

## API Endpoints

### 1. Аутентификация

#### POST /register
Регистрация нового пользователя
- **Request**: `request.SignUpInput`
  - `username` (string, required) - имя пользователя
  - `password` (string, required, min=8) - пароль
  - `passwordConfirm` (string, required) - подтверждение пароля
- **Response**: 
  - 201: `{"message": "Registration successful."}`
  - 400: `{"error": "Invalid input format"}`
  - 409: `{"error": "..."}`

#### POST /login
Вход пользователя и получение JWT токена
- **Request**: `request.SignInInput`
  - `username` (string, required) - имя пользователя
  - `password` (string, required) - пароль
- **Response**: 
  - 200: `{"token": "jwt_token_here"}`
  - 400: `{"error": "Invalid input format"}`
  - 401: `{"error": "..."}`

### 2. Доставка (публичные endpoints)

#### GET /tariffslist
Получение списка доступных тарифов СДЭК
- **Request**: `request.CDEKRequestList` (query params или body)
  - `weight` (int, required) - вес в граммах
  - `length` (int, required) - длина в см
  - `width` (int, required) - ширина в см
  - `height` (int, required) - высота в см
  - `from_address` (string, required) - адрес отправления
  - `to_address` (string, required) - адрес доставки
- **Response**: 
  - 200: `response.CDEKTariffListResponse`
  - 400: `{"error": "invalid input format"}`
  - 500: `{"error": "..."}`

#### GET /tariffs
Расчет стоимости доставки по конкретному тарифу СДЭК
- **Request**: `request.CDEKRequest` (query params или body)
  - `tariff_code` (int, required) - код тарифа СДЭК
  - `weight` (int, required) - вес в граммах
  - `length` (int, required) - длина в см
  - `width` (int, required) - ширина в см
  - `height` (int, required) - высота в см
  - `from_address` (string, required) - адрес отправления
  - `to_address` (string, required) - адрес доставки
- **Response**: 
  - 200: `response.CDEKTariffCalcResponse`
  - 400: `{"error": "invalid input format"}`
  - 500: `{"error": "..."}`

#### POST /delivery/calculate
Единый расчет вариантов доставки от всех провайдеров
- **Request**: `request.DeliveryCalcRequest`
  - `length_cm` (int, required) - длина в см
  - `width_cm` (int, required) - ширина в см
  - `height_cm` (int, required) - высота в см
  - `weight_kg` (float64, required) - вес в кг
  - `delivery_type` (string, required) - "pickup" или "door"
  - `speed` (string) - "economy" / "express" / "urgent"
  - `from_address` (string, required) - адрес отправления
  - `to_address` (string, required) - адрес доставки
  - `shipment_date` (string) - дата отгрузки "YYYY-MM-DD"
  - `extra_services` (object) - дополнительные услуги
    - `insurance_value` (int64) - объявленная стоимость
    - `need_packing` (bool) - упаковка
    - `need_courier` (bool) - курьер
    - `need_documents` (bool) - работа с документами
    - `need_storage` (bool) - хранение
- **Response**: 
  - 200: `domain.FilterResult`
    - `status` (string) - "ok" или "error"
    - `options` (array) - массив `domain.DeliveryOption`
      - `provider` (string) - "cdek" или "dellin"
      - `tariff_code` (string) - код тарифа
      - `name` (string) - название тарифа
      - `delivery_type` (string) - "pickup" или "door"
      - `price` (int64) - цена в копейках
      - `currency` (string) - валюта (обычно "RUB")
      - `eta_from` (time, optional) - дата доставки "от"
      - `eta_to` (time, optional) - дата доставки "до"
  - 400: `{"error": "invalid input format"}`
  - 500: `{"error": "..."}`

### 3. Пользователь (защищенные endpoints, требуют JWT)

#### GET /api/user
Получение профиля текущего пользователя
- **Headers**: `Authorization: Bearer <token>`
- **Response**: 
  - 200: `response.UserResponse`
    - `id` (uuid) - UUID пользователя
    - `username` (string) - имя пользователя
    - `role` (string) - роль пользователя
    - `created_at` (time) - дата создания
    - `updated_at` (time) - дата обновления
  - 401: `{"error": "User not authenticated"}`
  - 404: `{"error": "User not found"}`
  - 500: `{"error": "Internal server error"}`

#### DELETE /api/user
Удаление профиля текущего пользователя
- **Headers**: `Authorization: Bearer <token>`
- **Response**: 
  - 200: `{"success": "User has deleted"}`
  - 401: `{"error": "User not authenticated"}`
  - 404: `{"error": "User not found"}`
  - 500: `{"error": "Internal server error"}`

## DTO Структуры

### Request DTOs

#### request.SignUpInput
```go
{
  "username": "string",
  "password": "string",
  "passwordConfirm": "string"
}
```

#### request.SignInInput
```go
{
  "username": "string",
  "password": "string"
}
```

#### request.DeliveryCalcRequest
```go
{
  "length_cm": 30,
  "width_cm": 20,
  "height_cm": 15,
  "weight_kg": 2.5,
  "delivery_type": "door",
  "speed": "economy",
  "from_address": "Москва, ул. Ленина, д. 1",
  "to_address": "Санкт-Петербург, Невский пр., д. 1",
  "shipment_date": "2024-01-15",
  "extra_services": {
    "insurance_value": 5000,
    "need_packing": false,
    "need_courier": false,
    "need_documents": false,
    "need_storage": false
  }
}
```

### Response DTOs

#### response.UserResponse
```go
{
  "id": "uuid",
  "username": "string",
  "role": "string",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

#### domain.FilterResult
```go
{
  "status": "ok",
  "options": [
    {
      "provider": "cdek",
      "tariff_code": "139",
      "name": "Экспресс-лайт",
      "delivery_type": "door",
      "price": 150000,
      "currency": "RUB",
      "eta_from": "2024-01-20T00:00:00Z",
      "eta_to": "2024-01-22T00:00:00Z"
    }
  ]
}
```

## Swagger документация

После генерации документации (см. SWAGGER.md), Swagger UI будет доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

В Swagger UI можно:
- Просмотреть все endpoints
- Увидеть структуры запросов и ответов
- Протестировать API прямо из браузера
- Авторизоваться с помощью JWT токена через кнопку "Authorize"

