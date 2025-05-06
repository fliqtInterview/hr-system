.PHONY: build run test docker-build docker-up docker-down

build:
	go build -o main .

run:
	go run main.go

test:
	go test -v ./...

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -f main 