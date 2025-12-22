# vLLM 部署指南 - 精简实用版

## 📋 你的部署方案

```
┌─────────────────────────┐      HTTP      ┌─────────────────────────┐
│  macOS 本地              │  ◄────────────►  │  CentOS 服务器           │
│  ┌───────────────────┐  │                 │  ┌───────────────────┐  │
│  │ gin-vue-admin     │  │                 │  │ vLLM + Qwen3      │  │
│  │ docker-compose    │  │                 │  │ docker-compose    │  │
│  │ 端口: 8888        │  │                 │  │ 端口: 8000        │  │
│  └───────────────────┘  │                 │  └───────────────────┘  │
│                         │                 │         GPU              │
└─────────────────────────┘                 └─────────────────────────┘
```

---

## 第一步：CentOS 服务器部署 vLLM

### 1.1 环境准备

```bash
# SSH 登录服务器
ssh root@your-centos-server

# 检查 GPU
nvidia-smi

# 安装 Docker（如果没有）
curl -fsSL https://get.docker.com | bash
sudo systemctl start docker
sudo systemctl enable docker
```

### 1.2 安装 NVIDIA Container Toolkit

```bash
# 方法1：使用通用 RPM 仓库（推荐）
curl -fsSL https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo | \
  sudo tee /etc/yum.repos.d/nvidia-container-toolkit.repo

sudo yum clean expire-cache
sudo yum install -y nvidia-container-toolkit

# 配置 Docker 运行时
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker

# 验证
docker run --rm --gpus all nvidia/cuda:12.0.0-base-centos7 nvidia-smi

# 方法2：如果方法1失败，直接安装（CentOS 7/8）
# sudo yum install -y yum-utils
# sudo yum-config-manager --add-repo https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo
# sudo yum install -y nvidia-container-toolkit
# sudo nvidia-ctk runtime configure --runtime=docker
# sudo systemctl restart docker
```

### 1.3 创建部署目录

```bash
mkdir -p ~/vllm-deploy/models
cd ~/vllm-deploy
```

### 1.4 下载模型

```bash
# 方式1: huggingface-cli
pip3 install huggingface-hub
huggingface-cli download Qwen/Qwen2.5-7B-Instruct-AWQ \
  --local-dir ./models/Qwen2.5-7B-Instruct-AWQ

# 方式2: ModelScope（国内快）
pip3 install modelscope
python3 << EOF
from modelscope import snapshot_download
snapshot_download('Qwen/Qwen2.5-7B-Instruct-AWQ', cache_dir='./models')
EOF
```

### 1.5 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  vllm:
    image: vllm/vllm-openai:latest
    container_name: vllm-qwen
    restart: unless-stopped
    runtime: nvidia
    
    ports:
      - "8000:8000"
    
    volumes:
      - ./models:/models:ro
      - ./cache:/root/.cache
    
    environment:
      - NVIDIA_VISIBLE_DEVICES=all
      - CUDA_VISIBLE_DEVICES=0
    
    command: >
      --model /models/Qwen2.5-7B-Instruct-AWQ
      --quantization awq
      --host 0.0.0.0
      --port 8000
      --trust-remote-code
      --max-model-len 4096
      --gpu-memory-utilization 0.9
      --dtype auto
```

### 1.6 启动服务

```bash
# 启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 验证
curl http://localhost:8000/health
curl http://localhost:8000/v1/models
```

---

## 第二步：macOS 本地配置 gin

### 2.1 获取服务器 IP

```bash
# 在 CentOS 服务器上执行
ip addr show | grep "inet " | grep -v 127.0.0.1
# 记录 IP，例如: 192.168.1.100
```

### 2.2 测试连通性

```bash
# 在 macOS 上测试
curl http://192.168.1.100:8000/health
curl http://192.168.1.100:8000/v1/models
```

### 2.3 修改 gin 配置

编辑 `server/config.yaml`:

```yaml
# 添加 vLLM 配置
vllm:
  enabled: true
  base-url: "http://192.168.1.100:8000"  # 改成你的服务器IP
  model: "Qwen2.5-7B-Instruct-AWQ"
  timeout: 300
  max-tokens: 2000
  temperature: 0.7
```

### 2.4 重启 gin 容器

```bash
# 在 macOS 上
cd /Users/markxiu/workspace/www/gin-vue-admin
docker-compose restart server
# 或
docker-compose up -d
```

---

## 第三步：编写 gin 集成代码

### 3.1 创建配置结构

`server/config/vllm.go`:

```go
package config

