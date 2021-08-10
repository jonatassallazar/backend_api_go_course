FROM golang:1.16.7-alpine3.14 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY models ./models

RUN go build ./cmd/api/

EXPOSE 4000 4000

CMD [ "./api" ]