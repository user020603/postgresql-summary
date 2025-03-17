FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]