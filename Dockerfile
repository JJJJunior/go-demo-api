# 使用 golang:alpine 作为基础镜像进行构建
FROM golang:alpine AS builder

LABEL authors="Maxwellsim"

# 环境变量设置
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制当前目录的所有文件
COPY . .

# 构建可执行文件
RUN go build -o main

# 使用 scratch 作为基础镜像进行发布
FROM scratch
WORKDIR /app

# 从 builder 镜像中复制构建好的可执行文件
COPY --from=builder /build/main /app/main

# 运行可执行文件
CMD ["./main"]
