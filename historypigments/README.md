# HistoryPigments

## Что было сделано

### 1) Создан мобильный mockup по теме "История и культура №23 — Реконструкция первоначального цвета пигмента"

Реализованы три страницы:
- Лента
- Добавление
- Плитка

Файлы:
- `index.html`
- `add.html`
- `tiles.html`
- `assets/css/styles.css`

Цветовая система: светлый нейтральный фон, тёмный текст, акценты в цветах исторических пигментов.
Стиль: музейный/каталожный, вдохновлённый Classic Pigments.

### 2) Настроен MinIO для хранения изображений и видео

Минимальная локальная инфраструктура S3-хранилища:
- `docker-compose.yml`
- `.env`
- `scripts/init-minio.sh`

Параметры MinIO:
- console: http://localhost:9001
- API: http://localhost:9000
- login: root
- password: rootpassword

Используется публичный bucket:
- `historypigments`

### 3) Подготовлен Go-сервер с коллекцией пигментов

Файлы:
- `go.mod`
- `main.go`
- `templates/lenta.html`
- `templates/add.html`
- `templates/tiles.html`

Реализованы маршруты:
- `/lenta?id=1` — страница ленты по ID
- `/lenta?id=1&next=true` — следующая услуга/пигмент
- `/add` — черновик
- `/tiles?wavelength=610` — список карточек с фильтром по длине волны

### 4) Реализована модель данных

```go
type Pigment struct {
    ID          int
    Name        string
    Description string
    Color       string
    Wavelength  float64
    Reflection  float64
    Absorption  float64
    ImageKey    string
    VideoKey    string
    Status      string
    Likes       int
}
```

Примеры данных:
- Киноварь
- Жёлтая охра
- Малахит
- Лазурит
- Ультрамарин
- Умбра

Черновик хранится отдельно и показывается только на странице добавления.
Публикации отображаются на плитке и в ленте.

## Как запустить

### MinIO

```bash
cd "/Users/min/Desktop/rip/historypigments"
docker compose up -d
```

Проверка:
```bash
curl -I http://localhost:9000/minio/health/live
```

Console:
- http://localhost:9001
- login: root
- password: rootpassword

### Go-приложение

```bash
cd "/Users/min/Desktop/rip/historypigments"
go run .
```

Проверка:
```bash
curl -I http://localhost:8080/lenta?id=1
curl -I http://localhost:8080/add
curl -I "http://localhost:8080/tiles?wavelength=610"
```

## Что важно

- MinIO — это storage, не отдельная страница UI.
- image_key и video_key хранятся как отдельные поля модели.
- Для UI достаточно показывать блоки “Изображение” и “Видео”.
- В этом проекте нет JavaScript, только HTML/CSS/Go templates.
- Все данные в этом варианте являются коллекцией/массивом в Go, без БД.

## Примечание

В процессе работы проект несколько раз переименовывался (`lab1-service`, `historypigments`, затем в новый путь `/Users/min/Desktop/rip/historypigments`). Все команды и пути в текущей версии должны быть запущены из финального каталога:

```bash
cd "/Users/min/Desktop/rip/historypigments"
```

## Минимальные команды для работы

```bash
docker compose up -d
go run .
```

## Полезные ссылки

- MinIO API: http://localhost:9000
- MinIO Console: http://localhost:9001
- Go app: http://localhost:8080

## Итог

Создана рабочая базовая структура проекта по лабораторной работе:
- статический интерфейс/макет
- объектное хранилище MinIO
- Go-сервер с моделью пигментов, маршрутизацией и шаблонами
- тема исторических пигментов, как это требуется по методичке.
