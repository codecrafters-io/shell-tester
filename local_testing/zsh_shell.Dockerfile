FROM golang:1.24-bookworm

# Install required packages
RUN apt-get update && apt-get install -y zsh make

# Set working directory
WORKDIR /app

# Starting from Go 1.20, the go standard library is no longer compiled.
# Setting GODEBUG to "installgoroot=all" restores the old behavior
RUN GODEBUG="installgoroot=all" go install std

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Set the default shell to zsh
SHELL ["/bin/zsh", "-c"]

# Default command
CMD ["/bin/zsh"]