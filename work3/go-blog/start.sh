#!/bin/bash

echo "加载环境变量..."
set -a
source .env
set +a

echo "环境变量加载成功!"
echo "启动服务器..."

go run main.go 