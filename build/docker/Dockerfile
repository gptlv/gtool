FROM ubuntu:20.04

# Install necessary dependencies
RUN apt-get update && apt-get install -y wget build-essential

# Download and install Go
RUN wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.18.10.linux-amd64.tar.gz && \
    rm go1.18.10.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin

# Install wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb && \
    apt install -y ./wkhtmltox_0.12.6-1.focal_amd64.deb && \
    rm wkhtmltox_0.12.6-1.focal_amd64.deb

# Set working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o scripts cmd/scripts/main.go

# Set the entrypoint to the built binary
CMD ["./scripts"]
