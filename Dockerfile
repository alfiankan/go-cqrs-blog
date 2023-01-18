FROM golang:1.19-alpine as builder
WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app-api ./cmd/api/main.go
RUN go build -o app-cli ./cmd/cli/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app-api .
COPY --from=builder /app/app-cli .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/articles_seed.json ./articles_seed.json



ENV HTTP_API_PORT=3000

EXPOSE 3000

ENTRYPOINT ["/bin/sh", "-c", "/app/app-cli migrate up && /app/app-api"]

