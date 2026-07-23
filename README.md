# Room Booker

Сервис бронирования переговорных комнат на Go и PostgreSQL.

## Возможности

* создание переговорных комнат;
* настройка расписания;
* просмотр свободных слотов;
* бронирование и отмена брони;
* регистрация и авторизация;
* роли `admin` и `user`;
* защита от двойного бронирования.

## Запуск

Требуются Docker и Docker Compose.

```bash
git clone https://github.com/NutcrackerCom/room-booker.git
cd room-booker
docker compose up --build
```

Сервис будет доступен по адресу:

```text
http://localhost:8080
```

Проверка:

```bash
curl http://localhost:8080/_info
```

## Авторизация

Получить тестовый токен администратора:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}'
```

Токен пользователя:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role":"user"}'
```

В защищённых запросах передавайте токен:

```http
Authorization: Bearer <token>
```

## Основные эндпоинты

| Метод  | Путь                              | Описание            |
| ------ | --------------------------------- | ------------------- |
| `POST` | `/register`                       | Регистрация         |
| `POST` | `/login`                          | Авторизация         |
| `GET`  | `/rooms/list`                     | Список комнат       |
| `POST` | `/rooms/create`                   | Создание комнаты    |
| `POST` | `/rooms/{roomId}/schedule/create` | Создание расписания |
| `GET`  | `/rooms/{roomId}/slots/list`      | Свободные слоты     |
| `POST` | `/bookings/create`                | Создание брони      |
| `GET`  | `/bookings/my`                    | Мои бронирования    |
| `POST` | `/bookings/{bookingId}/cancel`    | Отмена брони        |

Полное описание API находится в файле [`api.yaml`](api.yaml).

## Тесты

```bash
go test ./...
```

E2E-тесты:

```bash
docker compose up -d --build
go test ./tests/e2e -v
```

## Стек

* Go;
* PostgreSQL;
* JWT;
* Docker Compose;
* OpenAPI.
