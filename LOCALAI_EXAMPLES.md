# 本地AI系统使用示例

这个文档提供了本地AI系统的完整使用示例和最佳实践。

## 目录

- [快速开始](#快速开始)
- [基础对话](#基础对话)
- [RAG检索增强](#rag检索增强)
- [知识库管理](#知识库管理)
- [训练数据收集](#训练数据收集)
- [实战场景](#实战场景)

## 快速开始

### 1. 一键设置（推荐）

```bash
# 运行设置脚本
./setup_localai.sh
```

### 2. 手动设置

```bash
# 安装 Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 启动服务
ollama serve &

# 下载模型
ollama pull qwen2.5:7b

# 创建数据目录
cd server
mkdir -p data/vector data/training knowledge_base

# 启动服务器
go run main.go
```

## 基础对话

### 简单对话

```bash
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好，请介绍一下你自己",
    "use_rag": false
  }'
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "answer": "你好！我是一个AI助手，基于本地大语言模型运行...",
    "sources": [],
    "use_rag": false,
    "session_id": "",
    "time": "2.3s"
  },
  "msg": "请求成功"
}
```

### 编程问题

```bash
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "用Go语言实现一个简单的HTTP服务器",
    "use_rag": false
  }'
```

## RAG检索增强

### 场景1：企业知识库问答

#### 步骤1：导入公司文档

```bash
# 导入公司简介
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "我们公司成立于2020年，总部位于北京。主营业务包括人工智能技术研发、企业数字化转型咨询和云计算服务。公司目前有员工300人，其中研发人员占70%。",
    "metadata": {
      "title": "公司简介",
      "category": "企业信息",
      "date": "2024-01-01"
    },
    "source": "company_docs"
  }'

# 导入产品信息
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "我们的主打产品AI-Chat是一款企业级智能客服系统，支持多轮对话、知识库管理和情感分析。产品采用本地部署方式，确保数据安全。年费为5万元，包含首年的技术支持。",
    "metadata": {
      "title": "产品：AI-Chat",
      "category": "产品",
      "price": "50000"
    }
  }'

# 导入联系方式
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "联系我们：\n客服电话：400-123-4567\n工作时间：周一至周五 9:00-18:00\n邮箱：service@example.com\n地址：北京市朝阳区xxx大厦8层",
    "metadata": {
      "title": "联系方式",
      "category": "联系信息"
    }
  }'
```

#### 步骤2：使用RAG进行问答

```bash
# 问题1：公司信息
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你们公司是什么时候成立的？有多少员工？",
    "use_rag": true,
    "top_k": 3
  }'

# 问题2：产品定价
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "AI-Chat产品的价格是多少？包含哪些服务？",
    "use_rag": true,
    "top_k": 3
  }'

# 问题3：联系方式
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "如何联系你们的客服？",
    "use_rag": true,
    "top_k": 3
  }'
```

### 场景2：技术文档问答

```bash
# 导入技术文档
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "API接口说明：\n1. 认证方式：使用Bearer Token\n2. 请求格式：JSON\n3. 响应格式：JSON\n4. 错误码：400=参数错误，401=未授权，500=服务器错误\n5. 限流：每分钟最多100次请求",
    "metadata": {
      "title": "API文档",
      "version": "v1.0"
    }
  }'

# 查询API信息
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "API的认证方式是什么？限流策略是怎样的？",
    "use_rag": true,
    "top_k": 5
  }'
```

## 知识库管理

### 搜索文档

```bash
# 搜索与"价格"相关的文档
curl -X POST http://localhost:8888/localai/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "产品价格",
    "top_k": 10
  }'
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "results": [
      {
        "document": {
          "id": "doc-uuid-1",
          "content": "我们的主打产品AI-Chat是一款企业级智能客服系统...",
          "metadata": {
            "title": "产品：AI-Chat",
            "category": "产品",
            "price": "50000"
          }
        },
        "score": 0.89,
        "distance": 0.11
      }
    ],
    "total": 1,
    "time": "0.05s"
  },
  "msg": "搜索成功"
}
```

### 删除文档

```bash
# 删除指定文档
curl -X POST http://localhost:8888/localai/document/delete \
  -H "Content-Type: application/json" \
  -d '{
    "id": "doc-uuid-1"
  }'
```

### 批量导入（使用脚本）

```bash
# 创建批量导入脚本
cat > import_docs.sh << 'EOF'
#!/bin/bash

# 从文件导入
for file in knowledge_base/*.txt; do
  content=$(cat "$file")
  filename=$(basename "$file")
  
  curl -X POST http://localhost:8888/localai/ingest \
    -H "Content-Type: application/json" \
    -d "{
      \"content\": \"$content\",
      \"metadata\": {
        \"filename\": \"$filename\",
        \"imported_at\": \"$(date -Iseconds)\"
      }
    }"
  
  echo "Imported: $filename"
done
EOF

chmod +x import_docs.sh
./import_docs.sh
```

## 训练数据收集

### 提交用户反馈

```bash
# 高质量反馈（评分5分）
curl -X POST http://localhost:8888/localai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "question": "你们公司主营业务是什么？",
    "answer": "我们公司主营业务包括人工智能技术研发、企业数字化转型咨询和云计算服务。",
    "rating": 5,
    "feedback": "回答准确且全面"
  }'

# 中等质量反馈（评分3分）
curl -X POST http://localhost:8888/localai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "question": "产品价格是多少？",
    "answer": "产品价格为5万元。",
    "rating": 3,
    "feedback": "回答正确但不够详细，应该说明包含的服务"
  }'

# 低质量反馈（评分1分）
curl -X POST http://localhost:8888/localai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "question": "如何联系客服？",
    "answer": "不知道。",
    "rating": 1,
    "feedback": "完全没有回答问题"
  }'
```

### 查看统计信息

```bash
curl http://localhost:8888/localai/stats
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "total_documents": 3,
    "total_examples": 15,
    "model_info": {
      "name": "qwen2.5:7b",
      "backend": "ollama",
      "status": "running"
    },
    "vector_store_status": "healthy"
  },
  "msg": "查询成功"
}
```

### 开始训练

```bash
# 确保有足够的训练样本（默认需要100个）
# 然后启动训练

curl -X POST http://localhost:8888/localai/train/start \
  -H "Content-Type: application/json" \
  -d '{
    "method": "lora",
    "epochs": 3,
    "batch_size": 4,
    "learning_rate": 0.0002,
    "use_gpu": false
  }'
```

### 监控训练进度

```bash
# 轮询训练状态
while true; do
  curl http://localhost:8888/localai/train/status
  sleep 5
done
```

## 实战场景

### 场景1：智能客服系统

```bash
# 1. 导入常见问题库
cat > faq.json << 'EOF'
[
  {
    "question": "如何重置密码？",
    "answer": "点击登录页面的"忘记密码"，输入注册邮箱，系统会发送重置链接。"
  },
  {
    "question": "支持哪些支付方式？",
    "answer": "支持支付宝、微信支付、银行卡支付和企业对公转账。"
  },
  {
    "question": "退款需要多久？",
    "answer": "退款申请提交后，一般3-7个工作日到账。"
  }
]
EOF

# 导入FAQ
jq -c '.[]' faq.json | while read item; do
  question=$(echo "$item" | jq -r '.question')
  answer=$(echo "$item" | jq -r '.answer')
  
  curl -X POST http://localhost:8888/localai/ingest \
    -H "Content-Type: application/json" \
    -d "{
      \"content\": \"问题：$question\n答案：$answer\",
      \"metadata\": {
        \"type\": \"faq\",
        \"question\": \"$question\"
      }
    }"
done

# 2. 客户咨询
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我忘记密码了，怎么办？",
    "use_rag": true,
    "top_k": 3
  }'
```

### 场景2：内部知识管理

```bash
# 1. 导入会议纪要
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "2024年1月产品会议纪要：\n1. 决定在Q2推出AI-Chat 2.0版本\n2. 新增多语言支持功能\n3. 优化响应速度，目标1秒内\n4. 预算增加30%用于研发",
    "metadata": {
      "type": "meeting",
      "date": "2024-01-15",
      "participants": "产品团队"
    }
  }'

# 2. 查询历史决策
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我们什么时候决定推出2.0版本的？预算是多少？",
    "use_rag": true,
    "top_k": 5
  }'
```

### 场景3：代码助手

```bash
# 1. 导入项目文档
curl -X POST http://localhost:8888/localai/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "content": "项目架构说明：\n- 前端：Vue 3 + TypeScript\n- 后端：Go + Gin框架\n- 数据库：MySQL 8.0\n- 缓存：Redis\n- 部署：Docker + K8s\n\n目录结构：\n- /server：后端代码\n- /web：前端代码\n- /deploy：部署配置",
    "metadata": {
      "type": "documentation",
      "topic": "architecture"
    }
  }'

# 2. 询问项目信息
curl -X POST http://localhost:8888/localai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "这个项目用什么技术栈？前端和后端分别在哪个目录？",
    "use_rag": true,
    "top_k": 3
  }'
```

## 健康检查

```bash
# 检查系统状态
curl http://localhost:8888/localai/health

# 检查Ollama状态
curl http://localhost:11434/api/tags

# 查看日志
tail -f /tmp/ollama.log
```

## Python集成示例

```python
import requests
import json

class LocalAIClient:
    def __init__(self, base_url="http://localhost:8888"):
        self.base_url = base_url
    
    def chat(self, message, use_rag=False, top_k=5):
        """发送聊天消息"""
        url = f"{self.base_url}/localai/chat"
        data = {
            "message": message,
            "use_rag": use_rag,
            "top_k": top_k
        }
        
        response = requests.post(url, json=data)
        return response.json()
    
    def ingest(self, content, metadata=None):
        """导入文档"""
        url = f"{self.base_url}/localai/ingest"
        data = {
            "content": content,
            "metadata": metadata or {}
        }
        
        response = requests.post(url, json=data)
        return response.json()
    
    def search(self, query, top_k=10):
        """搜索文档"""
        url = f"{self.base_url}/localai/search"
        data = {
            "query": query,
            "top_k": top_k
        }
        
        response = requests.post(url, json=data)
        return response.json()
    
    def feedback(self, question, answer, rating, feedback=""):
        """提交反馈"""
        url = f"{self.base_url}/localai/feedback"
        data = {
            "question": question,
            "answer": answer,
            "rating": rating,
            "feedback": feedback
        }
        
        response = requests.post(url, json=data)
        return response.json()

# 使用示例
client = LocalAIClient()

# 1. 导入文档
result = client.ingest(
    content="公司成立于2020年，员工300人。",
    metadata={"title": "公司简介"}
)
print("导入结果:", result)

# 2. RAG对话
result = client.chat(
    message="公司有多少员工？",
    use_rag=True,
    top_k=5
)
print("回答:", result['data']['answer'])

# 3. 搜索文档
results = client.search(query="员工", top_k=5)
print("搜索结果:", results)

# 4. 提交反馈
client.feedback(
    question="公司有多少员工？",
    answer="公司有300名员工。",
    rating=5,
    feedback="回答准确"
)
```

## JavaScript/Node.js集成示例

```javascript
const axios = require('axios');

class LocalAIClient {
  constructor(baseURL = 'http://localhost:8888') {
    this.baseURL = baseURL;
    this.client = axios.create({ baseURL });
  }

  async chat(message, useRAG = false, topK = 5) {
    const response = await this.client.post('/localai/chat', {
      message,
      use_rag: useRAG,
      top_k: topK
    });
    return response.data;
  }

  async ingest(content, metadata = {}) {
    const response = await this.client.post('/localai/ingest', {
      content,
      metadata
    });
    return response.data;
  }

  async search(query, topK = 10) {
    const response = await this.client.post('/localai/search', {
      query,
      top_k: topK
    });
    return response.data;
  }

  async feedback(question, answer, rating, feedback = '') {
    const response = await this.client.post('/localai/feedback', {
      question,
      answer,
      rating,
      feedback
    });
    return response.data;
  }

  async getStats() {
    const response = await this.client.get('/localai/stats');
    return response.data;
  }
}

// 使用示例
(async () => {
  const client = new LocalAIClient();

  // 导入文档
  await client.ingest(
    '公司成立于2020年，员工300人。',
    { title: '公司简介' }
  );

  // RAG对话
  const result = await client.chat('公司有多少员工？', true, 5);
  console.log('回答:', result.data.answer);

  // 获取统计
  const stats = await client.getStats();
  console.log('统计:', stats.data);
})();
```

## 性能优化建议

1. **批量导入优化**
   - 使用并发请求
   - 分批处理大文件
   - 添加重试机制

2. **查询优化**
   - 减少 `top_k` 值
   - 使用缓存机制
   - 预处理常见问题

3. **模型优化**
   - 使用量化模型
   - 调整 max_tokens
   - 启用GPU加速

更多信息请参考 `server/api/v1/localai/README.md`

