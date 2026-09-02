FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./main.go

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server .

COPY --from=builder /app/migrations ./migrations
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./server"]