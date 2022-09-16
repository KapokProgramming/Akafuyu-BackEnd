FROM golang:1.18.1-alpine
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o ./build/server.exe ./server.go

EXPOSE 7700
CMD ["./build/server.exe"]