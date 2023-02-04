# build
FROM golang:1.19-alpine AS builder

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./main ./main.go

# deploy
FROM scratch

WORKDIR /app
COPY --from=builder /build/main main
COPY ./data /app/data/

CMD ["./main"]
