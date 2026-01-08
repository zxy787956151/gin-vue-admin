# Ollama 本地 AI 部署

> 使用 Ollama 运行免费的本地大模型

## 快速开始

### 1. 启动 Ollama 容器

```bash
docker-compose -f docker-compose-ollama.yml up -d
```

### 2. 下载模型

```bash
# 推荐：千问 7B 模型（中文效果好）
docker exec ollama-qwen ollama pull qwen2.5:7b

# 其他选择
docker exec ollama-qwen ollama pull qwen2.5:3b        # 更快，占用更少内存
docker exec ollama-qwen ollama pull qwen2.5-coder:7b  # 代码生成专用
```

### 3. 配置 gin 服务

确认 `server/config.yaml` 配置正确：

```yaml
local-ai:
  llm:
    backend: ollama
    base-url: http://ollama-qwen:11434  # 容器模式
    # base-url: http://localhost:11434  # 本地模式
    model: qwen2.5:7b
    timeout: 300
    max-tokens: 2000
```

### 4. 启动 gin 服务

**容器内：**
```bash
docker exec -it gin-backend sh
cd /app
go run main.go
```

**或者宿主机：**
```bash
cd server
go run main.go
```

## 测试

### 健康检查
```bash
curl http://localhost:8888/localai/health
```

### AI 对话
```bash
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好", "use_rag": false}'
```

## 常用命令

```bash
# 查看容器状态
docker ps | grep ollama

# 查看已安装模型
docker exec ollama-qwen ollama list

# 删除模型
docker exec ollama-qwen ollama rm qwen2.5:7b

# 重启容器
docker restart ollama-qwen

# 停止服务
docker-compose -f docker-compose-ollama.yml down
```

## 网络配置

### 如果 gin 也在容器里

确保两个容器在同一网络：

```bash
# 连接到 gin 网络
docker network connect gin_default ollama-qwen

# 验证连通性
docker exec gin-backend ping ollama-qwen
```

配置使用容器名：
```yaml
base-url: http://ollama-qwen:11434
```

### 如果 gin 在宿主机

配置使用 localhost：
```yaml
base-url: http://localhost:11434
```

## 模型推荐

| 模型 | 大小 | 内存需求 | 适用场景 |
|------|------|---------|---------|
| qwen2.5:3b | 2GB | 4-8GB | 快速响应 |
| qwen2.5:7b | 5GB | 8-16GB | **推荐**，最佳平衡 |
| qwen2.5:14b | 9GB | 16-32GB | 高质量输出 |
| qwen2.5-coder:7b | 5GB | 8-16GB | 代码生成 |

## 故障排查

### 1. 容器无法启动
```bash
docker logs ollama-qwen
docker restart ollama-qwen
```

### 2. gin 无法连接 ollama
```bash
# 检查网络
docker network ls
docker network inspect gin_default

# 测试连通性
docker exec gin-backend curl http://ollama-qwen:11434/api/version
```

### 3. 模型下载失败
```bash
# 检查磁盘空间
df -h

# 手动进入容器下载
docker exec -it ollama-qwen sh
ollama pull qwen2.5:7b
```

### 4. 服务报错"模型未配置"
- 确保配置文件路径正确
- 重启 gin 服务让配置生效
- 检查 base-url 是否正确（容器名 vs localhost）

## API 接口

### POST /localai/chat
```json
{
  "message": "你的问题",
  "use_rag": false
}
```

### GET /localai/health
健康检查

### POST /localai/ingest
上传文档到知识库

### POST /localai/search
搜索知识库

## 参考

- Ollama 官网: https://ollama.com
- 模型库: https://ollama.com/library
- API 文档: https://github.com/ollama/ollama/blob/main/docs/api.md

