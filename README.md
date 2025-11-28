# Делимед API

REST API для расчета стоимости доставки через различные транспортные компании (СДЭК, Деловые Линии). Проект предоставляет единый интерфейс для получения вариантов доставки от всех провайдеров с фильтрацией и сортировкой по цене.

## 🚀 Быстрый старт

### Требования

- Docker и Docker Compose
- Git

### Запуск проекта через Docker

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd Делимед
   ```

2. **Запустите проект:**
   ```bash
   docker-compose up -d
   ```

3. **Проверьте статус контейнеров:**
   ```bash
   docker-compose ps
   ```

4. **Просмотрите логи приложения:**
   ```bash
   docker-compose logs -f app
   ```

5. **Остановите проект:**
   ```bash
   docker-compose down
   ```

### Доступ к сервисам

После запуска проект будет доступен по следующим адресам:

- **API:** http://localhost:8080
- **Swagger UI:** http://localhost:8080/swagger/index.html
- **PostgreSQL:** localhost:5432
  - User: `admin`
  - Password: `mysteam2006`
  - Database: `delimed`

## 📋 API Endpoints

### Публичные endpoints

- `POST /register` - Регистрация нового пользователя
- `POST /login` - Вход и получение JWT токена
- `GET /tariffslist` - Получение списка тарифов СДЭК
- `GET /tariffs` - Расчет конкретного тарифа СДЭК
- `POST /delivery/calculate` - **Расчет вариантов доставки от всех провайдеров**

### Защищенные endpoints (требуют JWT токен)

- `GET /api/user` - Получение профиля пользователя
- `DELETE /api/user` - Удаление профиля пользователя

## 🧪 Примеры запросов для Postman

### 1. Регистрация пользователя

**Endpoint:** `POST http://localhost:8080/register`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON):**
```json
{
  "username": "testuser",
  "password": "password123",
  "passwordConfirm": "password123"
}
```

**Ожидаемый ответ (201):**
```json
{
  "message": "Registration successful."
}
```

---

### 2. Вход пользователя

**Endpoint:** `POST http://localhost:8080/login`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON):**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Ожидаемый ответ (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Сохраните токен для использования в защищенных запросах!**

---

### 3. Расчет вариантов доставки (ОБЯЗАТЕЛЬНЫЙ ТЕСТ)

**Endpoint:** `POST http://localhost:8080/delivery/calculate`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON) - Пример 1: Доставка склад → дверь**
```json
{
  "length_cm": 30,
  "width_cm": 20,
  "height_cm": 15,
  "weight_kg": 2.5,
  "to": "дверь",
  "from_address": "Москва, ул. Ленина, д. 1",
  "to_address": "Санкт-Петербург, Невский пр., д. 1",
  "shipment_date": "2024-01-15",
  "speed": "economy",
  "extra_services": {
    "insurance_value": 5000,
    "need_packing": false,
    "need_courier": false,
    "need_documents": false,
    "need_storage": false
  }
}
```

**Body (JSON) - Пример 2: Самовывоз (склад → склад)**
```json
{
  "length_cm": 50,
  "width_cm": 40,
  "height_cm": 30,
  "weight_kg": 5.0,
  "to": "склад",
  "from_address": "Москва, ул. Тверская, д. 10",
  "to_address": "Санкт-Петербург, Невский пр., д. 28",
  "shipment_date": "2024-01-20",
  "extra_services": {
    "insurance_value": 10000,
    "need_packing": true,
    "need_courier": false,
    "need_documents": false,
    "need_storage": false
  }
}
```

**Body (JSON) - Пример 3: Минимальный запрос (только обязательные поля)**
```json
{
  "length_cm": 30,
  "width_cm": 20,
  "height_cm": 15,
  "weight_kg": 2.5,
  "to": "дверь",
  "from_address": "Москва",
  "to_address": "Санкт-Петербург"
}
```

