# go-messenger-microservices

Микросервисная архитектура для мессенджера на Go.

## Архитектура

```
Client → Nginx (80) → API Gateway (8080) → Микросервисы
```

### Сервисы

- **API Gateway** - Единая точка входа для всех запросов
- **Auth Service** - Аутентификация и авторизация (gRPC: 50051)
- **Chat Service** - Управление чатами (gRPC: 50052)
- **User Service** - Управление пользователями (gRPC: 50053, HTTP: 8090)

### Infrastructure

- **Nginx** - Reverse proxy и load balancer
- **PostgreSQL** - Базы данных для каждого сервиса
- **Redis** - Кэширование и сессии
- **Kafka** - Message broker для асинхронной обработки
- **Prometheus** - Мониторинг метрик
- **Grafana** - Визуализация метрик
- **Jaeger** - Distributed tracing

## Запуск проекта

### Все сервисы (Docker Compose)

```bash
docker-compose up -d
```

### Запуск отдельных сервисов

```bash
# Auth Service
cd auth && make run

# Chat Service
cd chat && make run

# User Service
cd user_service && make run

# API Gateway
cd api-gateway && make run
```

## Endpoints

### Nginx (порт 80)
- `GET /` - Главная страница
- `GET /api/v1/*` - API endpoints
- `GET /health` - Health check

### API Gateway (порт 8080)
- `GET /health` - Проверка работоспособности
- `GET /status` - Статус всех сервисов
- `/api/v1/users/*` - User Service endpoints

### User Service
- `GET http://localhost:8090/swagger/` - Swagger UI
- `POST http://localhost:8090/api/v1/users` - Создать пользователя
- `GET http://localhost:8090/api/v1/users/{id}` - Получить пользователя

## Структура проекта

```
├── api-gateway/       # API Gateway сервис
├── auth/              # Сервис аутентификации
├── chat/              # Сервис чатов
├── user_service/      # Сервис пользователей
├── nginx/             # Конфигурация Nginx
└── docker-compose.yml # Общая конфигурация
```

## Переменные окружения

Создайте `.env` файл в корне проекта:

```env
# Database
PG_DATABASE_NAME=messenger_db
PG_USER=messenger_user
PG_PASSWORD=messenger_password
PG_PORT=5432

# Auth Service
AUTH_PG_DB=auth_db
AUTH_PG_USER=auth_user
AUTH_PG_PASSWORD=auth_password

# Chat Service
CHAT_PG_DB=chat_db
CHAT_PG_USER=chat_user
CHAT_PG_PASSWORD=chat_password

# User Service
USER_PG_DB=user_db
USER_PG_USER=user_user
USER_PG_PASSWORD=user_password

# Redis
REDIS_PORT=6379
AUTH_REDIS_PASSWORD=auth_redis_password
USER_REDIS_PASSWORD=user_redis_password

# JWT
JWT_SECRET=your-secret-key-here
```

## Разработка

```bash
# Запуск всех сервисов
make up

# Остановка всех сервисов
make down

# Логи
docker-compose logs -f

# Пересборка
docker-compose up -d --build
```

## Технологии

- **Go** - Основной язык
- **gRPC** - Межсервисное взаимодействие
- **Docker** - Контейнеризация
- **PostgreSQL** - База данных
- **Redis** - Кэширование
- **Kafka** - Event streaming
- **Nginx** - Reverse proxy
- **Prometheus** - Мониторинг
- **Jaeger** - Tracing


