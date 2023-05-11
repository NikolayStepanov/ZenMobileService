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
test:
	go test -v -count=1 ./...
test100:
	go test -v -count=100 ./...
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
gen:
	mockgen -source=internal/service/service.go \
	-destination=internal/service/mocks/mock_service.go
