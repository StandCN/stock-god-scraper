# 使用官方的 Go 语言镜像作为基础镜像
FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到工作目录中
COPY . .

# 下载并安装依赖项并构建可执行文件
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 检查 .env 文件是否存在，如果不存在则创建一个空的 .env 文件
RUN if [ ! -f /app/.env ]; then touch /app/.env; fi

# 使用一个较小的基础镜像来运行应用程序
FROM alpine:3.14

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 安装 tzdata 包以支持时区
RUN apk add --no-cache tzdata

# 设置工作目录
WORKDIR /root/

# 将构建阶段的可执行文件复制到运行阶段
COPY --from=builder /app/main .

# 将 .env 文件从 builder 阶段复制到运行阶段
COPY --from=builder /app/.env /root/.env

# 将 config.toml 文件从 builder 阶段复制到运行阶段
COPY --from=builder /app/config/config.toml /root/config/config.toml

# 运行可执行文件
CMD ["./main"]