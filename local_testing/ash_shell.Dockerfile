FROM golang:1.24-alpine

# Install required packages
RUN apk add --no-cache \
    # Required for make test
    make \
    # Required to run bash tests
    bash \
    # Required to run zsh tests
    zsh \
    # Required for fixtures
    python3

# Set working directory
WORKDIR /app

# Starting from Go 1.20, the go standard library is no loger compiled.
# Setting GODEBUG to "installgoroot=all" restores the old behavior
RUN GODEBUG="installgoroot=all" go install std

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Set the default shell to ash
SHELL ["/bin/ash", "-c"]

# Default command
CMD ["/bin/ash"]