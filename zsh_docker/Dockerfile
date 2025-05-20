FROM ubuntu:latest

# Set Go version
ARG GO_VERSION=1.23.0
ENV GO_VERSION=${GO_VERSION}

# Install required packages
RUN apt-get update && apt-get install -y \
    wget \
    git \
    gcc \
    zsh \
    make \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

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

# Copy the project files
COPY . .

# Build the tester
RUN make build

# Set the default shell to zsh
SHELL ["/bin/zsh", "-c"]

# Default command
CMD ["/bin/zsh"] 