#!/bin/bash

# Interactive setup script for TaskBot
# Creates .env file with user input and default values

echo "TaskBot Setup Wizard"
echo "======================"
echo ""

# Function to read user input with default value
read_input() {
    local prompt="$1"
    local default="$2"
    local var_name="$3"
    
    if [ -n "$default" ]; then
        read -p "$prompt [$default]: " input
        input="${input:-$default}"
    else
        while [ -z "$input" ]; do
            read -p "$prompt: " input
            if [ -z "$input" ]; then
                echo "This field is required!"
            fi
        done
    fi
    
    eval "$var_name=\"$input\""
}

# Function to validate port
validate_port() {
    local port="$1"
    if [[ ! "$port" =~ ^[0-9]+$ ]] || [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
        echo "Invalid port number. Please enter a number between 1 and 65535."
        return 1
    fi
    return 0
}

# Function to validate UID/GID
validate_uid_gid() {
    local value="$1"
    if [[ ! "$value" =~ ^[0-9]+$ ]] || [ "$value" -lt 1 ]; then
        echo "Invalid UID/GID. Please enter a positive number."
        return 1
    fi
    return 0
}

# Function to validate bot token format
validate_bot_token() {
    local token="$1"
    if [[ ! "$token" =~ ^[0-9]+:[a-zA-Z0-9_-]+$ ]]; then
        echo "❌ Invalid bot token format. Should be like: 123456789:ABCdefGHIjklMNOpqrSTUVwxyZ"
        return 1
    fi
    return 0
}

# Function to validate URL format
validate_url() {
    local url="$1"
    if [[ ! "$url" =~ ^https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/[^[:space:]]*)?$ ]]; then
        echo "Invalid URL format. Should start with http:// or https://"
        return 1
    fi
    return 0
}

echo "Telegram Bot Configuration (Both fields required)"
echo "--------------------------------------------------"

# Bot Token - ask first and validate format
while true; do
    read_input "Telegram Bot Token (get from @BotFather)" "" "BOT_TOKEN"
    if validate_bot_token "$BOT_TOKEN"; then
        break
    fi
done

# Bot Webhook - required with validation
while true; do
    read -p "Telegram Bot Webhook URL (e.g., https://yourdomain.com/bot-webhook): " BOT_WEBHOOK
    if [ -z "$BOT_WEBHOOK" ]; then
        echo "Webhook URL is required!"
    elif ! validate_url "$BOT_WEBHOOK"; then
        # Validation failed, continue loop
        continue
    else
        break
    fi
done

echo ""
echo "📋 System Configuration"
echo "--------------------"

# User UID/GID
read_input "User UID (for container permissions)" "1000" "USER_UID"
while ! validate_uid_gid "$USER_UID"; do
    read_input "User UID (for container permissions)" "" "USER_UID"
done

read_input "User GID (for container permissions)" "1000" "USER_GID"
while ! validate_uid_gid "$USER_GID"; do
    read_input "User GID (for container permissions)" "" "USER_GID"
done

# App port
read_input "Web server port" "8000" "APP_PORT"
while ! validate_port "$APP_PORT"; do
    read_input "Web server port" "" "APP_PORT"
done

echo ""
echo "Database Configuration (PostgreSQL)"
echo "------------------------------------"

read_input "PostgreSQL user" "user" "PG_USER"
read_input "PostgreSQL password" "p4sswordgh12345" "PG_PASSWORD"
read_input "PostgreSQL host" "postgres" "PG_HOST"
read_input "PostgreSQL port" "5436" "PG_PORT"
while ! validate_port "$PG_PORT"; do
    read_input "PostgreSQL port" "" "PG_PORT"
done

read_input "Database name" "bottask" "DB_NAME"

echo ""
echo "Redis Configuration"
echo "-------------------"

read_input "Redis host" "redis" "REDIS_HOST"
read_input "Redis port" "6379" "REDIS_PORT"
while ! validate_port "$REDIS_PORT"; do
    read_input "Redis port" "" "REDIS_PORT"
done

read_input "Redis password" "redispass" "REDIS_PASSWORD"

echo ""
echo "Logging Configuration"
echo "----------------------"

read_input "Log directory (local path)" "./logs" "LOG_DIR"
read_input "Log file name" "logs.log" "LOG_FILE"

echo ""
echo "Configuration Summary"
echo "---------------------"
echo "Telegram Bot:"
echo "  Token: ${BOT_TOKEN:0:10}..."
echo "  Webhook: ${BOT_WEBHOOK:-Not set}"
echo ""
echo "System:"
echo "  UID: $USER_UID, GID: $USER_GID"
echo "  App Port: $APP_PORT"
echo ""
echo "PostgreSQL:"
echo "  User: $PG_USER"
echo "  Host: $PG_HOST:$PG_PORT"
echo "  Database: $DB_NAME"
echo ""
echo "Redis:"
echo "  Host: $REDIS_HOST:$REDIS_PORT"
echo ""
echo "Logging:"
echo "  Directory: $LOG_DIR"
echo "  File: $LOG_FILE"
echo ""

# Confirm before writing
read -p "Do you want to create .env with this configuration? [Y/n] " confirm
confirm="${confirm:-Y}"

if [[ "$confirm" =~ ^[Yy]$ ]]; then
    echo "Creating .env file..."
    
    # Create .env file
    cat > .env << EOF
USER_UID=$USER_UID
USER_GID=$USER_GID

APP_PORT=$APP_PORT

PG_USER=$PG_USER
PG_PASSWORD=$PG_PASSWORD
PG_PORT=$PG_PORT
PG_HOST=$PG_HOST
DB_NAME=$DB_NAME

REDIS_HOST=$REDIS_HOST
REDIS_PORT=$REDIS_PORT
REDIS_PASSWORD=$REDIS_PASSWORD

BOT_TOKEN=$BOT_TOKEN
BOT_WEBHOOK=$BOT_WEBHOOK

LOG_DIR=$LOG_DIR
LOG_FILE=$LOG_FILE
EOF
    
    echo ".env file created successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Review .env file: cat .env"
    echo "2. Build and start: make up"
    echo "3. Run migrations: make migrate-up"
    echo ""
else
    echo "Setup cancelled. No changes were made."
    exit 1
fi