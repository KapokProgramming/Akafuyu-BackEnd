FROM golang:1.18.1-alpine
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

<<<<<<< HEAD
ENV DBHOST ${{ secrets.DB_HOST }}
ENV DBUSER ${{ secrets.DB_USER }}
ENV DBPASSWORD ${{ secrets.DB_PASSWORD }}
=======
ENV DBHOST 172.104.35.73
ENV DBUSER root
ENV DBPASSWORD 4321
>>>>>>> 65b29f158a14f0220f5a022bccade16de28231a1
ENV DBPORT 3306
ENV DB akafuyu

COPY ./ ./

RUN go build -o ./build/server.exe ./server.go

EXPOSE 7700
CMD ["./build/server.exe"]
