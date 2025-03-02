# 构建阶段
FROM --platform=linux/amd64 golang:1.22-alpine as builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev

WORKDIR /usr/src/game-fun-be

# 设置 Go 环境
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    GOSUMDB=off \
    GOTOOLCHAIN=local \
    GOPATH=/go \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# 先复制依赖文件
COPY go.mod go.sum ./

# 下载依赖并显示进度
RUN go mod download -x && \
    go list -m all

# 再复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 go build -o game-fun-be main.go

# 运行阶段
FROM --platform=linux/amd64 golang:1.22-alpine

# 安装运行时依赖
RUN apk add --no-cache tzdata ca-certificates

# 创建应用目录
RUN mkdir -p /opt/game-fun-be

# 复制二进制文件和配置
COPY --from=builder /usr/src/game-fun-be/game-fun-be /opt/game-fun-be/game-fun-be
COPY --from=builder /usr/src/game-fun-be/docs /opt/game-fun-be/docs

# 设置工作目录
WORKDIR /opt/game-fun-be

# 创建日志目录
RUN mkdir -p /opt/game-fun-be/logs && \
    chmod -R 755 /opt/game-fun-be/logs

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 暴露端口
EXPOSE 4880

# 启动命令
CMD ["./game-fun-be"]
