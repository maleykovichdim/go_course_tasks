
FROM golang:1.22

WORKDIR /app

COPY . .  

WORKDIR /app
RUN go mod download

WORKDIR /app/cmd/shortener

RUN GOOS=linux go build -o app

EXPOSE 8080

CMD ["./app"]

#docker build -f .\build\Dockerfile --progress=plain -t maleykovich/link-shortener:latest .
#docker run --rm -p 8080:8080 maleykovich/link-shortener:latest

