# SFTP Config Encrypt

SFTP 管理系统配置文件加密工具，用于加密 `config.yml` 中的敏感信息（如数据库密码、LDAP 密码、邮件密码等），避免明文存储。

## 特性

- 跨平台支持：Windows / Linux / macOS
- RSA-OAEP + SHA-256 加密算法，安全可靠
- 加密格式 `ENC[base64]`，可直接嵌入 YAML 配置文件
- 交互式密码输入（不回显），防止肩窥
- 后端程序启动时自动检测并解密

## 快速开始

### 编译

```bash
# 编译当前平台
go build -o sftp-config-encrypt .

# 交叉编译（所有平台）
make build-all
```

### 使用流程

```bash
# 1. 生成密钥对（只需执行一次）
./sftp-config-encrypt keygen

# 2. 将私钥复制到后端 key/ 目录
cp config-key.pem /path/to/backend/key/config-private.pem

# 3. 加密 config.yml 中的敏感密码
./sftp-config-encrypt encrypt -pub config-key.pub.pem "your_db_password"
# 输出: ENC[xxxxx]

# 4. 将加密字符串填入 config.yml
#    password: "ENC[xxxxx]"

# 5. 解密验证
./sftp-config-encrypt decrypt -key config-key.pem "ENC[xxxxx]"
```

## 命令说明

### keygen - 生成密钥对

```bash
# 默认生成 2048 位密钥
./sftp-config-encrypt keygen

# 指定密钥位数和输出文件名
./sftp-config-encrypt keygen -bits 4096 -o mykey
```

输出文件：
- `config-key.pem` - RSA 私钥（保密，不提交 Git）
- `config-key.pub.pem` - RSA 公钥（用于加密）

### encrypt - 加密密码

```bash
# 交互式输入（密码不回显）
./sftp-config-encrypt encrypt -pub config-key.pub.pem

# 命令行直接传入
./sftp-config-encrypt encrypt -pub config-key.pub.pem "my_password"
```

### decrypt - 解密密码

```bash
./sftp-config-encrypt decrypt -key config-key.pem "ENC[xxxxx]"
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

## 交叉编译

```bash
# 使用 Makefile 编译所有平台
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
