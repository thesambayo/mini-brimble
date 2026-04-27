FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/main.go

FROM debian:bookworm-slim AS runner
WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    curl \
    ca-certificates \
    docker.io \
    && rm -rf /var/lib/apt/lists/*

# Install mise
RUN curl https://mise.jdx.dev/install.sh | sh \
    && cp /root/.local/bin/mise /usr/local/bin/mise \
    && chmod +x /usr/local/bin/mise

# Install Railpack
RUN curl -sSL https://railpack.com/install.sh | sh

# Create the symlink Railpack expects for mise
RUN mkdir -p /tmp/railpack/mise \
    && ln -s /usr/local/bin/mise /tmp/railpack/mise/mise-2026.3.17 \
    && chmod -R 777 /tmp/railpack

COPY --from=builder /app/server .

RUN mkdir -p /app/data /app/workspace

EXPOSE 3001

CMD ["./server"]
