# URL Shortener

HTTP-сервис для сокращения URL на Go.

Сервис позволяет создать короткую ссылку для исходного URL и получить исходный URL по короткой ссылке.

## Возможности

- создание короткой ссылки;
- получение исходного URL по короткой ссылке;
- одинаковый исходный URL всегда получает одну и ту же короткую ссылку;
- длина короткой ссылки — 10 символов;
- используются латинские буквы, цифры и `_`;
- два варианта хранения данных:
  - PostgreSQL;
  - in-memory storage;

- поддержка конкурентного доступа;
- unit-тесты;
- запуск в Docker.

## Стек

- Go
- PostgreSQL
- Docker / Docker Compose
- pgx
- net/http

## Архитектура

Сервис разделён на несколько слоёв:

Handler отвечает за HTTP-запросы и ответы.

Service содержит основную бизнес-логику.

Storage предоставляет единый интерфейс для работы с данными. Благодаря этому можно использовать PostgreSQL или in-memory хранилище без изменения бизнес-логики.

## Запуск через Docker

Для запуска нужен Docker.

Запустить сервис:

```bash
docker compose up -d
```

После запуска:

- API доступен на `http://localhost:8080`
- PostgreSQL доступен с хоста на порту `5433`

Миграция базы данных выполняется автоматически при запуске.

Проверить контейнеры:

```bash
docker ps
```

Остановить сервис:

```bash
docker compose down
```

Для полной очистки базы данных:

```bash
docker compose down -v
```

## Конфигурация

Используются переменные окружения:

```text
STORAGE_TYPE
DATABASE_URL
```

По умолчанию используется in-memory storage:

```text
STORAGE_TYPE=memory
```

Для PostgreSQL:

```text
STORAGE_TYPE=postgres
DATABASE_URL=postgres://shortener:shortener@postgres:5432/shortener
```

В Docker Compose эти значения уже настроены.

## API

### Создание короткой ссылки

```http
POST /shorten
Content-Type: application/json
```

Тело запроса:

```json
{
	"url": "https://example.com"
}
```

Пример через `curl`:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

Ответ:

```json
{
	"short_url": "JZD48piQv4"
}
```

### Получение исходного URL

```http
GET /{short_url}
```

Например:

```bash
curl http://localhost:8080/JZD48piQv4
```

Ответ:

```json
{
	"url": "https://example.com"
}
```

Если короткая ссылка не существует, сервис возвращает `404 Not Found`.

## Генерация коротких ссылок

Короткая ссылка имеет длину 10 символов.

Используется следующий алфавит:

```text
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_
```

Для генерации используется `crypto/rand`.

Пространство возможных ссылок:

```text
63^10 ≈ 9.85 × 10^17
```

Уникальность обеспечивается не только генерацией случайной строки, но и проверкой ограничения `UNIQUE` в хранилище.

При коллизии короткой ссылки сервис генерирует новую.

## Одинаковые URL

При повторном добавлении одного и того же URL существующая короткая ссылка возвращается повторно.

Например:

```text
POST https://example.com
→ JZD48piQv4

POST https://example.com
→ JZD48piQv4
```

Это поведение обеспечивается уникальностью исходного URL в хранилище.

## Конкурентный доступ

Сервис рассчитан на одновременную обработку нескольких запросов.

Для in-memory storage используется `sync.RWMutex`.

В PostgreSQL уникальность исходного и короткого URL дополнительно обеспечивается ограничениями базы данных.

Также есть тест конкурентного создания одной и той же ссылки.

## Тесты

Запустить все тесты:

```bash
go test ./...
```

Проверить код с помощью `go vet`:

```bash
go vet ./...
```

Для PostgreSQL-тестов PostgreSQL должен быть запущен.

## Запуск без Docker

Для запуска с in-memory storage достаточно:

```bash
go run ./cmd/server
```

После запуска сервер будет доступен на:

```text
http://localhost:8080
```

Для PostgreSQL необходимо задать:

```text
STORAGE_TYPE=postgres
DATABASE_URL=postgres://shortener:shortener@localhost:5433/shortener
```

и предварительно создать таблицу с помощью миграции:

```text
migrations/001_init.sql
```

## Завершение работы

Сервис поддерживает graceful shutdown. При получении сигнала завершения сервер прекращает принимать новые соединения и корректно завершает текущие запросы.

При использовании PostgreSQL также закрывается connection pool.
