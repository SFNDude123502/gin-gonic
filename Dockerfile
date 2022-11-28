FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build

EXPOSE 8080

CMD [ "./gingonic" ]