#!/bin/sh

# 修复 Go 依赖和清理缓存
# 在 Docker 容器中运行

echo "清理 Go 缓存..."
go clean -cache -modcache -i -r

echo "清理 vendor 目录..."
rm -rf vendor

echo "更新依赖..."
go mod tidy

echo "下载依赖..."
go mod download

echo "验证依赖..."
go mod verify

echo "完成！现在可以运行 go run main.go"

