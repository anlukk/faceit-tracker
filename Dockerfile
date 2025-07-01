FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest

COPY --from=builder /app/locales /app/locales
COPY --from=builder /app/main /app/main

WORKDIR /app
EXPOSE 8080
CMD ["./main"]