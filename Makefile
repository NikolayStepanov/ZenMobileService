up:
	docker-compose up
stop:
	docker-compose stop
down:
	docker-compose down
build:
	docker-compose build
swag:
	swag init -g ./cmd/app/main.go