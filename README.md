# 配置加密工具

配置文件加密工具，用于加密 `config.yml` 中的敏感信息（如数据库密码、LDAP 密码、邮件密码等），避免明文存储。支持命令行（CLI）和图形界面（GUI）两种使用方式。

## 特性

- 双模式支持：命令行（CLI）+ 图形界面（GUI）
- 跨平台支持：Windows / Linux / macOS
- RSA-OAEP + SHA-256 加密算法，安全可靠
- 加密格式 `ENC[base64]`，可直接嵌入 YAML 配置文件
- 交互式密码输入（不回显），防止肩窥
- GUI 支持 macOS/Windows 原生文件选择器
- 后端程序启动时自动检测并解密

## 编译

### CLI 命令行版本

```bash
# 编译当前平台
go build -o config-encrypt .

# 交叉编译（所有平台）
make build-all
```

### GUI 图形界面版本

GUI 版本基于 [Fyne](https://fyne.io/) 框架，依赖 CGO 和 OpenGL，编译前需确保本地有 C 编译器。

#### 前置条件

| 平台 | 所需工具 |
|------|----------|
| macOS | Xcode Command Line Tools（`xcode-select --install`） |
| Windows | [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 或 [MSYS2](https://www.msys2.org/) |
| Linux | `gcc` + OpenGL 开发库（`libgl1-mesa-dev xorg-dev`） |

#### 编译当前平台

```bash
# macOS / Linux
go build -ldflags "-s -w" -o config-encrypt-gui ./cmd/gui

# Windows（需设置 CGO）
set CGO_ENABLED=1
go build -ldflags "-s -w -H windowsgui" -o config-encrypt-gui.exe ./cmd/gui
```

#### 交叉编译 Windows 版本（在 macOS 上）

需要安装 [mingw-w64](https://www.mingw-w64.org/)：

```bash
# macOS 安装 mingw-w64
brew install mingw-w64

# 编译 Windows GUI（含图标）
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
  go build -ldflags "-s -w -H windowsgui" -o config-encrypt-gui-windows-amd64.exe ./cmd/gui
```

> **注意**：Windows 图标通过 `.syso` 资源文件嵌入（`cmd/gui/icon_windows.amd64.syso`），Go 构建时自动链接，无需额外操作。

#### 交叉编译 macOS 版本（在 macOS 上）

```bash
# Intel (amd64)
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
  go build -ldflags "-s -w" -o config-encrypt-gui-darwin-amd64 ./cmd/gui

# Apple Silicon (arm64)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
  go build -ldflags "-s -w" -o config-encrypt-gui-darwin-arm64 ./cmd/gui
```

#### macOS 打包为 .app 应用

macOS 用户无法直接运行二进制文件，需打包为 `.app` 应用包并签名：

```bash
# 创建 .app 目录结构
mkdir -p config-encrypt-gui.app/Contents/{MacOS,Resources}
cp config-encrypt-gui config-encrypt-gui.app/Contents/MacOS/config-encrypt-gui
cp gui/assets/icon.png config-encrypt-gui.app/Contents/Resources/icon.png

# 创建 Info.plist
cat > config-encrypt-gui.app/Contents/Info.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>配置加密工具</string>
    <key>CFBundleDisplayName</key>
    <string>配置加密工具</string>
    <key>CFBundleIdentifier</key>
    <string>com.config-encrypt.gui</string>
    <key>CFBundleVersion</key>
    <string>2.0.3</string>
    <key>CFBundleShortVersionString</key>
    <string>2.0.3</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleExecutable</key>
    <string>config-encrypt-gui</string>
    <key>CFBundleIconFile</key>
    <string>icon.png</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
</dict>
</plist>
EOF

# Ad-hoc 代码签名（解决 Gatekeeper "已损坏" 提示）
codesign --force --deep --sign - config-encrypt-gui.app

# 打包为 zip
ditto -c -k --sequesterRsrc --keepParent config-encrypt-gui.app config-encrypt-gui.zip
```

#### 使用 Makefile 一键编译

```bash
# CLI 所有平台
make build-all

# GUI 当前平台
make build-gui

# GUI 所有平台（需 fyne-cross + Docker）
make build-gui-all
```

### 使用流程

```bash
# 1. 生成密钥对（只需执行一次）
./config-encrypt keygen

# 2. 将私钥复制到后端 key/ 目录
cp config-key.pem /path/to/backend/key/config-private.pem

# 3. 加密 config.yml 中的敏感密码
./config-encrypt encrypt -pub config-key.pub.pem "your_db_password"
# 输出: ENC[xxxxx]

# 4. 将加密字符串填入 config.yml
#    password: "ENC[xxxxx]"

# 5. 解密验证
./config-encrypt decrypt -key config-key.pem "ENC[xxxxx]"
```

## 命令说明

### keygen - 生成密钥对

```bash
# 默认生成 2048 位密钥
./config-encrypt keygen

# 指定密钥位数和输出文件名
./config-encrypt keygen -bits 4096 -o mykey
```

输出文件：
- `config-key.pem` - RSA 私钥（保密，不提交 Git）
- `config-key.pub.pem` - RSA 公钥（用于加密）

### encrypt - 加密密码

```bash
# 交互式输入（密码不回显）
./config-encrypt encrypt -pub config-key.pub.pem

# 命令行直接传入
./config-encrypt encrypt -pub config-key.pub.pem "my_password"
```

### decrypt - 解密密码

```bash
./config-encrypt decrypt -key config-key.pem "ENC[xxxxx]"
```

## config.yml 配置示例

```yaml
# 加密前（明文，不安全）
database:
  password: "my_db_password"

# 加密后（安全）
database:
  password: "ENC[KALfN/8NNaHFxnW6s65VuzOyCNbZNFqs...]"
```

后端程序启动时会自动检测 `ENC[...]` 格式并使用私钥解密。

## 需要加密的配置字段

| 配置段 | 字段 | 说明 |
|--------|------|------|
| `database` | `password` | 数据库密码 |
| `email` | `password` | 邮箱 SMTP 密码 |
| `ldap` | `password` | LDAP 绑定密码 |
| `jwt` | `secret` | JWT 签名密钥 |

## GUI 使用说明

GUI 版本提供三个功能页面，通过左侧边栏切换：

| 页面 | 功能 |
|------|------|
| 🔐 密钥生成 | 生成 RSA 密钥对，选择保存目录 |
| 🔒 加密 | 选择公钥文件，输入密码，生成 `ENC[...]` 加密结果 |
| 🔓 解密 | 选择私钥文件，输入密文，解密查看原始密码 |

**文件选择器**：macOS 和 Windows 使用系统原生文件浏览器（支持排序、侧边栏），Linux 回退到 Fyne 自带对话框。

## 交叉编译

```bash
# CLI 所有平台
make build-all

# 或单独编译
make build-linux    # Linux amd64
make build-windows  # Windows amd64 (.exe)
make build-darwin   # macOS amd64
make build-darwin-arm64  # macOS Apple Silicon
```

## 安全建议

1. **私钥文件**（`*.pem`）必须加入 `.gitignore`，绝不能提交到 Git 仓库
2. 私钥文件权限应设为 `600`（仅所有者可读写）
3. 生产环境中，私钥应通过安全渠道分发到服务器
4. 定期轮换密钥对

## License

MIT
