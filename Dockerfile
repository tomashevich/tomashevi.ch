FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache build-base

COPY . .

RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o tomashevich .

FROM alpine:latest
RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/tomashevich .

EXPOSE 8037

ENTRYPOINT ["./tomashevich"]
