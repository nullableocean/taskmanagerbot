#!/bin/bash

# Load vars from .env file
if [ -f "./.env" ]; then
    set -a
    source ./.env
    set +a
else
    error_exit "Файл .env не найден. Пожалуйста, создайте его на основе .env.example"
fi

# POSTGRES DSN
PG_HOST="localhost"
DB_DSN="postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${DB_NAME}?sslmode=disable"
MIGRATIONS_PATH="./migrations"


error_exit() {
    echo -e "$Ошибка: $1" >&2
    exit 1
}

check_tool_installed() {
    if ! command -v migrate &> /dev/null; then
        echo -e "Утилита migrate не найдена. Установка"
        install_migrate
    fi
}

# install migrate util
install_migrate() {
    echo "Установка утилиты migrate..."
    
    # Try to install using go install
    if command -v go &> /dev/null; then
        echo "Установка через go install..."
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        
        # Add GOPATH/bin to PATH if not already there
        GOPATH_BIN="$(go env GOPATH)/bin"
        if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]]; then
            export PATH="$PATH:$GOPATH_BIN"
        fi
        
        if command -v migrate &> /dev/null; then
            echo "Утилита migrate успешно установлена"
            return 0
        fi
    fi
    
    echo "Попытка загрузки бинарного файла..."
    
    # Detect OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) ARCH="amd64" ;;
    esac
    
    # Download and install
    URL="https://github.com/golang-migrate/migrate/releases/latest/download/migrate.${OS}-${ARCH}.tar.gz"
    echo "Загрузка с $URL..."
    
    if command -v curl &> /dev/null; then
        curl -L $URL | tar xvz
    elif command -v wget &> /dev/null; then
        wget -qO- $URL | tar xvz
    else
        error_exit "Не удалось установить migrate. Установите Go или curl/wget."
    fi
    
    # Move to /usr/local/bin
    if [ -f ./migrate ]; then
        sudo mv ./migrate /usr/local/bin/migrate
        sudo chmod +x /usr/local/bin/migrate
        echo "Утилита migrate установлена в /usr/local/bin/migrate"
    else
        error_exit "Не удалось установить утилиту migrate"
    fi
}

check_empty_dsn() {
    if [ -z "$DB_DSN" ]; then
        error_exit "DSN базы данных не задан\n\n"
    fi
}

check_migrations_dir() {
    if [ ! -d "$MIGRATIONS_PATH" ]; then
        error_exit "Директория с миграциями не найдена: $MIGRATIONS_PATH"
    fi
}

run_migration() {
    local command=$1
    local arg=$2
    
    case $command in
        up)
            if [ -n "$arg" ]; then
                echo "UP $arg..."
                migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" up "$arg"
            else
                echo "UP..."
                migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" up
            fi
            ;;
        down)
            if [ -n "$arg" ]; then
                echo "DOWN $arg ..."
                migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" down "$arg"
            else
                echo "DOWN 1 ..."
                migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" down 1
            fi
            ;;
        force)
            if [ -n "$arg" ]; then
                echo "Принудительная установка версии $arg..."
                migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" force "$arg"
            else
                error_exit "Для команды force необходимо указать версию"
            fi
            ;;
        version)
            echo "Текущая версия миграций:"
            migrate -path "$MIGRATIONS_PATH" -database "$DB_DSN" version
            ;;
        *)
            error_exit "Неизвестная команда: $command. Команды: up down force version"
            ;;
    esac
}

main() {
    if [ $# -lt 1 ]; then
        echo "Использование: $0 {up|down|force|count} [count]"
        echo "Примеры:"
        echo "  $0 up                    # Применить все миграции"
        echo "  $0 up 3                  # Применить 3 миграции"
        echo "  $0 down                  # Откатить на 1 миграцию"
        echo "  $0 down 2                # Откатить на 2 миграции"
        echo "  $0 force 5               # Принудительно установить версию бд 5"
        echo "  $0 version               # Показать текущую версию"
        echo ""
        exit 1
    fi
    
    local command=$1
    local version=$2

    check_empty_dsn
    check_tool_installed
    check_migrations_dir
    
    run_migration "$command" "$version"
}

main "$@"