type VLLM struct {
	Enabled     bool    `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	BaseURL     string  `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	Model       string  `mapstructure:"model" json:"model" yaml:"model"`
	Timeout     int     `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	MaxTokens   int     `mapstructure:"max-tokens" json:"max-tokens" yaml:"max-tokens"`
	Temperature float64 `mapstructure:"temperature" json:"temperature" yaml:"temperature"`
}
```

在 `server/config/config.go` 添加:

```go
type Server struct {
    // ... 其他配置 ...
    VLLM VLLM `mapstructure:"vllm" json:"vllm" yaml:"vllm"`
}
```

### 3.2 创建服务层

`server/service/vllm/client.go`:

```go
package vllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type VLLMClient struct {
	BaseURL string
	Client  *http.Client
}

func NewVLLMClient() *VLLMClient {
	return &VLLMClient{
		BaseURL: global.GVA_CONFIG.VLLM.BaseURL,
		Client: &http.Client{
			Timeout: time.Duration(global.GVA_CONFIG.VLLM.Timeout) * time.Second,
		},
	}
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *VLLMClient) Chat(ctx context.Context, message string) (string, error) {
	reqBody := ChatRequest{
		Model: global.GVA_CONFIG.VLLM.Model,
		Messages: []Message{
			{Role: "user", Content: message},
		},
		MaxTokens:   global.GVA_CONFIG.VLLM.MaxTokens,
		Temperature: global.GVA_CONFIG.VLLM.Temperature,
	}

	jsonData, _ := json.Marshal(reqBody)
	
	req, err := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/v1/chat/completions",
		bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
```

`server/service/vllm/enter.go`:

```go
package vllm

type ServiceGroup struct {
	VLLMClient
}
```

### 3.3 创建 API 层

`server/api/v1/vllm/vllm.go`:

```go
package vllm

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vllm"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VLLMApi struct{}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (v *VLLMApi) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := vllm.NewVLLMClient()
	reply, err := client.Chat(c.Request.Context(), req.Message)
	if err != nil {
		global.GVA_LOG.Error("vLLM调用失败!", zap.Error(err))
		response.FailWithMessage("调用失败: "+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"reply": reply}, c)
}
```

`server/api/v1/vllm/enter.go`:

```go
package vllm

type ApiGroup struct {
	VLLMApi
}
```

### 3.4 注册路由

`server/router/vllm/vllm.go`:

```go
package vllm

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type VLLMRouter struct{}

func (r *VLLMRouter) InitVLLMRouter(Router *gin.RouterGroup) {
	vllmRouter := Router.Group("vllm")
	vllmApi := v1.ApiGroupApp.VLLMApiGroup
	{
		vllmRouter.POST("chat", vllmApi.Chat)
	}
}
```

`server/router/vllm/enter.go`:

```go
package vllm

type RouterGroup struct {
	VLLMRouter
}
```

### 3.5 注册到主程序

修改 `server/api/v1/enter.go`:

```go
import (
    // ... 其他导入
    "github.com/flipped-aurora/gin-vue-admin/server/api/v1/vllm"
)

type ApiGroup struct {
    // ... 其他
    VLLMApiGroup vllm.ApiGroup
}
```

修改 `server/service/enter.go`:

```go
import (
    // ... 其他导入
    "github.com/flipped-aurora/gin-vue-admin/server/service/vllm"
)

type ServiceGroup struct {
    // ... 其他
    VLLMServiceGroup vllm.ServiceGroup
}
```

修改 `server/router/enter.go`:

```go
import (
    // ... 其他导入
    "github.com/flipped-aurora/gin-vue-admin/server/router/vllm"
)

type RouterGroup struct {
    // ... 其他
    VLLM vllm.RouterGroup
}
```

修改 `server/initialize/router.go`:

```go
func Routers() *gin.Engine {
    // ... 现有代码 ...
    
    vllmRouter := router.RouterGroupApp.VLLM
    
    // 在 PublicGroup 中添加
    {
        // ... 其他路由
        vllmRouter.InitVLLMRouter(PublicGroup) // vLLM 路由
    }
    
    // ... 其他代码
}
```

---

## 第四步：测试

### 4.1 测试 vLLM 直连

```bash
# 测试服务器
curl -X POST http://192.168.1.100:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "Qwen2.5-7B-Instruct-AWQ",
    "messages": [{"role": "user", "content": "你好"}],
    "max_tokens": 100
  }'
```

### 4.2 测试 gin 接口

```bash
# 测试 gin
curl -X POST http://localhost:8888/vllm/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下你自己"}'
```

### 4.3 查看日志

```bash
# CentOS 服务器
docker-compose logs -f vllm

# macOS gin
docker-compose logs -f server
```

---

## 常见问题

### Q1: 连接超时

```bash
# 检查防火墙
sudo firewall-cmd --add-port=8000/tcp --permanent
sudo firewall-cmd --reload

# 或关闭防火墙
sudo systemctl stop firewalld
```

### Q2: 显存不足

```yaml
# 修改 docker-compose.yml
command: >
  --gpu-memory-utilization 0.8  # 降低显存使用
  --max-model-len 2048          # 减小上下文长度
```

### Q3: 模型下载慢

```bash
# 使用国内镜像
export HF_ENDPOINT=https://hf-mirror.com
huggingface-cli download Qwen/Qwen2.5-7B-Instruct-AWQ \
  --local-dir ./models/Qwen2.5-7B-Instruct-AWQ
```

---

## 快速命令参考

```bash
# ===== CentOS 服务器 =====
# 启动
cd ~/vllm-deploy && docker-compose up -d

# 停止
docker-compose stop

# 重启
docker-compose restart

# 查看日志
docker-compose logs -f

# 查看 GPU
nvidia-smi

# ===== macOS 本地 =====
# 重启 gin
docker-compose restart server

# 查看日志
docker-compose logs -f server

# 测试接口
curl -X POST http://localhost:8888/vllm/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}'
```

---

## 性能优化

### 服务器端

```yaml
# docker-compose.yml 优化参数
command: >
  --model /models/Qwen2.5-7B-Instruct-AWQ
  --quantization awq
  --max-model-len 4096
  --gpu-memory-utilization 0.9
  --max-num-seqs 256
  --max-num-batched-tokens 8192
  --enable-prefix-caching
```

### gin 端

```go
// 连接池优化
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 100,
    IdleConnTimeout:     90 * time.Second,
}

Client: &http.Client{
    Transport: transport,
    Timeout:   300 * time.Second,
}
```

---

**完成！** 🎉

现在你可以：
1. CentOS 服务器运行 vLLM + GPU
2. macOS 本地运行 gin
3. 通过 HTTP 通信
4. 调用本地大模型

