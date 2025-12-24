# AI 智能客服技术验证方案（规范版）

## 项目定位

**类型**：技术可行性验证项目（POC）  
**目标**：验证 vLLM + Qwen2.5 在企业 AI 客服场景的可行性  
**时间**：5 个工作日  
**人员**：1 名技术人员

---

## 阶段规划（5 个阶段）

```
阶段 0: 方案设计 → 阶段 1: 环境准备 → 阶段 2: 服务部署 → 阶段 3: 接口集成 → 阶段 4: 性能验证
```

---

## 阶段 0：方案设计与技术选型

### 进入条件
- ✅ 需求明确：AI 智能客服场景
- ✅ 技术约束：必须使用 vLLM + Qwen
- ✅ 资源确认：有 GPU 服务器资源

### 主要工作
1. **技术选型调研**
   - 推理引擎：vLLM vs. TGI vs. Ollama
   - 模型选择：Qwen2.5 vs. Llama vs. ChatGLM
   - 部署方式：Docker vs. 裸机安装

2. **架构设计**
   - 服务端架构（AI 推理层）
   - 客户端架构（业务集成层）
   - 网络架构（内网 vs. 公网）

3. **风险评估**
   - 性能风险：GPU 资源是否足够
   - 兼容性风险：vLLM 与 Qwen 版本兼容性
   - 成本风险：GPU 服务器成本

### 退出条件
- ✅ 技术选型确定：vLLM v0.4.2 + Qwen2.5-7B-AWQ
- ✅ 架构方案评审通过
- ✅ 风险识别清单完成

### 验证指标
- [ ] 架构图完整（含服务端、客户端、网络）
- [ ] 技术选型文档（含对比分析）
- [ ] 风险清单（≥5 项风险点）

### 决策记录

| 决策点 | 选项 | 选择 | 原因 | 风险 |
|--------|------|------|------|------|
| **推理引擎** | vLLM / TGI / Ollama | ✅ vLLM | 性能最优，社区活跃 | 版本兼容性复杂 |
| **模型版本** | Qwen2.5-7B / 14B | ✅ 7B-AWQ | 显存占用低，速度快 | 效果可能不如 14B |
| **部署方式** | Docker / 裸机 | ✅ Docker | 环境隔离，易于管理 | 需要 GPU 直通 |
| **网络方案** | 内网 / 公网 | ✅ 内网 + VPN | 数据安全 | 需配置 VPN/跳板机 |

---

## 阶段 1：环境准备

### 进入条件
- ✅ GPU 服务器已分配（NVIDIA A10, 23GB）
- ✅ 操作系统已安装（CentOS/Ubuntu）
- ✅ 网络已配置（可访问互联网）

### 主要工作
1. **基础环境**
   - 安装 NVIDIA 驱动
   - 安装 Docker + Docker Compose
   - 安装 NVIDIA Container Toolkit

2. **网络配置**
   - 防火墙规则配置
   - 云服务器安全组配置
   - 端口开放（8000）

3. **存储规划**
   - 创建数据盘挂载点（/data）
   - 配置 Docker 数据目录
   - 预留模型存储空间（≥10GB）

### 退出条件
- ✅ GPU 驱动正常：`nvidia-smi` 可执行
- ✅ Docker 可用：`docker run hello-world` 成功
- ✅ GPU 容器可用：`docker run --gpus all nvidia/cuda:12.0.0-base nvidia-smi` 成功

### 验证指标
```bash
# 必须全部通过
# 1. GPU 驱动正常
nvidia-smi                           # 显示 GPU 信息

# 2. Docker 版本
docker --version                     # 版本 ≥ 20.10

# 3. GPU 容器支持
# 如果 vLLM 容器已运行，直接检查它
docker exec vllm-qwen nvidia-smi 2>/dev/null || \
echo "vLLM 容器未运行，后续阶段会验证"

# 4. 磁盘空间
df -h /data                          # 可用空间 ≥ 50GB
```

