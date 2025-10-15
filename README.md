# MyGoMessenger

Простое приложение для обмена сообщениями, созданное с использованием Go и микросервисов на базе gRPC с асинхронной обработкой уведомлений через Apache Kafka.

## Технологии

- **Go 1.25.2**: Основной язык программирования.
- **gRPC**: Для синхронного взаимодействия между сервисами.
- **Protocol Buffers**: Для определения контрактов API.
- **Apache Kafka**: Для асинхронной обработки сообщений (уведомлений).
- **PostgreSQL 15**: В качестве основной базы данных.
- **Docker & Docker Compose**: Для контейнеризации и оркестрации приложения.

## Архитектура

Приложение состоит из трех микросервисов, организованных как отдельные Go-модули:

- **Сервис пользователей** (`users/`): Управляет регистрацией пользователей и профилем.
- **Сервис сообщений** (`messages/`): Обрабатывает отправку и получение сообщений.
- **Сервис уведомлений** (`notifications/`): Прослушивает новые сообщения в Kafka и логирует их.

Сервисы взаимодействуют друг с другом с помощью gRPC для синхронных вызовов и Kafka для асинхронной обработки уведомлений. Данные хранятся в базе данных PostgreSQL с автоматическим применением миграций.

## Структура проекта

```
├── test_client/            # Тестовый клиент для комплексного тестирования
├── client/                 # Простой клиент (опциональный)
├── users/                  # Сервис пользователей (gRPC порт 8081)
├── messages/               # Сервис сообщений (gRPC порт 8082)
├── notifications/          # Сервис уведомлений (Kafka consumer)
├── proto/                  # Protocol Buffers определения
├── gen/                    # Сгенерированные gRPC файлы
├── db/migrations/          # Миграции базы данных PostgreSQL
├── github.com/leshkoan/MyGoMessenger/  # Локальные Go модули
└── docker-compose.yml      # Оркестрация всех сервисов
```

## Предварительные требования

- **Docker** и **Docker Compose**
- **Go 1.25.2** (для запуска клиентов и генерации кода)
- **protoc** (компилятор Protocol Buffers)

## Начало работы

### 1. Клонируйте репозиторий

### 2. Создайте файл `.env` в корне проекта:

```env
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=mydatabase
```

### 3. Запустите все сервисы:

```bash
docker-compose up --build -d
```

После запуска сервисы будут доступны на следующих портах:

- **Сервис пользователей (gRPC)**: localhost:8081
- **Сервис сообщений (gRPC)**: localhost:8082
- **База данных PostgreSQL**: localhost:5432
- **Kafka**: localhost:9092

### 4. Проверьте работу сервисов:

```bash
# Проверка таблиц в базе данных
docker-compose exec postgres psql -U user -d mydatabase -c "SELECT tablename FROM pg_tables;"
```

## Тестирование

### Тестовый клиент

Отдельный модуль для комплексного тестирования мессенджера:

```bash
cd test_client
go mod tidy
go run test_messenger.go
```

Клиент автоматически:

1. Создает трех пользователей (user_a, user_b, user_c)
2. Организует полный обмен сообщениями между всеми пользователями
3. Проверяет корректность работы всех микросервисов

## API интерфейсы

### gRPC API

Сервисы предоставляют gRPC интерфейсы, определенные в `proto/*.proto` файлах.

**Генерация gRPC кода:**

```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

**Доступные сервисы:**

#### User Service (порт 8081)

- `RegisterUser(username)` - Регистрация нового пользователя

#### Message Service (порт 8082)

- `SendMessage(fromUserId, toUserId, text)` - Отправить сообщение

## Разработка

### Структура модулей

Проект использует Go-модули с локальными заменами:

```go
// go.mod
replace github.com/leshkoan/MyGoMessenger/gen/go/users => ./gen/go/users
replace github.com/leshkoan/MyGoMessenger/gen/go/messages => ./gen/go/messages
```

### Миграции базы данных

Миграции автоматически применяются при запуске PostgreSQL контейнера из директории `db/migrations/`.

### Логирование

Сервисы предоставляют детальное логирование своей работы. Для просмотра логов:

```bash
docker-compose logs user-service
docker-compose logs message-service
docker-compose logs notification-service
```

## Мониторинг и отладка

### Проверка логов сервисов:

```bash
# Логи всех сервисов
docker-compose logs

# Логи конкретного сервиса
docker-compose logs user-service
docker-compose logs message-service
docker-compose logs notification-service
```

### Проверка базы данных:

```bash
# Список таблиц
docker-compose exec postgres psql -h localhost -U user -d mydatabase -c "\dt"

# Просмотр пользователей
docker-compose exec postgres psql -h localhost -U user -d mydatabase -c "SELECT * FROM users;"

# Просмотр сообщений
docker-compose exec postgres psql -h localhost -U user -d mydatabase -c "SELECT * FROM messages;"
```

## Результаты тестирования

Проект успешно протестирован с использованием расширенного тестового клиента:

✅ **Микросервисы**: Все сервисы (users, messages, notifications) работают корректно
✅ **База данных**: Пользователи и сообщения сохраняются в PostgreSQL
✅ **Kafka интеграция**: Асинхронная обработка уведомлений функционирует
✅ **gRPC коммуникация**: Межсервисное взаимодействие работает стабильно

**Тестовый сценарий**: Создание 3 пользователей и полный обмен сообщениями между всеми парами (6 сообщений) прошел успешно.