**Body (JSON) - Пример 4: С дополнительными услугами**
```json
{
  "length_cm": 25,
  "width_cm": 20,
  "height_cm": 10,
  "weight_kg": 1.5,
  "to": "дверь",
  "from_address": "Москва, ул. Арбат, д. 5",
  "to_address": "Санкт-Петербург, ул. Садовая, д. 12",
  "shipment_date": "2024-01-18",
  "extra_services": {
    "insurance_value": 15000,
    "need_packing": true,
    "need_courier": true,
    "need_documents": false,
    "need_storage": false
  }
}
```

**Ожидаемый ответ (200):**
```json
{
  "status": "ok",
  "options": [
    {
      "provider": "cdek",
      "tariff_code": "139",
      "name": "Экспресс-лайт",
      "delivery_type": "door",
      "from_type": "склад",
      "to_type": "дверь",
      "delivery_mode": 2,
      "price": 150000,
      "currency": "RUB",
      "eta_from": "2024-01-20T00:00:00Z",
      "eta_to": "2024-01-22T00:00:00Z"
    },
    {
      "provider": "dellin",
      "tariff_code": "auto",
      "name": "Деловые линии (Авто)",
      "delivery_type": "door",
      "from_type": "склад",
      "to_type": "дверь",
      "price": 145000,
      "currency": "RUB",
      "eta_from": "2024-01-21T00:00:00Z",
      "eta_to": "2024-01-23T00:00:00Z"
    }
  ]
}
```

**Важные замечания:**
- Поле `from` **не требуется** - всегда используется "склад" по умолчанию
- Указывается только поле `to`: "склад" или "дверь"
- Варианты доставки **автоматически сортируются по цене** (от меньшей к большей)
- Тарифы с "постомат" в названии **автоматически исключаются**
- Фильтрация происходит по значению `to`:
  - `to="дверь"` → только тарифы "склад-дверь"
  - `to="склад"` → только тарифы "склад-склад"

---

### 4. Получение профиля пользователя (защищенный endpoint)

**Endpoint:** `GET http://localhost:8080/api/user`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <ваш_jwt_токен>
```

**Ожидаемый ответ (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "testuser",
  "role": "user",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

### 5. Получение списка тарифов СДЭК

**Endpoint:** `GET http://localhost:8080/tariffslist`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON):**
```json
{
  "weight": 2500,
  "length": 30,
  "width": 20,
  "height": 15,
  "from_address": "Москва",
  "to_address": "Санкт-Петербург"
}
```

---

## 🏗️ Архитектура проекта

Проект следует принципам **Clean Architecture** и разделен на следующие слои:

```
delimed/
├── cmd/
│   └── api/              # Точка входа приложения
│       ├── main.go       # Инициализация и запуск сервера
│       └── swagger.go    # Swagger метаданные
│
├── config/               # Конфигурационные файлы
│   └── config.yaml       # Настройки приложения
│
├── docs/                 # Swagger документация (автогенерируемая)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal/             # Внутренняя логика приложения
│   ├── config/           # Загрузка и парсинг конфигурации
│   │   └── config.go
│   │
│   ├── domain/           # Доменные модели (бизнес-логика)
│   │   └── delivery.go   # Модели доставки
│   │
│   ├── repository/       # Слой доступа к данным
│   │   ├── models/       # Модели БД (GORM)
│   │   ├── postgres/     # Реализация для PostgreSQL
│   │   └── repository.go # Интерфейсы репозиториев
│   │
│   ├── service/          # Бизнес-логика (Use Cases)
│   │   ├── authservice/  # Сервис аутентификации
│   │   ├── userservice/  # Сервис пользователей
│   │   ├── deliveryservice/ # Сервис расчета доставки
│   │   │   ├── deliveryservice.go  # Основная логика
│   │   │   ├── mapper.go           # Маппинг данных
│   │   │   ├── dellin_builder.go   # Построение запросов к Деловым Линиям
│   │   │   └── extra_services_mapper.go # Маппинг доп. услуг
│   │   └── service.go    # Агрегатор сервисов
│   │
│   ├── transport/        # Слой транспорта (HTTP)
│   │   ├── dto/          # Data Transfer Objects
│   │   │   ├── request/  # Структуры запросов
│   │   │   └── response/ # Структуры ответов
│   │   ├── handler/      # HTTP обработчики (Controllers)
│   │   │   ├── authHandler.go
│   │   │   ├── userHandler.go
│   │   │   └── handler.go
│   │   ├── httpserver/   # Настройка HTTP сервера (Echo)
│   │   │   └── server.go
│   │   └── middleware/   # Middleware (JWT, CORS и т.д.)
│   │       └── authMiddleware.go
│   │
│   └── utils/            # Вспомогательные утилиты
│       ├── cdek/         # Клиент для API СДЭК
│       ├── jwt/          # Работа с JWT токенами
│       ├── logger/       # Настройка логирования
│       └── password/     # Хеширование паролей
│
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile            # Docker образ приложения
└── go.mod                # Go зависимости
```