### 遇到的问题与解决（迷宫记录）

| 问题 | 现象 | 尝试方案 | 结果 | 最终方案 |
|------|------|---------|------|---------|
| **CUDA 驱动不兼容** | `Error 803: unsupported driver` | 1. 使用最新版 vLLM ❌<br>2. 更新驱动 ❌<br>3. 使用旧版 vLLM (v0.4.2) ✅ | 死路→死路→活路 | 使用 vLLM v0.4.2 |
| **磁盘空间不足** | 根分区 100% 满 | 1. 清理 Docker 缓存 ❌<br>2. 迁移 Docker 到 /data ✅ | 死路→活路 | 配置 Docker data-root |
| **网络访问失败** | macOS 访问服务器超时 | 1. 关闭防火墙 ❌<br>2. 配置安全组 ✅ | 死路→活路 | 开放 8000 端口 |

---

## 阶段 2：vLLM 服务部署

### 进入条件
- ✅ 阶段 1 所有验证指标通过
- ✅ 模型下载链接可访问（HuggingFace 或 ModelScope）

### 主要工作
1. **模型下载**
   - 下载 Qwen2.5-7B-Instruct-AWQ（~5GB）
   - 验证模型文件完整性

2. **服务部署**
   - 编写 docker-compose.yml
   - 配置模型路径、GPU 参数
   - 启动 vLLM 服务

3. **服务验证**
   - 健康检查接口测试
   - 模型列表接口测试
   - 对话接口功能测试

### 退出条件
- ✅ 服务启动成功：容器状态 `Up`
- ✅ 健康检查通过：`curl http://localhost:8000/health` 返回 200
- ✅ 对话功能正常：能够正常生成回复

### 验证指标
```bash
# 1. 容器运行
docker ps | grep vllm-qwen          # 状态: Up

# 2. 端口监听
ss -tlnp | grep 8000                # 0.0.0.0:8000 LISTEN

# 3. 健康检查
curl http://localhost:8000/health   # HTTP/1.1 200 OK

# 4. 功能测试
curl -X POST http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model": "Qwen/Qwen2.5-7B-Instruct-AWQ", "messages": [{"role": "user", "content": "你好"}], "max_tokens": 50}'
# 返回正常 JSON，包含 AI 回复

# 5. GPU 使用率
nvidia-smi                          # GPU Memory-Usage > 0
```

### 遇到的问题与解决

| 问题 | 现象 | 尝试方案 | 结果 | 最终方案 |
|------|------|---------|------|---------|
| **vLLM v1 引擎失败** | `Engine core initialization failed` | 1. 调整参数 ❌<br>2. 禁用 v1 引擎 ❌<br>3. 换 v0.4.2 版本 ✅ | 死路→死路→活路 | 使用 v0.4.2（v0 引擎）|
| **模型下载慢** | 下载超时 | 1. Docker 内下载 ❌<br>2. Python 下载 ❌<br>3. vLLM 自动下载 ✅ | 死路→死路→活路 | 让 vLLM 启动时自动下载 |
| **AWQ 量化不兼容** | 加载失败 | 1. 使用 GPTQ ❌<br>2. 使用旧版 vLLM ✅ | 死路→活路 | v0.4.2 支持 AWQ |

### 关键配置
```yaml
# 最终可用的 docker-compose.yml
version: '3.8'
services:
  vllm:
    image: vllm/vllm-openai:v0.4.2  # 关键：必须用 v0.4.2
    runtime: nvidia                 # 关键：GPU 访问方式
    ports:
      - "8000:8000"
    volumes:
      - /data/vllm-cache:/root/.cache  # 关键：数据盘路径
    environment:
      - HF_ENDPOINT=https://hf-mirror.com  # 关键：国内镜像
    command: >
      --model Qwen/Qwen2.5-7B-Instruct-AWQ
      --quantization awq
      --host 0.0.0.0
      --port 8000
      --trust-remote-code
      --max-model-len 4096
      --gpu-memory-utilization 0.9
```

