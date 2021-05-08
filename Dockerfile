FROM golang:alpine

ENV CGO_ENABLED=0

RUN mkdir /data
RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go test ./...
RUN go build ./cmd/obibot

CMD ["/app/obibot"]