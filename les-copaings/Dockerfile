FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app .

ENV TOKEN = ""

CMD sleep 2 && ./app $TOKEN
