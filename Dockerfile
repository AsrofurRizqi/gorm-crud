# syntax=docker/dockerfile:1

FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /gorm-crud

# Copy go mod and sum files
COPY . ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o apps

# Command to run the executable
ENTRYPOINT ["/gorm-crud/apps"]