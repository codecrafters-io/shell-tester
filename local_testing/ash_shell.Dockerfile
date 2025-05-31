FROM alpine:latest

# Set Go version
ARG GO_VERSION=1.23.0
ENV GO_VERSION=${GO_VERSION}

# Install required packages
RUN apk add --no-cache \
    wget \
    git \
    gcc \
    musl-dev \
    make \
    ca-certificates

# Install Go
RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz" && \
    rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Set working directory
WORKDIR /app

# Set Go environment variables to handle certificate issues
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV GOINSECURE=*
ENV GOTOOLCHAIN=local

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Set the default shell to ash
SHELL ["/bin/ash", "-c"]

# Default command
CMD ["/bin/ash"] 