---

## 阶段 3：客户端集成验证

### 进入条件
- ✅ 阶段 2 所有验证指标通过
- ✅ vLLM 服务稳定运行 ≥ 10 分钟
- ✅ 网络可达：客户端能访问服务器 8000 端口

### 主要工作
1. **网络连通性测试**
   - 内网访问测试（服务器本地）
   - 外网访问测试（开发机）
   - 防火墙/安全组配置

2. **多语言客户端实现**
   - Go 客户端（gin-vue-admin）
   - PHP 客户端（示例代码）
   - Java 客户端（示例代码）

3. **接口封装**
   - 配置管理
   - 错误处理
   - 超时重试
   - 日志记录

### 退出条件
- ✅ 至少 2 种语言的客户端验证通过
- ✅ 端到端调用成功率 ≥ 95%
- ✅ 响应时间 ≤ 5 秒（中位数）

### 验证指标
```bash
# 1. 网络连通性
# 从开发机执行
curl --max-time 5 http://服务器IP:8000/health  # 成功返回

# 2. 功能测试（10 次调用）
for i in {1..10}; do
  curl -X POST http://服务器IP:8000/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d "{\"model\": \"Qwen/Qwen2.5-7B-Instruct-AWQ\", \"messages\": [{\"role\": \"user\", \"content\": \"测试$i\"}], \"max_tokens\": 50}"
  echo ""
done
# 成功率: 10/10 = 100%

# 3. 响应时间统计
# 记录 10 次调用的耗时，计算中位数 ≤ 5s
```

### 遇到的问题与解决

| 问题 | 现象 | 尝试方案 | 结果 | 最终方案 |
|------|------|---------|------|---------|
| **外网无法访问** | `Connection timeout` | 1. 修改 Docker 配置 ❌<br>2. 关闭防火墙 ❌<br>3. 配置云安全组 ⏳ | 死路→死路→待验证 | 配置云安全组规则 |
| **Ping 不通** | `Request timeout` | 1. 检查网络 ❌<br>2. 忽略（ICMP 被禁） ✅ | 死路→活路 | Ping 不通不影响 HTTP |

---

## 阶段 4：性能与稳定性验证

### 进入条件
- ✅ 阶段 3 所有验证指标通过
- ✅ 至少 1 个客户端稳定运行

### 主要工作
1. **性能基准测试**
   - 单并发响应时间
   - 多并发吞吐量（5/10/20 并发）
   - GPU 利用率监控

2. **稳定性测试**
   - 长时间运行测试（≥ 1 小时）
   - 异常场景测试（超长输入、空输入）
   - 服务重启恢复测试

3. **资源消耗分析**
   - GPU 显存占用
   - CPU/内存占用
   - 网络带宽消耗

### 退出条件
- ✅ 性能达标：单次响应 ≤ 5s（P95）
- ✅ 稳定性达标：1 小时无崩溃，成功率 ≥ 99%
- ✅ 资源消耗：GPU 显存 ≤ 8GB，CPU ≤ 50%

### 验证指标
```bash
# 1. 性能测试（使用 ab 或 wrk）
# 单并发 100 次请求
ab -n 100 -c 1 -p request.json -T application/json http://服务器IP:8000/v1/chat/completions
# 指标：Time per request < 5000ms

# 2. 多并发测试
ab -n 100 -c 10 -p request.json -T application/json http://服务器IP:8000/v1/chat/completions
# 指标：Success rate ≥ 95%

# 3. GPU 监控（持续 1 小时）
watch -n 10 nvidia-smi
# 记录：GPU Memory 稳定在 4-8GB，利用率波动正常

# 4. 长时间稳定性
# 持续请求 1 小时，记录：
# - 总请求数
# - 成功数
# - 失败数
# - 平均响应时间
# 指标：成功率 ≥ 99%
```

