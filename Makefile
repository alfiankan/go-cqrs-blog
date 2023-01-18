migrate:
	go run ./cmd/cli/... migrate up

migrate-down:
	go run ./cmd/cli/... migrate down

seed:
	go run ./cmd/cli/... seed

flush-index:
	curl --user elastic:elastic -X DELETE "localhost:9200/articles"
