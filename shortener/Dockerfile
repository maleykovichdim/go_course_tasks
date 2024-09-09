FROM golang:1.22

WORKDIR /app

COPY ./shortener .  

WORKDIR /app
RUN go mod download

WORKDIR /app/cmd/server

RUN GOOS=linux go build -o shortener

EXPOSE 8080

CMD ["./shortener"]

