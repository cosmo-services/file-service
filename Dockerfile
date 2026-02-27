FROM golang:alpine as BUILDER

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

ENV CONFIG_PATH=./config.yaml

COPY . .

RUN swag init

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
#CMD ["go", "run", "main.go"]