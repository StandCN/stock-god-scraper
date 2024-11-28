# 使用官方的 Go 语言镜像作为基础镜像
FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到工作目录中
COPY . .

# 下载并安装依赖项并构建可执行文件
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用一个较小的基础镜像来运行应用程序
FROM alpine:3.14

# 设置工作目录
WORKDIR /root/

# 将构建阶段的可执行文件复制到运行阶段
COPY --from=builder /app/main .

# 复制 .env 文件和 config/config.toml 文件（如果存在）
RUN if [ -f /app/.env ]; then cp /app/.env /root/.env; fi && \
	if [ -f /app/config/config.toml ]; then cp /app/config/config.toml /root/config/config.toml; fi

# 运行可执行文件
CMD ["./main"]
