FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o echo-server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/echo-server .
EXPOSE 8080
CMD ["./echo-server"]