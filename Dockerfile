# Build stage
FROM golang:1.25-alpine AS builder

# Install git and SSL certificates
RUN apk update && apk add --no-cache git ca-certificates

# Set the working directory
WORKDIR /app

# Copy module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Set build-time variables
ARG VERSION
ARG BUILD_DATE

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
  -ldflags "-X github.com/particledecay/slackmoji-notifier/pkg/build.Version=${VERSION} \
  -X github.com/particledecay/slackmoji-notifier/pkg/build.Date=${BUILD_DATE}" \
  -o main .

# Final stage
FROM alpine:latest

# Install SSL certificates
RUN apk update && apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/main .

# Executable
ENTRYPOINT [ "/app/main" ]
CMD [ "listen" ]
