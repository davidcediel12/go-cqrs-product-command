FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o product-command ./cmd/app


FROM alpine:3

# Add CA certificates (for HTTPS calls to SNS, DB, etc.)
RUN apk --no-cache add ca-certificates

WORKDIR /app


# Use a non-root user
RUN adduser -D appuser

# Copy the binary from builder
COPY --from=builder --chown=appuser:appuser /app/product-command .

# USER appuser

EXPOSE 8080

ENTRYPOINT [ "./product-command" ]