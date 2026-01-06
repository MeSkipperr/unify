FROM golang:1.24.5-alpine

WORKDIR /app

RUN apk add --no-cache iputils

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

EXPOSE 8080

CMD ["./app"]