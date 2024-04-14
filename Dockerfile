FROM golang:1.22.2-alpine3.19

WORKDIR /app

RUN go mod init avito_app
#RUN go get -u github.com/lib/pq github.com/golang-jwt/jwt/v5 golang.org/x/crypto/bcrypt github.com/swaggo/http-swagger/v2 go get github.com/redis/go-redis/v9
RUN go get -u github.com/lib/pq  golang.org/x/crypto/bcrypt github.com/redis/go-redis/v9 github.com/justinas/alice

COPY ./ ./

EXPOSE 5000

ENTRYPOINT [ "go", "run", "./cmd" ]