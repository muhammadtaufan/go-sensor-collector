setup:
	go mod tidy
	cp .env.sample .env
build:
	go build -o bin/go-sensor-collector cmd/main.go

run: build
	./bin/go-sensor-collector

.PHONY: proto