### Слои архитектуры

1. **Transport Layer** (`internal/transport/`)
   - HTTP handlers (Echo framework)
   - DTO (Data Transfer Objects)
   - Middleware (JWT, CORS, логирование)

2. **Service Layer** (`internal/service/`)
   - Бизнес-логика приложения
   - Интеграция с внешними API (СДЭК, Деловые Линии)
   - Маппинг и трансформация данных

3. **Repository Layer** (`internal/repository/`)
   - Доступ к базе данных (PostgreSQL через GORM)
   - Абстракция над БД

4. **Domain Layer** (`internal/domain/`)
   - Доменные модели и бизнес-сущности
   - Независимы от инфраструктуры

### Поток данных

```
HTTP Request
    ↓
Handler (transport/handler)
    ↓
Service (internal/service)
    ↓
Repository (internal/repository) ←→ PostgreSQL
    ↓
External APIs (СДЭК, Деловые Линии)
    ↓
Response
```

### Особенности реализации

- **Единый расчет доставки:** Сервис `deliveryservice` параллельно запрашивает варианты у всех провайдеров и объединяет результаты
- **Фильтрация:** Варианты фильтруются по значению `to` (откуда всегда "склад" по умолчанию) и исключаются тарифы с "постомат"
- **Сортировка:** Все варианты сортируются по цене независимо от провайдера
- **Маппинг:** Единый формат `DeliveryOption` для всех провайдеров

## 🔧 Конфигурация

Основные настройки находятся в `config/config.yaml`:

```yaml
env: local                    # Окружение (local, dev, prod)
jwt_secret: andrey            # Секретный ключ для JWT
http_server:
  port: :8080                 # Порт сервера
database:
  dsn: "host=db..."          # Строка подключения к БД
cdek:
  client_id: ...              # ID клиента СДЭК
  client_secret: ...          # Секрет СДЭК
dellin:
  app_key: ...                # API ключ Деловых Линий
```

## 📚 Дополнительная документация

- **Swagger UI:** http://localhost:8080/swagger/index.html
- Полная документация API доступна в Swagger UI после запуска проекта

## 🛠️ Разработка

### Локальная разработка (без Docker)

1. Установите Go 1.25.1+
2. Установите PostgreSQL
3. Настройте `config/config.yaml`
4. Запустите:
   ```bash
   go run cmd/api/main.go
   ```

### Генерация Swagger документации

```bash
# Установите swag
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерируйте документацию
swag init -g cmd/api/main.go -o docs
```

## 📝 Примечания

- Пароли пользователей хешируются с использованием bcrypt
- JWT токены используются для аутентификации защищенных endpoints
- Все цены возвращаются в копейках (int64)
- Даты доставки возвращаются в формате ISO 8601

## 🔒 Безопасность

- Пароли хранятся в хешированном виде
- JWT токены используются для аутентификации
- CORS настроен для разрешенных источников
- Валидация входных данных на всех уровнях

## 📄 Лицензия

[Укажите лицензию проекта]

