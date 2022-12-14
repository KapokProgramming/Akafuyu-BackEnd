FROM golang:1.18.1-alpine
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ENV DBPORT 3306
ENV DB akafuyu
ENV ORIGIN_ALLOWED *

COPY ./ ./

RUN go build -o ./build/server.exe

EXPOSE 7700
CMD ["./build/server.exe"]
