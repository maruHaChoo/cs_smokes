FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/cs-smokes-bot ./cmd/bot

FROM alpine:3.20

RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app

COPY --from=builder /bin/cs-smokes-bot /app/cs-smokes-bot

ENTRYPOINT ["/app/cs-smokes-bot"]
