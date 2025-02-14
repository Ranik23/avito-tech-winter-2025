FROM golang:1.23-alpine

WORKDIR /avito-tech

RUN apk add --no-cache make

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN make run

EXPOSE 8080
