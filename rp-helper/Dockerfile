FROM golang:1.19-alpine

WORKDIR /app

RUN mkdir /config

COPY . .

RUN go mod tidy && go build -o app .

ENV TOKEN = ""

CMD sleep 2 && ./app $TOKEN
