FROM ubuntu:20.04

# Install necessary dependencies
RUN apt-get update && apt-get install -y wget build-essential software-properties-common

# Download and install Go
RUN wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.18.10.linux-amd64.tar.gz && \
    rm go1.18.10.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin

# Install wkhtmltopdf and its dependencies
RUN apt-get update && apt-get install -y \
    libc6 \
    libfreetype6 \
    libjpeg-turbo8 \
    libpng16-16 \
    libssl1.1 \
    libstdc++6 \
    libx11-6 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    zlib1g

RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb && \
    dpkg -i wkhtmltox_0.12.6-1.focal_amd64.deb && \
    apt-get -f install -y && \
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
