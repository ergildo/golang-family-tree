FROM golang:1.21-alpine3.18 AS builder

WORKDIR  /build

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/golang-family-tree-api ./cmd/...

FROM alpine:3.18

WORKDIR  /app

COPY --from=builder /build/golang-family-tree-api /app/golang-family-tree-api
COPY --from=builder /build/migrations /app/migrations

EXPOSE 8080

CMD ["/app/golang-family-tree-api"]