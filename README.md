# MyGoMessenger

Простое приложение для обмена сообщениями, созданное с использованием Go и микросервисов с поддержкой как gRPC, так и REST API интерфейсов.

## Технологии

- **Go 1.25.2**: Основной язык программирования.
- **gRPC**: Для синхронного взаимодействия между сервисами.
- **Protocol Buffers**: Для определения контрактов API.
- **Kafka**: Для асинхронной обработки сообщений (уведомлений).
- **PostgreSQL 15**: В качестве основной базы данных.
- **Docker & Docker Compose**: Для контейнеризации и оркестрации приложения.
- **OpenAPI 3.0**: Спецификация REST API.

## Архитектура

Приложение состоит из трех микросервисов, организованных как отдельные Go-модули:

- **Сервис пользователей** (`users/`): Управляет регистрацией пользователей и профилем.
- **Сервис сообщений** (`messages/`): Обрабатывает отправку и получение сообщений.
- **Сервис уведомлений** (`notifications/`): Прослушивает новые сообщения в Kafka и логирует их.

Сервисы взаимодействуют друг с другом с помощью gRPC для синхронных вызовов и Kafka для асинхронных. Данные хранятся в базе данных PostgreSQL с автоматическим применением миграций.

## Структура проекта

```
├── client/                 # Простой тестовый клиент
├── test_client/            # Расширенный тестовый клиент (отдельный модуль)
├── users/                  # Сервис пользователей
├── messages/               # Сервис сообщений
├── notifications/          # Сервис уведомлений
├── proto/                  # Protocol Buffers определения
├── gen/                    # Сгенерированные gRPC файлы
├── db/migrations/          # Миграции базы данных
└── github.com/leshkoan/MyGoMessenger/  # Локальные модули
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

- **Сервис пользователей**: http://localhost:8081
- **Сервис сообщений**: http://localhost:8082
- **База данных PostgreSQL**: localhost:5432
- **Kafka**: localhost:9092

### 4. Проверьте работу сервисов:

```bash
# Проверка таблиц в базе данных
docker-compose exec postgres psql -U user -d mydatabase -c "SELECT tablename FROM pg_tables;"
```

## Тестирование

### Простой тестовый клиент

Базовый клиент демонстрирует работу с gRPC-сервисами:

```bash
go run client/main.go
```

Клиент последовательно:

1. Создает двух пользователей
2. Отправляет несколько сообщений между ними

### Расширенный тестовый клиент

Отдельный модуль с дополнительными возможностями тестирования:

```bash
cd test_client
go mod tidy
go run test_messenger.go
```

### Автоматизированные тесты

```bash
go test ./...
```

## API интерфейсы

### gRPC API

Сервисы предоставляют gRPC интерфейсы, определенные в `proto/*.proto` файлах.

**Генерация gRPC кода:**

```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### REST API

Проект включает OpenAPI спецификацию (`openapi.yaml`) с поддержкой HTTP endpoints:

#### Пользователи

- `POST /users/register` - Регистрация нового пользователя

#### Сообщения

- `POST /messages/send` - Отправить сообщение
- `GET /messages/history?user1={id}&user2={id}` - Получить историю сообщений

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

## Отладка

Для локальной разработки сервисов без Docker:

```bash
# Сервис пользователей
cd users && go run main.go server.go

# Сервис сообщений
cd messages && go run main.go server.go

# Сервис уведомлений
cd notifications && go run main.go
```

Убедитесь, что база данных и Kafka запущены через Docker Compose.
