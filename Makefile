include .env

swag:
	swag init -g cmd/server/main.go -o docs
