# Use the official Golang image as the base image
FROM golang:1.24-alpine AS go-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.22
RUN apk --no-cache add ca-certificates libc6-compat
RUN mkdir -p /app/data /app/uploads
WORKDIR /app
COPY --from=go-builder /app/main .
COPY --from=go-builder /app/internal/db/migration /app/internal/db/migration
COPY --from=go-builder /app/static /app/static
COPY --from=go-builder /app/inspirations /app/inspirations
EXPOSE 8080
CMD ["./main"]
