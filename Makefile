export $(shell grep -v '^#' .env | xargs)

create-database:
	@echo "Checking if database $(POSTGRES_DBNAME) exists..."
	@psql postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT) -tc "SELECT 1 FROM pg_database WHERE datname = '$(POSTGRES_DBNAME)'" | grep -q 1 || \
		psql postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT) -c "CREATE DATABASE $(POSTGRES_DBNAME);" && \
		echo "Database $(POSTGRES_DBNAME) created successfully." || \
		echo "Database $(POSTGRES_DBNAME) already exists."

run:
	HTTP_SERVER_HOST=localhost HTTP_SERVER_PORT=8080 POSTGRES_HOST=localhost POSTGRES_PORT=5432 POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DBNAME=avito REDIS_HOST=localhost REDIS_PORT=6379 CACHE_PASSWORD=your_redis_password LOG_LEVEL=info JWT_DURATION=24h JWT_SECRET=lol go run cmd/main/main.go

migrate:
	HTTP_SERVER_HOST=localhost HTTP_SERVER_PORT=8080 POSTGRES_HOST=localhost POSTGRES_PORT=5432 POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DBNAME=avito REDIS_HOST=localhost REDIS_PORT=6379 CACHE_PASSWORD=your_redis_password LOG_LEVEL=info JWT_DURATION=24h JWT_SECRET=lol go run cmd/migrator/main.go

migrate-up: create-database
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DBNAME)?sslmode=disable" -path internal/migrations up

migrate-down:
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DBNAME)?sslmode=disable" -path internal/migrations down

run-server:
	echo "Starting server on $(HTTP_SERVER_HOST):$(HTTP_SERVER_PORT)"

print-config:
	@echo "HTTP Server: $(HTTP_SERVER_HOST):$(HTTP_SERVER_PORT)"
	@echo "Database: $(POSTGRES_USER)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DBNAME)"
	@echo "Redis: $(REDIS_HOST):$(REDIS_PORT)"
	@echo "Log Level: $(LOG_LEVEL)"
	@echo "JWT Duration: $(JWT_DURATION)"
