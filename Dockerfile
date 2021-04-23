FROM golang:1.16.0-alpine3.13

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build ./cmd/obibot

CMD ["/app/obibot"]