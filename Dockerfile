FROM golang:1.18.1-alpine
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ENV DBHOST ${{ secrets.DB_HOST }}
ENV DBUSER ${{ secrets.DB_USER }}
ENV DBPASSWORD ${{ secrets.DB_PASSWORD }}
ENV DBPORT 3306
ENV DB akafuyu

COPY ./ ./

RUN go build -o ./build/server.exe ./server.go

EXPOSE 7700
CMD ["./build/server.exe"]