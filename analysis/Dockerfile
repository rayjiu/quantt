# Stage 1: Modules caching
FROM golang:1.23.2-bullseye as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Stage 2: Build
FROM golang:1.23.2-bullseye as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /analysis
WORKDIR /analysis
# Build your app (only in production)
ARG ENV=production
RUN if [ "$ENV" = "production" ]; then GOOS=linux GOARCH=amd64 go build -o build/analysis ./cmd/analysis; fi

# Stage 3: Final
FROM ubuntu:jammy
# Tsinghua Ubuntu source
RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|https://mirrors.tuna.tsinghua.edu.cn/ubuntu/|g' /etc/apt/sources.list

# Copy binaries or source code based on the environment
COPY --from=builder /analysis /analysis
# Alway copy to go the final stage.
COPY --from=builder /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

# COPY --from=builder /go/bin/playwright /crawler /
ARG ENV=production

# Install dependencies and all Playwright browsers
RUN apt-get update && apt-get install -y ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /analysis

RUN if [ "$ENV" = "production" ]; then \
    cp /analysis/build/analysis /analysis; \
    rm -rf /usr/local/go/; \
fi

CMD if [ "$ENV" = "production" ]; then /analysis; else go run ./cmd/analysis; fi
