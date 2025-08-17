include .env.local

DB_URL=mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)
MIGRATE_PATH=database/migration

build:
	@go build -o bin/MVC cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/MVC
	
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

docker-up:
	docker compose --env-file .env.docker up --build

docker-down:
	docker compose down --volumes