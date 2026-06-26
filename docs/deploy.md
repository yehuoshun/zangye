# 部署指南

## 单 exe 部署

### 构建

```bash
# 1. 构建前端
cd frontend
npm install
npm run build

# 2. 构建 Go 二进制
cd ..
go build -o zangye.exe -ldflags="-s -w" .
```

### 运行

```bash
# Windows
zangye.exe

# 指定端口
set ZANGYE_ADDR=0.0.0.0:27138
zangye.exe

# 指定数据库
set ZANGYE_DSN=user:password@tcp(host:port)/zangye?charset=utf8mb4&parseTime=True&loc=Local
zangye.exe
```

### 交叉编译

```bash
# Windows amd64
GOOS=windows GOARCH=amd64 go build -o zangye.exe -ldflags="-s -w" .

# Linux amd64
GOOS=linux GOARCH=amd64 go build -o zangye -ldflags="-s -w" .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o zangye-mac -ldflags="-s -w" .
```

## CI/CD

GitHub Actions 配置在 `.github/workflows/build.yml`：

- 触发：push 到 main 分支
- 步骤：构建前端 → 交叉编译 Windows exe → 上传 artifact

## 数据库初始化

```bash
mysql -u root -p < internal/db/schema.sql
```

## 环境要求

- MySQL 8.0+
- 操作系统：Windows/Linux/macOS
- 内存：最低 64MB
- 磁盘：二进制文件约 15MB
