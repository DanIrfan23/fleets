FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk update && apk add --no-cache ca-certificates

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main .

FROM scratch

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"]
