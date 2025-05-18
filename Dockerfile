FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /rentapp

FROM alpine:latest

COPY --from=builder /app/config-ci.json /config.json

COPY --from=builder /rentapp /rentapp

USER 1000

ENTRYPOINT ["/rentapp"]



