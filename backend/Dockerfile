# --- Build stage ---
FROM golang:1.27-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o calculator-server .

# --- Run stage ---
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/calculator-server .
EXPOSE 8080
CMD ["./calculator-server"]