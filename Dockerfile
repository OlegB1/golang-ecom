FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o ./bin/main ./cmd

FROM gcr.io/distroless/base

WORKDIR /root/

COPY --from=builder /app/bin/main .

COPY .env .

EXPOSE 8080

CMD ["./main"]
