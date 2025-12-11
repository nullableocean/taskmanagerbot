# TaskBot - Менеджер задач с управлением через Telegram

TaskBot - это приложение для управления задачами через Telegram-бота с поддержкой базы данных PostgreSQL и кэширования Redis.

## Конфигурация

### 1. Настройка окружения


```bash
cp .env.example .env
```

В конфиге `.env`:

- `BOT_TOKEN` - токен бота (в BotFather)
- `BOT_WEBHOOK` - HTTPS вебхук для получения сообщений от тг
- `PG_USER`, `PG_PASSWORD`, `PG_HOST`, `PG_PORT`, `DB_NAME` - конфигурация PostgreSQL, хост по умолчанию из docker-compose
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` - конфигурация Redis, хост по умолчанию из docker-compose
- `USER_UID`, `USER_GID` - хостовые UID/GID пользователя, нужны для создания юзера в контейнере

### 2. Сборка и запуск

```bash
# Сборка и запуск контейнеров
make up

# Остановка контейнеров
make down

# Просмотр логов
make logs
```

### 3. Миграции базы данных

Используется утилита: github.com/golang-migrate/migrate


```bash
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Принудительно установить версию
make migrate-force VERSION=1

# Показать текущую версию
make migrate-version
```

### 4. Полная настройка

Команда `setup` интерактивно создаст `.env` (если не существует) и применит миграции:

```bash
make setup
```

