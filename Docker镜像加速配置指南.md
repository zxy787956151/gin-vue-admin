# Docker Desktop 镜像加速配置指南（Mac 版）

## 📋 配置步骤

### 1. 打开 Docker Desktop

在 Mac 上找到并打开 Docker Desktop 应用。

### 2. 进入设置界面

点击右上角的 **⚙️ 图标**（齿轮），或者从菜单栏选择：
```
Docker Desktop → Preferences/Settings
```

### 3. 选择 Docker Engine

在左侧菜单中选择 **"Docker Engine"**

### 4. 编辑 JSON 配置

在右侧的 JSON 编辑器中，你会看到类似这样的内容：

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false
}
```

### 5. 添加镜像加速器配置

**将整个内容替换为**（直接复制下面的）：

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://dockerproxy.com",
    "https://docker.nju.edu.cn",
    "https://docker.mirrors.sjtug.sjtu.edu.cn"
  ],
  "dns": ["8.8.8.8", "8.8.4.4"],
  "experimental": false,
  "features": {
    "buildkit": true
  }
}
```

### 6. 应用并重启

1. 点击右下角的 **"Apply & Restart"** 按钮
2. 等待 Docker Desktop 重启（约 30-60 秒）
3. 确认 Docker 图标变为绿色（运行中）

### 7. 验证配置

在终端运行：

```bash
docker info | grep -A 10 "Registry Mirrors"
```

应该能看到你配置的镜像源。

### 8. 测试拉取

```bash
docker pull ollama/ollama:latest
```

如果配置成功，拉取速度会显著提升！

---

## 🔍 故障排查

### 问题1：配置后仍然很慢

尝试更换镜像源顺序，或添加其他镜像源：

```json
"registry-mirrors": [
    "https://dockerproxy.com",          // 试试把这个放第一
    "https://hub.rat.dev",
    "https://docker.m.daocloud.io"
]
```

### 问题2：Docker 重启失败

检查 JSON 格式是否正确：
- 是否有多余的逗号
- 括号是否配对
- 引号是否正确

可以使用在线 JSON 验证工具：https://jsonlint.com/

### 问题3：仍然超时

可能镜像源也被限制了，尝试方案2（使用代理）

---

## 🌐 替代方案：使用代理

如果你有 HTTP 代理（如 Clash、Shadowsocks 等）：

### 方式1：在 Docker Desktop 中配置

1. 打开 Docker Desktop → Settings/Preferences
2. 选择 **Resources** → **Proxies**
3. 启用 **"Manual proxy configuration"**
4. 填写代理信息：
   - HTTP Proxy: `http://127.0.0.1:7890` （根据你的代理端口）
   - HTTPS Proxy: `http://127.0.0.1:7890`
5. Apply & Restart

### 方式2：临时设置环境变量

```bash
export HTTP_PROXY=http://127.0.0.1:7890
export HTTPS_PROXY=http://127.0.0.1:7890
docker pull ollama/ollama:latest
```

---

## 📦 方案3：离线导入镜像

如果你有其他能正常访问 Docker Hub 的环境：

### 在能访问的机器上：

```bash
# 拉取镜像
docker pull ollama/ollama:latest

# 导出镜像
docker save ollama/ollama:latest | gzip > ollama-image.tar.gz

# 传输到你的 Mac（通过网盘、U盘等）
```

### 在你的 Mac 上：

```bash
# 导入镜像
docker load < ollama-image.tar.gz

# 验证
docker images | grep ollama
```

---

## ✅ 验证安装

配置成功后，运行：

```bash
# 查看镜像
docker images | grep ollama

# 应该看到类似输出：
# ollama/ollama    latest    xxxxx    2 days ago    1.5GB

# 启动容器测试
docker run -d --name test-ollama -p 11434:11434 ollama/ollama:latest

# 测试服务
curl http://localhost:11434/api/tags

# 清理测试容器
docker stop test-ollama && docker rm test-ollama
```

---

## 🎯 推荐镜像源（2024-12-24 更新）

可用性从高到低：

1. **dockerproxy.com** - 稳定性好
2. **docker.m.daocloud.io** - DaoCloud 官方
3. **docker.nju.edu.cn** - 南京大学
4. **docker.mirrors.sjtug.sjtu.edu.cn** - 上海交大

注意：镜像源的可用性会变化，建议配置多个作为备份。

---

## 💡 额外提示

1. **定期更新镜像源列表**：网上搜索"Docker 镜像加速器"获取最新可用源
2. **考虑使用代理**：如果经常需要拉取镜像，建议配置全局代理
3. **本地安装 Ollama**：对于 Mac 用户，本地安装比 Docker 更简单：
   ```bash
   brew install ollama
   ```

---

## 📞 需要帮助？

如果以上方案都不行，请：

1. 检查网络连接
2. 尝试连接手机热点（有时移动网络更好）
3. 考虑使用 VPN
4. 或者选择本地安装 Ollama（无需 Docker）

---

**配置完成后，记得重启 Docker Desktop！**

