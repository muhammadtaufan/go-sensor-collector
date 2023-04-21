setup:
	go mod tidy
	cp .env.sample .env
build:
	go build -o bin/go-sensor-collector cmd/main.go

run: build
	./bin/go-sensor-collector

# migration cmd
MYSQL_DSN=root:@tcp(localhost:3306)/
DATABASE_NAME=sensor_collector

migrate-up:
	@echo "Checking if database exists..."
	@mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS $(DATABASE_NAME);"
	@echo "Running database migrations..."
	@migrate -path migration -database "mysql://$(MYSQL_DSN)$(DATABASE_NAME)?parseTime=true" up

migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path migration -database "mysql://$(MYSQL_DSN)$(DATABASE_NAME)?parseTime=true" down

migrate-create:
	@echo "Creating new database migration..."
	@migrate create -ext sql -tz utc -dir migration $(name) -database "mysql://$(MYSQL_DSN)$(DATABASE_NAME)?parseTime=true"

.PHONY: proto
