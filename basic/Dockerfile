# Step 1: Use an official Golang image as a base
FROM golang:1.23-alpine AS builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files
COPY src/go.mod src/go.sum ./

# Step 4: Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Step 5: Copy the source code into the container
COPY src/ .

# Step 6: Build the Go application
RUN go build -gcflags=all="-N -l" -o main .

# Step 7: Multi-stage build to reduce the final image size
FROM alpine:latest

# Step 8: Set the Current Working Directory inside the container
WORKDIR /app

# Step 9: Copy the binary from the builder stage
COPY --from=builder /app/main .

# Step 10: Copy templates and static files into the container
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Step 11: Expose port 8080 to the outside world
EXPOSE 8080

# Step 12: Command to run the executable
CMD ["./main"]