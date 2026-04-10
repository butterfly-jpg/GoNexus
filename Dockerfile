# ===== 构建阶段 =====
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装依赖（alpine 需要 gcc 支持 cgo）
RUN apk add --no-cache git

# 优先复制 go.mod / go.sum，利用缓存加速依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码
COPY . .

# 编译 gonexus 主服务（CGO_ENABLED=0 生成静态二进制，适合 alpine 运行）
RUN CGO_ENABLED=0 GOOS=linux go build -o gonexus .

# ===== 运行阶段 =====
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制二进制
COPY --from=builder /app/gonexus .
# 复制配置目录（config.toml 通过 volume 挂载覆盖）
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./gonexus"]
