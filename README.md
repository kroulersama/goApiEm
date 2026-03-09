# Subscriptions API

API для управления подписками пользователей.

## Запуск

```
docker compose up --build -d
```

## Конфигурация (.env)

```
DB_HOST=db
DB_PORT=5432
DB_USER=demo
DB_PASSWORD=secret
DB_NAME=demo_db
DB_SSLMODE=disable
SERVER_PORT=8080
```


## Эндпоинты

| Метод  | Путь                | Описание                    |
|--------|---------------------|-----------------------------|
| POST   | /subs               | Создать подписку            |
| GET    | /subs               | Получить все подписки       |
| GET    | /subs/{id}          | Получить подписку по ID     |
| POST   | /subs/{id}          | Обновить подписку           |
| DELETE | /subs/{id}          | Удалить подписку            |
| GET    | /subs/{id}/prices   | Сумма за период             |

## Документация

Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/
```
## Назначение файлов

```
cmd      - вход
docs     - swagger
internal - вся логика
  - config     - логгер
  - handler    - обработка http
  - repository - обработка с базой данных
  - service    - бизнесс логика 
```
