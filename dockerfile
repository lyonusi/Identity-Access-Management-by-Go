FROM golang:1.17.2-bullseye
WORKDIR /usr/app
COPY . .
RUN go build -o main .
EXPOSE 1323

CMD ["./main"]
