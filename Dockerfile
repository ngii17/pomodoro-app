# Pomodoro App
FROM golang:1.21-alpine

WORKDIR /app

COPY pomodoro/go.mod pomodoro/go.sum ./
RUN go mod download

COPY pomodoro/ .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]