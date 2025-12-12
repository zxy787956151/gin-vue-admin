# 本地AI系统 - 可训练的RAG模型

## 功能特性

这是一个完整的本地AI系统，集成了以下功能：

✅ **本地大模型推理** - 支持 Ollama、llama.cpp、vLLM 等多种后端  
✅ **RAG检索增强** - 向量存储 + 语义搜索  
✅ **知识库管理** - 文档导入、搜索、删除  
✅ **训练数据收集** - 自动收集用户反馈用于训练  
✅ **模型训练** - 支持 LoRA 微调（需外部工具）  
✅ **完全本地化** - 所有数据和处理都在本地完成  


## 与通义千问接口的区别

| 特性 | 通义千问 (`/qianwen/chat`) | 本地AI (`/localai/*`) |
|------|---------------------------|----------------------|
| 部署方式 | 云端API | 本地部署 |
| 数据隐私 | 数据上传到云端 | 数据完全本地 |
| 成本 | 按调用次数收费 | 硬件成本（一次性） |
| RAG支持 | ❌ | ✅ |
| 知识库 | ❌ | ✅ |
| 可训练 | ❌ | ✅ |
| 响应速度 | 快（云端算力） | 取决于本地硬件 |

## 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                     本地AI系统                            │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │ 用户输入  │  │ RAG检索  │  │ 本地LLM  │              │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘              │
│       │             │              │                     │
│       ▼             ▼              ▼                     │
│  ┌─────────────────────────────────────┐               │
│  │          RAG引擎 (rag.go)           │               │
│  ├─────────────────────────────────────┤               │
│  │  • 检索相关文档                      │               │
│  │  • 构建增强提示                      │               │
│  │  • 调用本地模型生成回答              │               │
│  │  • 收集训练数据                      │               │
│  └─────────────────────────────────────┘               │
│       │             │              │                     │
│       ▼             ▼              ▼                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │向量存储  │  │本地LLM   │  │训练服务  │              │
│  │(vector)  │  │(llm)     │  │(training)│              │
│  └──────────┘  └──────────┘  └──────────┘              │
└─────────────────────────────────────────────────────────┘
```

## 快速开始

### 1. 安装 Ollama（推荐）

```bash
# macOS/Linux
curl -fsSL https://ollama.com/install.sh | sh

# Windows
# 访问 https://ollama.com/download 下载安装包

# 启动 Ollama 服务
ollama serve

# 下载模型（在新终端）
ollama pull qwen2.5:7b    # 或其他模型
ollama pull llama3.1:8b
```

### 2. 配置系统

编辑 `server/config.yaml` 文件：

```yaml
local-ai:
  llm:
    backend: "ollama"
    base-url: "http://localhost:11434"
    model: "qwen2.5:7b"
    timeout: 300
    max-tokens: 2000
  
  vector-store:
    type: "local"
    data-path: "./data/vector"
  
  training:
    enabled: true
    data-path: "./data/training"
    auto-train: false
    min-examples: 100
```

### 3. 创建数据目录

```bash
cd server
mkdir -p data/vector data/training knowledge_base
```

### 4. 启动服务

```bash
cd server
go run main.go
```

## API 接口说明

### 1. 聊天接口（支持RAG）

**接口**: `POST /localai/chat`

**请求参数**:
```json
{
  "message": "什么是Go语言？",
  "use_rag": true,
  "top_k": 5,
  "stream": false,
  "session_id": "session-123"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "answer": "Go语言是由Google开发的...",
    "sources": [
      {
        "document": {
          "id": "doc-1",
          "content": "Go语言相关文档内容...",
          "metadata": {}
        },
        "score": 0.85,
        "distance": 0.15
      }
    ],
    "use_rag": true,
    "session_id": "session-123",
    "time": "1.5s"
  },
  "msg": "请求成功"
}
```

### 2. 导入文档到知识库

**接口**: `POST /localai/ingest`

**请求**:
```json
{
  "content": "Go语言是一种开源编程语言，由Google开发。它具有静态类型、编译型等特点...",
  "metadata": {
    "title": "Go语言简介",
    "source": "文档",
    "category": "编程语言"
  },
  "source": "manual_upload"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "doc_id": "uuid-xxx",
    "status": "success",
    "message": "文档已成功添加到知识库",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "msg": "文档导入成功"
}
```

### 3. 搜索文档

**接口**: `POST /localai/search`

**请求**:
```json
{
  "query": "Go语言的特点",
  "top_k": 10
}
```

### 4. 提交反馈（用于训练）

**接口**: `POST /localai/feedback`

**请求**:
```json
{
  "question": "什么是Go语言？",
  "answer": "Go语言是由Google开发的...",
  "rating": 4.5,
  "feedback": "回答很详细，但可以更简洁"
}
```

### 5. 开始训练

**接口**: `POST /localai/train/start`

**请求**:
```json
{
  "method": "lora",
  "epochs": 3,
  "batch_size": 4,
  "learning_rate": 0.0002,
  "use_gpu": true
}
```

### 6. 查询训练状态

**接口**: `GET /localai/train/status`

**响应**:
```json
{
  "code": 0,
  "data": {
    "status": "training",
    "progress": 45.5,
    "message": "训练进行中... 45.5%",
    "start_time": "2024-01-01T12:00:00Z",
    "metrics": {
      "loss": 0.35,
      "accuracy": 0.87
    }
  }
}
```

### 7. 获取统计信息

**接口**: `GET /localai/stats`

**响应**:
```json
{
  "code": 0,
  "data": {
    "total_documents": 150,
    "total_examples": 89,
    "model_info": {
      "name": "qwen2.5:7b",
      "backend": "ollama",
      "status": "running"
    },
    "vector_store_status": "healthy"
  }
}
```

### 8. 健康检查

**接口**: `GET /localai/health`

### 9. 删除文档

**接口**: `POST /localai/document/delete`

**请求**:
```json
{
  "id": "doc-uuid-xxx"
}
```

## 使用示例

### 基础对话（不使用RAG）

```bash
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，介绍一下你自己",
    "use_rag": false
  }'
