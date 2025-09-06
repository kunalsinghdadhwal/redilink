FROM golang:1.22-alpine AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init

RUN go build -o main .

# -------------------------
# Runtime image
# -------------------------
FROM alpine:latest

RUN adduser -S -D -H -h /app appuser
USER appuser

WORKDIR /app

COPY --from=builder /build/main .
COPY .env .
COPY docs ./docs

EXPOSE 3000

CMD ["./main"]
