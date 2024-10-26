FROM golang:1.23 AS builder

WORKDIR /project

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mainapp ./cmd/main.go

#========================================

FROM alpine:latest

WORKDIR /project

COPY --from=builder /project/mainapp .
COPY .env .

RUN apk update && apk add bash gcc

CMD ["./mainapp"]