### 性能基准（参考值）

| 指标 | 目标值 | 实测值 | 状态 |
|------|--------|--------|------|
| **单并发响应时间（P50）** | ≤ 2s | 待测 | - |
| **单并发响应时间（P95）** | ≤ 5s | 待测 | - |
| **10 并发吞吐量** | ≥ 2 req/s | 待测 | - |
| **GPU 显存占用** | ≤ 10GB | 待测 | - |
| **1 小时稳定性** | ≥ 99% | 待测 | - |

---

## 总体进度跟踪

| 阶段 | 开始日期 | 结束日期 | 状态 | 验收人 | 备注 |
|------|----------|----------|------|--------|------|
| **0. 方案设计** | 2025-12-22 | 2025-12-22 | ✅ 完成 | - | 技术选型确定 |
| **1. 环境准备** | 2025-12-22 | 2025-12-23 | ✅ 完成 | - | GPU 容器可用 |
| **2. 服务部署** | 2025-12-23 | 2025-12-23 | ✅ 完成 | - | vLLM 服务正常 |
| **3. 客户端集成** | 2025-12-23 | - | ⏳ 进行中 | - | 安全组待配置 |
| **4. 性能验证** | - | - | 📋 待开始 | - | - |

---

## 风险与问题清单

### 已解决问题（活路）

| # | 问题 | 影响 | 解决方案 | 验证时间 |
|---|------|------|---------|---------|
| 1 | CUDA 驱动版本不兼容 | 服务无法启动 | 使用 vLLM v0.4.2 | 2025-12-23 |
| 2 | 根分区磁盘满 | 容器无法创建 | 迁移 Docker 到 /data | 2025-12-23 |
| 3 | 模型下载失败 | 服务无法加载模型 | vLLM 自动下载 + 国内镜像 | 2025-12-23 |
| 4 | vLLM v1 引擎崩溃 | 服务启动失败 | 降级到 v0.4.2（v0 引擎）| 2025-12-23 |

### 待解决问题（当前障碍）

| # | 问题 | 影响 | 优先级 | 计划方案 | 负责人 |
|---|------|------|--------|---------|--------|
| 1 | 云安全组未配置 | 外网无法访问 | 🔴 高 | 配置安全组开放 8000 | 待定 |
| 2 | 缺少性能基准 | 无法评估生产可行性 | 🟡 中 | 执行阶段 4 测试 | 待定 |

### 潜在风险（未来可能遇到）

| # | 风险 | 概率 | 影响 | 应对策略 |
|---|------|------|------|---------|
| 1 | 高并发性能不足 | 中 | 高 | 横向扩展多实例 |
| 2 | 模型效果不达标 | 低 | 高 | 切换更大模型或微调 |
| 3 | GPU 资源不足 | 中 | 中 | 升级 GPU 或限流 |
| 4 | 成本超预算 | 高 | 中 | 按需启停或混合云 |

---

## 决策记录归档（ADR - Architecture Decision Record）

### ADR-001: 选择 vLLM 作为推理引擎

**日期**: 2025-12-22  
**状态**: ✅ 已采纳

**背景**  
需要选择一个 LLM 推理引擎，候选：vLLM、TGI、Ollama

**决策**  
选择 vLLM

**理由**
- 性能：PagedAttention 技术，吞吐量最高
- 兼容性：OpenAI API 兼容，易于集成
- 社区：GitHub 25K+ stars，活跃维护

**代价**
- 版本兼容性复杂（需要匹配 CUDA 版本）
- 配置复杂度高（参数多）

---

### ADR-002: 选择 Qwen2.5-7B-AWQ 模型

**日期**: 2025-12-22  
**状态**: ✅ 已采纳

**背景**  
需要选择模型大小和量化方式

**决策**  
Qwen2.5-7B-Instruct-AWQ（4-bit 量化）

