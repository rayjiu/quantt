# 第一阶段：构建应用程序
FROM golang:1.23.2-bullseye AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源代码
COPY . .

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crawler ./cmd/crawler

# 第二阶段：创建一个最小化的运行镜像
FROM alpine:latest

RUN mkdir -p /var/log/app

# 设置时区（可选）
ENV TZ=Asia/Shanghai

# 将可执行文件从构建阶段复制到运行阶段
COPY --from=builder /app/crawler /crawler

# 暴露端口（如果需要）
# EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["/crawler"]
