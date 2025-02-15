FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/sillysocks-GoTH ./cmd/sillysocks
RUN apk add --no-cache ca-certificates

FROM scratch
WORKDIR /root

COPY --from=builder /app/bin .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


EXPOSE 8080

CMD ["./sillysocks-GoTH"]
