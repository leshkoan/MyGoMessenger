# MyGoMessendger

Простое приложение для обмена сообщениями, созданное с использованием Go и микросервисов.

## Архитектура

Приложение состоит из трех микросервисов:
- **Сервис пользователей**: Управляет регистрацией и аутентификацией пользователей.
- **Сервис сообщений**: Обрабатывает отправку и получение сообщений.
- **Сервис уведомлений**: Прослушивает новые сообщения и логирует их.

Сервисы взаимодействуют друг с другом с помощью базы данных PostgreSQL и брокера сообщений Kafka.

## Предварительные требования

- Docker
- Docker Compose

## Начало работы

1. Клонируйте репозиторий.
2. Создайте файл `.env` в корне проекта со следующим содержимым:
   ```
   POSTGRES_USER=user
   POSTGRES_PASSWORD=password
   POSTGRES_DB=mydatabase
   ```
3. Выполните следующую команду, чтобы запустить приложение:
   ```
   docker-compose up --build -d
   ```

## Использование API

### Регистрация нового пользователя

```
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser"}' http://localhost:8081/users/register
```

### Отправка сообщения

```
curl -X POST -H "Content-Type: application/json" -d '{"from_user_id":"<USER_ID_1>","to_user_id":"<USER_ID_2>","text":"Hello, world!"}' http://localhost:8082/messages/send
```

### Получение истории сообщений

```
curl -X GET "http://localhost:8082/messages/history?user1=<USER_ID_1>&user2=<USER_ID_2>"
```