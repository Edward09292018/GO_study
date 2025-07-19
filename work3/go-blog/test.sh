#!/bin/bash

# 定义颜色常量
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 设置基础URL
BASE_URL="http://localhost:8080"

# 记录日志函数
log() {
    echo -e "$1" | tee -a "api_test.log"
}

# 检查命令是否成功执行
check_status() {
    if [ $? -ne 0 ]; then
        log "${RED}错误: 上一步操作失败，停止执行${NC}"
        exit 1
    fi
}

# 显示测试结果
check_result() {
    if [[ $2 == *"success"* ]]; then
        log "${GREEN}$1 成功: $2${NC}"
    else
        log "${RED}$1 失败: $2${NC}"
    fi
}

# 测试注册接口
test_register() {
    log "${YELLOW}测试注册接口...${NC}"

    # 测试成功注册
    register_data='{"username":"testuser","password":"password123","email":"test@example.com"}'
    response=$(curl -s -X POST "$BASE_URL/register" -H "Content-Type: application/json" -d "$register_data")
    check_result "注册" "$response"

    # 测试重复注册
    response=$(curl -s -X POST "$BASE_URL/register" -H "Content-Type: application/json" -d "$register_data")
    check_result "重复注册" "$response"
}

# 测试登录接口
test_login() {
    log "${YELLOW}测试登录接口...${NC}"

    # 测试成功登录
    login_data='{"username":"testuser","password":"password123"}'
    response=$(curl -s -X POST "$BASE_URL/login" -H "Content-Type: application/json" -d "$login_data")
    token=$(echo "$response" | grep -o '"data":"[^"]*' | cut -d'"' -f4)

    if [ -n "$token" ]; then
        log "${GREEN}登录成功，获取到token: $token${NC}"
    else
        log "${RED}登录失败，未获取到token${NC}"
        exit 1
    fi

    echo "$token" > ".token"

    # 测试错误密码登录
    login_data='{"username":"testuser","password":"wrongpassword"}'
    response=$(curl -s -X POST "$BASE_URL/login" -H "Content-Type: application/json" -d "$login_data")
    check_result "错误密码登录" "$response"

    return 0
}

# 测试文章接口
test_posts() {
    log "${YELLOW}测试文章接口...${NC}"

    # 读取token
    token=$(cat ".token")

    # 创建测试文章
    post_data='{"title":"Test Post","content":"This is a test post"}'
    response=$(curl -s -X POST "$BASE_URL/posts" -H "Content-Type: application/json" -H "Authorization: $token" -d "$post_data")

    # 修复了字段名匹配，使用正确的 "ID" 字段
    post_id=$(echo "$response" | grep -o '"ID":[0-9]*' | head -n 1 | grep -o '[0-9]*')

    if [ -n "$post_id" ]; then
        log "${GREEN}文章创建成功，ID: $post_id${NC}"
    else
        log "${RED}文章创建失败${NC}"
        exit 1
    fi
    echo "$post_id" > ".post_id"

    # 获取所有文章
    response=$(curl -s -X GET "$BASE_URL/posts" -H "Authorization: $token")
    check_result "获取所有文章" "$response"

    # 获取单个文章

    response=$(curl -s -X GET "$BASE_URL/posts/$post_id" -H "Authorization: $token")
    echo "$response"
    check_result "获取单个文章" "$response"

    # 更新文章
    updated_post_data='{"title":"Updated Test Post","content":"This is an updated test post"}'
    response=$(curl -s -X PUT "$BASE_URL/posts/$post_id" -H "Content-Type: application/json" -H "Authorization: $token" -d "$updated_post_data")
    check_result "更新文章" "$response"

}

# 测试评论接口
test_comments() {
    log "${YELLOW}测试评论接口...${NC}"

    # 读取token和文章ID
    token=$(cat ".token")
    post_id=$(cat ".post_id")

    # 创建评论
    comment_data='{"content":"This is a test comment","PostID":'"$post_id"'}'

    log "${YELLOW}发送请求: curl -s -X POST \"$BASE_URL/comments\" -H \"Content-Type: application/json\" -H \"Authorization: $token\" -d \"$comment_data\"${NC}"
    response=$(curl -s -X POST "$BASE_URL/comments" -H "Content-Type: application/json" -H "Authorization: $token" -d "$comment_data")
    check_result "创建评论" "$response"

    # 解析评论ID
    comment_id=$(echo "$response" | grep -o '"ID":[0-9]*' | grep -o '[0-9]*')

    # 获取文章的评论
    response=$(curl -s -X GET "$BASE_URL/posts/$post_id/comments" -H "Authorization: $token")
    check_result "获取文章的评论" "$response"

     # 删除文章
    response=$(curl -s -X DELETE "$BASE_URL/posts/$post_id" -H "Authorization: $token")
    check_result "删除文章" "$response"
}



# 主函数
main() {
    log "${BLUE}开始测试Go-Blog API接口$(date)${NC}"

    # 测试注册接口
    test_register

    # 测试登录接口
    test_login

    # 测试文章接口
    test_posts

    # 测试评论接口
    test_comments

    log "${BLUE}所有API接口测试完成${NC}"
}

# 执行主函数
main