```

### RAG增强对话

```bash
# 1. 先导入一些文档
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "公司成立于2020年，主营业务是人工智能技术研发。",
    "metadata": {"title": "公司简介"}
  }'

# 2. 使用RAG进行对话
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "公司是什么时候成立的？",
    "use_rag": true,
    "top_k": 5
  }'
```

### 收集反馈并训练

```bash
# 1. 提交用户反馈
curl -X POST http://localhost:8888/localai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "question": "公司主营业务是什么？",
    "answer": "公司主营业务是人工智能技术研发。",
    "rating": 5,
    "feedback": "回答准确"
  }'

# 2. 收集足够样本后，启动训练
curl -X POST http://localhost:8888/localai/train/start \
  -H "Content-Type: application/json" \
  -d '{
    "method": "lora",
    "epochs": 3
  }'

# 3. 查询训练状态
curl http://localhost:8888/localai/train/status
```

## 系统要求

### 最低配置
- **CPU**: 4核心以上
- **内存**: 16GB RAM
- **存储**: 20GB 可用空间
- **操作系统**: Linux / macOS / Windows

### 推荐配置（使用GPU）
- **GPU**: NVIDIA GPU with 8GB+ VRAM
- **内存**: 32GB RAM
- **CUDA**: 11.8+

## 支持的模型

### 通过 Ollama 运行的模型

```bash
# 中文模型
ollama pull qwen2.5:7b        # 阿里通义千问
ollama pull qwen2.5:14b       # 更大版本
ollama pull chatglm3:6b       # 清华智谱

# 英文模型
ollama pull llama3.1:8b       # Meta Llama
ollama pull mistral:7b        # Mistral AI
ollama pull phi3:mini         # Microsoft Phi

# 查看已安装的模型
ollama list
```

## 训练模型（高级功能）

### 准备训练环境

```bash
# 安装 Python 依赖
pip install torch transformers peft accelerate datasets

# 或使用 axolotl 训练框架
git clone https://github.com/OpenAccess-AI-Collective/axolotl
cd axolotl
pip install -e .
```

### 训练流程

1. **收集训练数据** - 通过反馈接口收集高质量的问答对
2. **准备数据集** - 系统自动生成 `./data/training/train.json`
3. **启动训练** - 调用训练接口或使用外部训练脚本
4. **评估模型** - 在测试集上评估模型效果
5. **部署模型** - 将训练好的模型部署到 Ollama

### 使用训练好的模型

```bash
# 1. 创建 Modelfile
cat > Modelfile <<EOF
FROM qwen2.5:7b
ADAPTER ./models/finetuned/adapter_model.bin
PARAMETER temperature 0.7
EOF

# 2. 创建自定义模型
ollama create my-custom-model -f Modelfile

# 3. 在配置中使用
# config.yaml:
# local-ai:
#   llm:
#     model: "my-custom-model"
```

## 性能优化

### 1. 使用量化模型

```bash
# Q4量化（推荐）
ollama pull qwen2.5:7b-q4_K_M

# Q8量化（更高精度）
ollama pull qwen2.5:7b-q8_0
```

### 2. 调整配置

```yaml
local-ai:
  llm:
    max-tokens: 1000  # 减少生成长度
    timeout: 120      # 减少超时时间
```

### 3. 向量存储优化

对于大规模知识库，建议使用专业的向量数据库：

```yaml
local-ai:
  vector-store:
    type: "qdrant"    # 或 chroma
    host: "localhost:6333"
```

## 故障排查

### 1. 无法连接到 Ollama

```bash
# 检查服务是否运行
ollama list

# 重启服务
ollama serve
```

### 2. 模型响应慢

- 使用量化模型（Q4/Q5）
- 减少 max-tokens
- 升级硬件（增加内存/使用GPU）

### 3. 内存不足

```bash
# 使用更小的模型
ollama pull qwen2.5:3b
ollama pull phi3:mini
```

## 路线图

- [ ] 支持流式输出
- [ ] 多轮对话上下文管理
- [ ] 文档批量导入
- [ ] PDF/Word 文件解析
- [ ] 向量数据库集成（Qdrant/Chroma）
- [ ] 嵌入模型服务
- [ ] 自动化训练流程
- [ ] Web UI 管理界面
- [ ] 性能监控和日志

## 最佳实践

1. **知识库管理**
   - 定期更新知识库内容
   - 删除过时或低质量的文档
   - 为文档添加有意义的 metadata

2. **训练数据收集**
   - 只保留评分 >= 3 的样本
   - 定期审核自动收集的数据
   - 手动标注高质量样本

3. **模型选择**
   - 开发环境：使用小模型（3B/7B）
   - 生产环境：根据需求选择（7B/14B）
   - 资源受限：使用量化版本

4. **安全建议**
   - 设置访问权限控制
   - 不要在知识库中存储敏感信息
   - 定期备份数据

## 参考资源

- [Ollama 官方文档](https://ollama.com/docs)
- [Qwen2.5 模型介绍](https://github.com/QwenLM/Qwen2.5)
- [LoRA 微调教程](https://github.com/huggingface/peft)
- [RAG 最佳实践](https://www.pinecone.io/learn/retrieval-augmented-generation/)

