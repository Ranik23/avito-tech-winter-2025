YQ := yq

HTTP_HOST := $(shell $(YQ) '.httpServerConfig.host' config.yaml)
HTTP_PORT := $(shell $(YQ) '.httpServerConfig.port' config.yaml)

DB_HOST := $(shell $(YQ) '.postgresConfig.host' config.yaml)
DB_PORT := $(shell $(YQ) '.postgresConfig.port' config.yaml)
DB_USER := $(shell $(YQ) '.postgresConfig.user' config.yaml)
DB_PASS := $(shell $(YQ) '.postgresConfig.password' config.yaml)
DB_NAME := $(shell $(YQ) '.postgresConfig.dbname' config.yaml)

REDIS_HOST := $(shell $(YQ) '.redisConfig.host' config.yaml)
REDIS_PORT := $(shell $(YQ) '.redisConfig.port' config.yaml)
REDIS_PASS := $(shell $(YQ) '.redisConfig.password' config.yaml)

LOG_LEVEL := $(shell $(YQ) '.logConfig.level' config.yaml)

JWT_DURATION := $(shell $(YQ) '.jwtConfig.duration' config.yaml)

create-database:
	@echo "Checking if database $(DB_NAME) exists..."
	@psql postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT) -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || \
		psql postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT) -c "CREATE DATABASE $(DB_NAME);" && \
		echo "Database $(DB_NAME) created successfully." || \
		echo "Database $(DB_NAME) already exists."

run:
	go run cmd/main/main.go

migrate-up: create-database
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/migrations up

migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/migrations down

run-server:
	echo "Starting server on $(HTTP_HOST):$(HTTP_PORT)"


print-config:
	@echo "HTTP Server: $(HTTP_HOST):$(HTTP_PORT)"
	@echo "Database: $(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"
	@echo "Redis: $(REDIS_HOST):$(REDIS_PORT)"
	@echo "Log Level: $(LOG_LEVEL)"
	@echo "JWT Duration: $(JWT_DURATION)"
