FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/sillysocks-GoTH ./cmd/app

FROM scratch
WORKDIR /root

COPY --from=builder /app/bin .

EXPOSE 3000

CMD ["./sillysocks-GoTH"]
