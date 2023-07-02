FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app .

ENV TOKEN = ""

CMD ./app $TOKEN
