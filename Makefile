YQ := yq

DB_HOST := $(shell $(YQ) '.database.host' config.yaml)
DB_PORT := $(shell $(YQ) '.database.port' config.yaml)
DB_USER := $(shell $(YQ) '.database.user' config.yaml)
DB_PASS := $(shell $(YQ) '.database.password' config.yaml)
DB_NAME := $(shell $(YQ) '.database.dbname' config.yaml)

migrate-up:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/migrations up
migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path internal/migrations down
