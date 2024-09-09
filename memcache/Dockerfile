FROM golang:1.22

WORKDIR /app

COPY ./memcache .  

WORKDIR /app
RUN go mod download

WORKDIR /app/cmd/server

RUN GOOS=linux go build -o memcache

EXPOSE 8081

CMD ["./memcache"]

