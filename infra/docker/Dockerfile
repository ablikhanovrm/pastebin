FROM golang:1.25-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/app
RUN go build -o migrate ./cmd/migrate

#runtime
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrate .

CMD ["./app"]