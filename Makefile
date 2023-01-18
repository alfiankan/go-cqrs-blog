migrate:
	go run ./cmd/cli/... migrate up

migrate-down:
	go run ./cmd/cli/... migrate down

seed:
	go run ./cmd/cli/... seed

test:
	go test ./article/tests/... -v

swagger:
	swag init -g cmd/api/main.go

docker:
	docker-compose up -d

run:
	go run ./cmd/api/main.go

