FROM golang:alpine

# Set necessary environment variables needed for the image
ENV GO111MODULE=on \
    GCO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
ADD build ./build
ADD internal ./internal
COPY main.go .

# Build the application
RUN go build -o trackmyfish main.go

# Move to /app directory as the place for resulting binary folder
WORKDIR /app

# Copy binary from build to main folder
RUN cp /build/trackmyfish .

# Define a volume to hold the config file
VOLUME /config

# Expose the web port
EXPOSE 9000

# Command to run when starting the container
CMD ["/app/trackmyfish"]
