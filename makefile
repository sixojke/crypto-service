build:
	go build -o app ./cmd/crypto-service

up: build
	sudo docker-compose up -d --build

down:
	sudo docker-compose down

restart: down build up

swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/crypto-service/main.go