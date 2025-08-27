SHELL := /bin/bash
ENV_FILE := .env

# Функция для чтения переменной из .env, если есть, иначе пусто
define get_env_var
$(shell grep -v '^#' $(ENV_FILE) | grep "^$(1)=" | cut -d '=' -f2-)
endef

# Читаем переменные из .env или ставим дефолты
DATABASE_URL := $(call get_env_var,DATABASE_URL)

DB_HOST := $(call get_env_var,DB_HOST)
DB_PORT := $(call get_env_var,DB_PORT)
DB_USER := $(call get_env_var,DB_USER)
DB_PASSWORD := $(call get_env_var,DB_PASSWORD)
DB_NAME := $(call get_env_var,DB_NAME)


# Формируем DATABASE_URL, если она не задана
ifeq ($(DATABASE_URL),)
	DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
endif

.PHONY: migrate-up migrate-down

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down
