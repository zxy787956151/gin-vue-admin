# 通义千问 API 接口

## 功能说明

这是一个调用阿里云通义千问大模型的接口，可以实现对话功能。

## 配置说明

在使用之前，需要在 `server/config.yaml` 中配置通义千问的 API Key：

```yaml
# 通义千问配置
qianwen:
    api-key: "your-api-key-here" # 请在 https://dashscope.console.aliyun.com/ 获取 API Key
    model: "qwen-turbo" # 使用的模型，qwen-turbo 为免费模型
```

### 获取 API Key 步骤：

1. 访问 [通义千问控制台](https://dashscope.console.aliyun.com/)
2. 登录后，在左侧菜单选择 "API-KEY 管理"
3. 创建或复制现有的 API Key
4. 将 API Key 填入配置文件中

## API 接口

### POST /qianwen/chat

与通义千问进行对话

#### 请求参数

```json
{
    "message": "你好，请介绍一下你自己"
}
```

#### 响应示例

```json
{
    "code": 0,
    "data": {
        "reply": "你好！我是通义千问，由阿里云开发的AI助手。我可以帮助你回答问题、提供信息、进行对话等。有什么我可以帮助你的吗？"
    },
    "msg": "请求成功"
}
```

## 使用示例

### curl 命令

```bash
curl -X POST http://localhost:8888/qianwen/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}'
```

### JavaScript (Fetch API)

```javascript
fetch('http://localhost:8888/qianwen/chat', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        message: '你好，请介绍一下你自己'
    })
})
.then(response => response.json())
.then(data => console.log(data.data.reply))
.catch(error => console.error('Error:', error));
```

### Python (requests)

```python
import requests

url = "http://localhost:8888/qianwen/chat"
payload = {
    "message": "你好，请介绍一下你自己"
}

response = requests.post(url, json=payload)
result = response.json()
print(result['data']['reply'])
```

## 支持的模型

- `qwen-turbo`: 通义千问超大规模语言模型（免费）
- `qwen-plus`: 通义千问增强版（付费）
- `qwen-max`: 通义千问最强版（付费）

更多模型请参考：https://help.aliyun.com/zh/dashscope/developer-reference/model-square

## 注意事项

1. 此接口为公开接口，不需要登录验证
2. 请妥善保管你的 API Key，不要泄露给他人
3. 免费模型 `qwen-turbo` 有调用次数限制，请查看阿里云控制台了解详情
4. 建议在生产环境中添加请求频率限制和鉴权机制