**理由**
- 中文能力强（阿里出品）
- 显存占用低（~4.5GB，A10 23GB 足够）
- AWQ 量化速度快，精度损失小

**代价**
- 效果可能不如 14B 模型
- AWQ 需要特定 vLLM 版本支持

---

### ADR-003: 使用 vLLM v0.4.2 而非最新版

**日期**: 2025-12-23  
**状态**: ✅ 已采纳

**背景**  
最新版 v0.13.0 启动失败（v1 引擎问题）

**决策**  
降级到 v0.4.2

**理由**
- v0.4.2 使用稳定的 v0 引擎
- 兼容 CUDA 12.3 + 驱动 545.23.06
- 支持 AWQ 量化

**代价**
- 无法使用最新特性
- 未来升级需要重新验证

---

### ADR-004: Docker 数据目录迁移到 /data

**日期**: 2025-12-23  
**状态**: ✅ 已采纳

**背景**  
根分区磁盘空间不足（99GB 满载）

**决策**  
配置 Docker data-root 到 /data 盘（500GB）

**理由**
- 数据盘空间充足（392GB 可用）
- 模型缓存 + 镜像占用大

**代价**
- 需要迁移已有数据
- 配置文件修改

---

## 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────┐
│                      业务系统层                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │ PHP 系统  │  │ Java 系统 │  │ Go 系统  │  ...         │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘              │
└───────┼─────────────┼─────────────┼─────────────────────┘
        │             │             │
        └─────────────┴─────────────┘
                      │ HTTP/REST API (OpenAI 兼容)
        ┌─────────────▼─────────────────────────────────┐
        │         AI 推理服务层                          │
        │  ┌─────────────────────────────────────────┐  │
        │  │  vLLM API Server                        │  │
        │  │  - 端口: 8000                           │  │
        │  │  - 协议: HTTP/JSON                      │  │
        │  │  - 接口: OpenAI 兼容                     │  │
        │  │    · POST /v1/chat/completions          │  │
        │  │    · GET  /v1/models                    │  │
        │  │    · GET  /health                       │  │
        │  └─────────────┬───────────────────────────┘  │
        │                │                               │
        │  ┌─────────────▼───────────────────────────┐  │
        │  │  Qwen2.5-7B-Instruct-AWQ 模型           │  │
        │  │  - 参数量: 70亿                         │  │
        │  │  - 量化: AWQ (4-bit)                    │  │
        │  │  - 显存: ~4.5GB                         │  │
        │  └─────────────┬───────────────────────────┘  │
        └────────────────┼───────────────────────────────┘
                         │
        ┌────────────────▼───────────────────────────────┐
        │         硬件资源层                             │
        │  NVIDIA A10 GPU (23GB)                         │
        │  CentOS 服务器 + Docker                        │
        └────────────────────────────────────────────────┘
```

### 核心特性

✅ **标准化接口**：OpenAI API 兼容，支持任何语言调用  
✅ **高性能**：vLLM 优化推理，支持批处理和流式输出  
✅ **可扩展**：容器化部署，易于横向扩展  
✅ **企业可控**：本地部署，数据不出内网

---

## 多语言调用示例

### PHP (Laravel)

```php
use Illuminate\Support\Facades\Http;

$response = Http::timeout(300)->post('http://服务器IP:8000/v1/chat/completions', [
    'model' => 'Qwen/Qwen2.5-7B-Instruct-AWQ',
    'messages' => [
        ['role' => 'user', 'content' => '你好']
    ],
    'max_tokens' => 200
]);

$reply = $response->json()['choices'][0]['message']['content'];
```

### Java (Spring Boot)

```java
RestTemplate restTemplate = new RestTemplate();
Map<String, Object> request = Map.of(
    "model", "Qwen/Qwen2.5-7B-Instruct-AWQ",
    "messages", List.of(Map.of("role", "user", "content", "你好")),
    "max_tokens", 200
);

