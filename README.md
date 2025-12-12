# Mood Tracker API
REST API для отслеживания настроения, реализованное на Go с использованием Gin.

## Функциональность
- Регистрация и аутентификация пользователей
- CRUD операции для записей настроения
- Добавление и обновление записей настроения по дате
- Просмотр всех записей настроения пользователя
- Генерация отчетов по настроению за месяц
- Отчеты по конкретным настроениям
- Смена активного пользователя

# Установка и запуск

1 Клонируйте репозиторий:
```bash
git clone https://github.com/ВАШ_USERNAME/mood-tracker-api.git
cd mood-tracker-api
```

2 Установите зависимости:
```bash
go mod download
```
3 
Запустите приложение:
```bash
go run main.go
```
# Использование API

- После запуска сервер будет доступен по адресу http://localhost:8080.

### Основные эндпоинты
* **Аутентификация**

     * POST /auth/register — Регистрация нового пользователя.

     * POST /auth/login — Вход в систему.

* **Записи настроения (Mood Entries)**

     * POST /api/entries — Добавить новую запись настроения.

     * GET /api/entries — Получить запись настроения по дате (параметры: month, day).

     * GET /api/entries/all — Получить все записи настроения текущего пользователя.

* **Отчеты (Reports)**

     * GET /api/reports/general — Общий отчет по настроениям за месяц (параметр: month).

     * GET /api/reports/mood — Отчет по конкретному настроению за месяц (параметры: month, mood).

* **Управление пользователями**

     * POST /api/switch-user — Сменить текущего пользователя.

     * GET /api/current-user — Получить информацию о текущем пользователе.

* **Информация**

     * GET / — Информация о доступных эндпоинтах.

     * GET /health — Статус сервера и статистика.

## Примеры запросов

### Регистрация нового пользователя
```bash
$body = @{
    username = "alex"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/auth/register" -Method Post -Body $body -ContentType "application/json"
  ```
### Вход пользователя
```bash
Invoke-RestMethod -Uri "http://localhost:8080/auth/login" -Method Post -Body $body -ContentType "application/json"
  ```
### Добавление записи настроения
```bash
$entryBody = @{
    day = 15
    month = 10
    mood = "Радостное"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/entries" -Method Post -Body $entryBody -ContentType "application/json"

  ```
### Получение всех записей
```bash
Invoke-RestMethod -Uri "http://localhost:8080/api/entries/all" -Method Get
```
### Отчёт по конкретной дате
```bash
Invoke-RestMethod -Uri "http://localhost:8080/api/entries?month=10&day=15"
```
### Общий отчет за месяц
```bash
Invoke-RestMethod -Uri "http://localhost:8080/api/reports/general?month=10" -Method Get
```
### Отчет по конкретному настроению
```bash
Invoke-RestMethod -Uri "http://localhost:8080/api/reports/mood?month=10&mood=Грустное" -Method Get
```
# Доступные настроения

API поддерживает следующие настроения:
- "Грустное"
-  "Спокойное"
- "Нейтральное"
- "Радостное"
- "Злое"

## Структура данных
### Запись настроения (MoodEntry)
```json
{
  "day": 15,
  "month": 10,
  "mood": "Радостное"
}
Пользователь (User)
json
{
  "username": "alex",
  "entries": [
    {
      "day": 15,
      "month": 10,
      "mood": "Радостное"
    }
  ]
}
```
## Хранение данных
### Данные сохраняются в файле mood_diary.json в формате JSON. При каждом изменении данных файл автоматически обновляется.

## Технологии
- Go — Язык программирования
- Gin — Веб-фреймворк
- JSON — Формат хранения данных
- Standard Library — Использование стандартных библиотек Go для работы с файлами
---
*Этот проект является учебным и создан для демонстрации построения REST API на Go.*