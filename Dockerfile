FROM golang:1.18.1-alpine
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ENV DBPORT 3306
ENV DB akafuyu
ENV ORIGIN_ALLOWED *
ENV DB_HOST 172.104.35.73
ENV DB_USER root
ENV DB_PASSWORD 4321

COPY ./ ./

RUN go build -o ./build/server.exe ./server.go

EXPOSE 7700
CMD ["./build/server.exe"]