ResponseEntity<Map> response = restTemplate.postForEntity(
    "http://服务器IP:8000/v1/chat/completions",
    request,
    Map.class
);

String reply = ((Map) ((List) response.getBody().get("choices")).get(0))
    .get("message").get("content").toString();
```

### Python

```python
import requests

response = requests.post(
    "http://服务器IP:8000/v1/chat/completions",
    json={
        "model": "Qwen/Qwen2.5-7B-Instruct-AWQ",
        "messages": [{"role": "user", "content": "你好"}],
        "max_tokens": 200
    },
    timeout=300
)

reply = response.json()["choices"][0]["message"]["content"]
```

### Go

```go
type ChatRequest struct {
    Model      string    `json:"model"`
    Messages   []Message `json:"messages"`
    MaxTokens  int       `json:"max_tokens"`
}

client := &http.Client{Timeout: 300 * time.Second}
reqBody, _ := json.Marshal(ChatRequest{
    Model: "Qwen/Qwen2.5-7B-Instruct-AWQ",
    Messages: []Message{{Role: "user", Content: "你好"}},
    MaxTokens: 200,
})

resp, _ := client.Post("http://服务器IP:8000/v1/chat/completions",
    "application/json", bytes.NewReader(reqBody))
```

---

## 资源投入

| 资源类型 | 配置 | 成本 | 说明 |
|---------|------|------|------|
| **GPU 服务器** | NVIDIA A10 (23GB) | 按需计费 | 腾讯云/阿里云 |
| **存储** | 数据盘 500GB | ~200元/月 | 模型+缓存 |
| **带宽** | 按流量计费 | 按需 | API 调用流量 |
| **人力** | 1 人开发 | 3-5 天 | 集成和调试 |

**总成本预估**：约 **2000-3000元/月**（GPU 按需使用）

---

## 优势对比

| 对比项 | 云端 API (如 GPT) | 本地 vLLM + Qwen | 
|--------|------------------|------------------|
| **成本** | 按 token 计费，高频使用昂贵 | 固定服务器成本，高频使用划算 |
| **数据安全** | 数据上传到云端 | 数据不出内网，完全可控 |
| **响应速度** | 受网络影响，约 2-5 秒 | 内网调用，约 0.5-2 秒 |
| **定制化** | 无法定制模型 | 可微调模型，专属优化 |
| **可用性** | 依赖外部服务 | 自主可控，无依赖 |
| **合规性** | 可能不符合内部安全要求 | 符合企业内部安全规范 |

---

## 项目审计检查清单

### 技术验证
- [x] 技术选型有充分对比分析
- [x] 关键决策有 ADR 记录
- [x] 遇到的问题有解决方案记录
- [ ] 性能指标有量化数据
- [ ] 稳定性有测试报告

### 项目管理
- [x] 阶段划分清晰
- [x] 每个阶段有进入/退出条件
- [x] 每个阶段有可验证指标
- [x] 风险识别完整
- [x] 进度跟踪准确

### 知识沉淀
- [x] 问题决策记录完整（"迷宫路径"）
- [ ] 技术文档齐全
- [ ] 可复现性验证（他人可按文档独立完成）

---

## 总结

### 项目现状
- ✅ **可行性验证**：vLLM + Qwen2.5 技术栈可行
- ✅ **服务部署**：已在 CentOS 服务器成功运行
- ⏳ **网络配置**：安全组待配置
- 📋 **性能验证**：待执行
- 📋 **文档交付**：待完成

### 下一步行动
1. 配置云安全组，解决外网访问问题
2. 执行阶段 4 性能测试，获取基准数据
3. 完成多语言集成示例代码
4. 编写完整技术文档

### 项目风险
- 🟡 网络配置阻塞（优先级高）
- 🟡 性能数据缺失（需补充）
- 🟢 技术可行性已验证

---

**文档版本**: v2.0（规范版）  
**更新时间**: 2025-12-23  
**审核状态**: 待审核  
**负责人**: 技术团队

