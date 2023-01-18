migrate:
	go run ./cmd/cli/... migrate up

migrate-down:
	go run ./cmd/cli/... migrate down

seed:
	go run ./cmd/cli/... seed

test:
	go test ./article/tests/... -v


