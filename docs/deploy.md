# 构建部署

## 环境要求

| 组件 | 版本要求 | 用途 |
|------|----------|------|
| Go | 1.22+ | 后端编译 |
| MySQL | 8.0+ | 数据存储 |
| Node.js | 18+ | 前端构建（仅开发/编译时） |
| npm | 9+ | 前端依赖管理 |

## 数据库准备

首次使用前，创建数据库和用户：

```sql
CREATE DATABASE IF NOT EXISTS zang_ye
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'zangye'@'localhost'
  IDENTIFIED BY 'zangye123';

GRANT ALL PRIVILEGES ON zang_ye.* TO 'zangye'@'localhost';
FLUSH PRIVILEGES;
```

> 注意：当前代码中默认连接使用 `root/root@127.0.0.1:3306/zang_ye`。
> 如需修改，请编辑 `internal/db/mysql.go` 中的 `DefaultConfig()`。

## 编译

### 完整编译（含前端）

```bash
# 1. 编译前端
cd frontend
npm install
npm run build
cd ..

# 2. 编译 Go 后端（自动嵌入前端 dist）
go build -o zangye .

# 3. 运行
./zangye
```

### 仅编译后端（无前端界面）

```bash
# 如果 frontend/dist 不存在，Go 编译会跳过前端嵌入
go build -o zangye .
./zangye
# 启动后仅 API 可用，前端访问返回 404
```

### 交叉编译

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o zangye-linux-amd64 .

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o zangye.exe .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o zangye-darwin-amd64 .
```

## 运行

### 基本运行

```bash
./zangye
# 监听 127.0.0.1:27138
```

### 自定义监听地址

```bash
# 通过环境变量修改监听地址
ZANGYE_ADDR="0.0.0.0:8080" ./zangye
```

### 系统服务（systemd）

创建 `/etc/systemd/system/zangye.service`：

```ini
[Unit]
Description=藏叶 个人文件管理器
After=network.target mysql.service

[Service]
Type=simple
User=zangye
WorkingDirectory=/opt/zangye
ExecStart=/opt/zangye/zangye
Environment=ZANGYE_ADDR=127.0.0.1:27138
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now zangye
sudo systemctl status zangye
```

## 开发模式

### 前后端分离开发

```bash
# 终端 1：启动后端
go run .

# 终端 2：启动前端开发服务器（自动代理 API）
cd frontend && npm run dev
# 访问 http://localhost:5173
```

前端开发服务器通过 Vite proxy 将 `/api` 请求转发到 `http://127.0.0.1:27138`，避免跨域问题。

## 配置说明

| 配置项 | 默认值 | 修改方式 |
|--------|--------|----------|
| 监听地址 | `127.0.0.1:27138` | 环境变量 `ZANGYE_ADDR` |
| 数据库主机 | `127.0.0.1` | 修改 `DefaultConfig()` |
| 数据库端口 | `3306` | 修改 `DefaultConfig()` |
| 数据库名 | `zang_ye` | 修改 `DefaultConfig()` |
| 数据库用户 | `root` | 修改 `DefaultConfig()` |
| 数据库密码 | `root` | 修改 `DefaultConfig()` |

## 部署检查清单

- [ ] MySQL 已安装并运行
- [ ] 数据库 `zang_ye` 已创建
- [ ] 数据库用户权限已配置
- [ ] 前端已编译（`npm run build`）
- [ ] Go 二进制已编译
- [ ] 监听地址配置正确（如需外网访问，改为 `0.0.0.0`）
- [ ] 防火墙已开放端口（如需要）
- [ ] 健康检查端点可访问（`GET /api/health` 返回 `{"status":"ok"}`）