# Base image
FROM golang:1.19-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod .
COPY go.sum .

# Download module dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /golang-fiber_app .

# Create minimal image
FROM scratch

# Copy the built binary into the image
COPY --from=0 /golang-fiber_app /golang-fiber_app

# Expose port 3000
EXPOSE 3000

# Run the app
CMD ["/golang-fiber_app"]
