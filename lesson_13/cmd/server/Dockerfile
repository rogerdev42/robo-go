FROM golang:1.23.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server-app ./cmd/server

EXPOSE 8080

CMD ["./server-app"]
