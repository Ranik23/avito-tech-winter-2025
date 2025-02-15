FROM golang:1.23-alpine AS builder

WORKDIR /avito-tech

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app ./cmd/main/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache postgresql-client curl

COPY --from=builder /avito-tech/app ./

EXPOSE 8080

CMD ./app
