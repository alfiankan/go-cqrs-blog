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


ENV HTTP_API_PORT=3000

EXPOSE 3000

ENTRYPOINT ["/app/app-api"]

