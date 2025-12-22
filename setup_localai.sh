#!/bin/bash

# 本地AI系统快速设置脚本
# 用于 gin-vue-admin 项目

set -e

echo "================================================"
echo "  本地AI系统 - 快速设置"
echo "================================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查是否已安装 Ollama
check_ollama() {
    echo "检查 Ollama 安装状态..."
    if command -v ollama &> /dev/null; then
        echo -e "${GREEN}✓ Ollama 已安装${NC}"
        return 0
    else
        echo -e "${RED}✗ Ollama 未安装${NC}"
        return 1
    fi
}

# 安装 Ollama
install_ollama() {
    echo ""
    echo "正在安装 Ollama..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        echo "检测到 macOS 系统"
        if command -v brew &> /dev/null; then
            brew install ollama
        else
            curl -fsSL https://ollama.com/install.sh | sh
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        echo "检测到 Linux 系统"
        curl -fsSL https://ollama.com/install.sh | sh
    else
        echo -e "${YELLOW}请手动访问 https://ollama.com/download 下载安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Ollama 安装完成${NC}"
}

# 启动 Ollama 服务
start_ollama() {
    echo ""
    echo "启动 Ollama 服务..."
    
    # 检查是否已经在运行
    if pgrep -x "ollama" > /dev/null; then
        echo -e "${GREEN}✓ Ollama 服务已在运行${NC}"
    else
        # 在后台启动
        nohup ollama serve > /tmp/ollama.log 2>&1 &
        sleep 3
        
        if pgrep -x "ollama" > /dev/null; then
            echo -e "${GREEN}✓ Ollama 服务已启动${NC}"
        else
            echo -e "${RED}✗ Ollama 服务启动失败${NC}"
            echo "请手动运行: ollama serve"
            exit 1
        fi
    fi
}

# 下载模型
download_model() {
    echo ""
    echo "可用的模型："
    echo "  1) qwen2.5:7b    - 通义千问 7B (推荐，约 4GB)"
    echo "  2) qwen2.5:3b    - 通义千问 3B (轻量，约 2GB)"
    echo "  3) llama3.1:8b   - Llama 3.1 8B (约 4.5GB)"
    echo "  4) phi3:mini     - Phi-3 Mini (约 2.5GB)"
    echo "  5) 跳过下载"
    echo ""
    read -p "请选择要下载的模型 (1-5): " choice
    
    case $choice in
        1)
            MODEL="qwen2.5:7b"
            ;;
        2)
            MODEL="qwen2.5:3b"
            ;;
        3)
            MODEL="llama3.1:8b"
            ;;
        4)
            MODEL="phi3:mini"
            ;;
        5)
            echo "跳过模型下载"
            return
            ;;
        *)
            echo "无效选择，使用默认: qwen2.5:7b"
            MODEL="qwen2.5:7b"
            ;;
    esac
    
    echo ""
    echo "正在下载模型: $MODEL ..."
    echo "这可能需要几分钟时间，请耐心等待..."
    
    if ollama pull $MODEL; then
        echo -e "${GREEN}✓ 模型下载完成: $MODEL${NC}"
        
        # 更新配置文件
        if [ -f "server/config.yaml" ]; then
            echo ""
            echo "更新配置文件..."
            # 这里可以用 sed 或其他工具更新配置
            echo -e "${YELLOW}请手动更新 server/config.yaml 中的 model 字段为: $MODEL${NC}"
        fi
    else
        echo -e "${RED}✗ 模型下载失败${NC}"
        exit 1
    fi
}

# 创建必要的目录
create_directories() {
    echo ""
    echo "创建数据目录..."
    
    cd server
    
    mkdir -p data/vector
    mkdir -p data/training
    mkdir -p knowledge_base
    mkdir -p models/finetuned
    
    echo -e "${GREEN}✓ 目录创建完成${NC}"
    
    cd ..
}

# 检查 Go 环境
check_go() {
    echo ""
    echo "检查 Go 环境..."
    
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        echo -e "${GREEN}✓ Go 已安装: $GO_VERSION${NC}"
    else
        echo -e "${RED}✗ Go 未安装${NC}"
        echo "请访问 https://go.dev/dl/ 下载安装"
        exit 1
    fi
}

# 测试连接
test_connection() {
    echo ""
    echo "测试 Ollama 连接..."
    
    if curl -s http://localhost:11434/api/tags > /dev/null; then
        echo -e "${GREEN}✓ Ollama API 可访问${NC}"
        
        echo ""
        echo "已安装的模型列表："
        ollama list
    else
        echo -e "${RED}✗ 无法连接到 Ollama${NC}"
        echo "请确保 Ollama 服务正在运行"
    fi
}

# 显示使用说明
show_usage() {
    echo ""
    echo "================================================"
    echo "  设置完成！"
    echo "================================================"
    echo ""
    echo "下一步操作："
    echo ""
    echo "1. 启动 gin-vue-admin 服务器："
    echo "   cd server"
    echo "   go run main.go"
    echo ""
    echo "2. 测试本地AI接口："
    echo "   curl -X POST http://localhost:8888/localai/chat \\"
    echo "     -H \"Content-Type: application/json\" \\"
    echo "     -d '{\"message\": \"你好\", \"use_rag\": false}'"
    echo ""
    echo "3. 导入文档到知识库："
    echo "   curl -X POST http://localhost:8888/localai/ingest \\"
    echo "     -H \"Content-Type: application/json\" \\"
    echo "     -d '{\"content\": \"你的文档内容\"}'"
    echo ""
    echo "4. 使用RAG进行对话："
    echo "   curl -X POST http://localhost:8888/localai/chat \\"
    echo "     -H \"Content-Type: application/json\" \\"
    echo "     -d '{\"message\": \"问题\", \"use_rag\": true, \"top_k\": 5}'"
    echo ""
    echo "详细文档请查看: server/api/v1/localai/README.md"
    echo ""
}

# 主流程
main() {
    echo "开始设置本地AI系统..."
    echo ""
    
    # 检查 Go
    check_go
    
    # 检查/安装 Ollama
    if ! check_ollama; then
        read -p "是否安装 Ollama? (y/n): " install
        if [[ $install == "y" || $install == "Y" ]]; then
            install_ollama
        else
            echo "跳过 Ollama 安装"
            exit 0
        fi
    fi
    
    # 启动服务
    start_ollama
    
    # 下载模型
    download_model
    
    # 创建目录
    create_directories
    
    # 测试连接
    test_connection
    
    # 显示使用说明
    show_usage
}

# 运行主流程
main


