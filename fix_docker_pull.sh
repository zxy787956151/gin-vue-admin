#!/bin/bash

# Docker Hub 镜像拉取问题修复脚本
# 解决国内网络访问 Docker Hub 超时问题

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[✓]${NC} $1"; }
print_error() { echo -e "${RED}[✗]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[!]${NC} $1"; }

echo ""
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║        Docker Hub 镜像拉取问题修复工具                    ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    print_error "未安装 Docker"
    exit 1
fi

if ! docker info &> /dev/null 2>&1; then
    print_error "Docker 未运行，请启动 Docker Desktop"
    exit 1
fi

print_success "Docker 运行正常"
echo ""

# 提供解决方案选择
echo "请选择解决方案:"
echo ""
echo "1) 推荐：改用本地安装 Ollama（无需 Docker，更简单）"
echo "2) 配置 Docker 镜像加速器（国内镜像源）"
echo "3) 手动拉取镜像（使用代理）"
echo "4) 从其他源导入镜像"
echo "5) 退出"
echo ""
read -p "请选择 [1-5]: " CHOICE

case $CHOICE in
    1)
        # 方案1: 本地安装
        print_info "推荐方案：本地安装 Ollama"
        echo ""
        echo "对于 Mac，本地安装 Ollama 更简单、更快、更稳定！"
        echo ""
        
        if command -v ollama &> /dev/null; then
            print_success "Ollama 已安装"
            ollama --version
        else
            if command -v brew &> /dev/null; then
                print_info "开始安装 Ollama..."
                brew install ollama
                print_success "安装完成"
            else
                print_warning "未安装 Homebrew"
                echo ""
                echo "请选择以下方式之一："
                echo "1. 安装 Homebrew 后再运行此脚本"
                echo "   /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
                echo ""
                echo "2. 直接下载 Ollama App"
                echo "   https://ollama.com/download/mac"
                exit 0
            fi
        fi
        
        # 启动服务
        print_info "启动 Ollama 服务..."
        brew services start ollama
        sleep 3
        
        # 验证
        if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
            print_success "Ollama 服务运行正常"
            echo ""
            print_info "下一步："
            echo "1. 拉取模型: ollama pull qwen2.5:7b"
            echo "2. 修改 config.yaml:"
            echo "   local-ai.llm.base-url: \"http://localhost:11434\""
            echo ""
            read -p "是否现在拉取模型? [Y/n]: " PULL_MODEL
            if [[ ! "$PULL_MODEL" =~ ^[Nn]$ ]]; then
                print_info "拉取 qwen2.5:7b 模型（约 4.7GB）..."
                ollama pull qwen2.5:7b
                print_success "模型下载完成！"
                
                # 测试
                print_info "运行测试..."
                ollama list
            fi
        else
            print_error "Ollama 服务启动失败"
        fi
        ;;
        
    2)
        # 方案2: 配置镜像加速
        print_info "配置 Docker 镜像加速器"
        echo ""
        print_warning "注意：Docker Desktop for Mac 需要手动配置"
        echo ""
        echo "请按以下步骤操作："
        echo ""
        echo "1. 打开 Docker Desktop"
        echo "2. 点击右上角 ⚙️ (Settings)"
        echo "3. 选择 'Docker Engine'"
        echo "4. 在 JSON 配置中添加以下内容："
        echo ""
        echo "{"
        echo "  \"registry-mirrors\": ["
        echo "    \"https://docker.m.daocloud.io\","
        echo "    \"https://dockerproxy.com\","
        echo "    \"https://docker.nju.edu.cn\","
        echo "    \"https://docker.mirrors.sjtug.sjtu.edu.cn\""
        echo "  ]"
        echo "}"
        echo ""
        echo "5. 点击 'Apply & Restart'"
        echo "6. 等待 Docker 重启"
        echo ""
        read -p "配置完成后按 Enter 继续..."
        
        # 验证
        print_info "验证配置..."
        if docker info | grep -A 10 "Registry Mirrors" > /dev/null; then
            print_success "镜像加速器配置成功"
            docker info | grep -A 10 "Registry Mirrors"
        else
            print_warning "未检测到镜像加速器配置"
        fi
        
        # 尝试拉取
        print_info "尝试拉取 Ollama 镜像..."
        if docker pull ollama/ollama:latest; then
            print_success "镜像拉取成功！"
        else
            print_error "拉取失败，请尝试其他方案"
        fi
        ;;
        
    3)
        # 方案3: 使用代理
        print_info "使用代理拉取镜像"
        echo ""
        echo "如果你有 HTTP/HTTPS 代理，请输入代理地址"
        echo "格式: http://127.0.0.1:7890"
        echo ""
        read -p "代理地址 (留空跳过): " PROXY_URL
        
        if [[ -n "$PROXY_URL" ]]; then
            export HTTP_PROXY="$PROXY_URL"
            export HTTPS_PROXY="$PROXY_URL"
            print_info "已设置代理: $PROXY_URL"
            
            print_info "尝试拉取镜像..."
            if docker pull ollama/ollama:latest; then
                print_success "镜像拉取成功！"
            else
                print_error "拉取失败"
            fi
        else
            print_warning "未设置代理"
        fi
        ;;
        
    4)
        # 方案4: 从其他源导入
        print_info "从其他源获取镜像"
        echo ""
        echo "可以尝试以下方式："
        echo ""
        echo "1. 使用阿里云容器镜像服务"
        echo "   docker pull registry.cn-hangzhou.aliyuncs.com/xxx/ollama:latest"
        echo ""
        echo "2. 使用 GitHub Container Registry"
        echo "   docker pull ghcr.io/ollama/ollama:latest"
        echo ""
        echo "3. 从已有镜像的机器导出/导入"
        echo "   docker save ollama/ollama:latest | gzip > ollama.tar.gz"
        echo "   docker load < ollama.tar.gz"
        echo ""
        print_warning "以上镜像源可能不是官方维护，请谨慎使用"
        ;;
        
    5)
        print_info "退出"
        exit 0
        ;;
        
    *)
        print_error "无效选择"
        exit 1
        ;;
esac

echo ""
echo "═══════════════════════════════════════════════════════════"
print_info "完成！"
echo ""
echo "💡 推荐："
echo "   对于 Mac 用户，强烈推荐使用本地安装而不是 Docker："
echo "   brew install ollama"
echo ""
echo "   优势："
echo "   - 安装更快（无需下载大镜像）"
echo "   - 启动更快（无容器开销）"
echo "   - 管理更简单（brew services）"
echo "   - 性能更好（原生运行）"
echo ""

