package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/YanGLweI/sftp-config-encrypt/crypto"
	"golang.org/x/term"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "keygen":
		cmdKeygen()
	case "encrypt":
		cmdEncrypt()
	case "decrypt":
		cmdDecrypt()
	case "help", "-h", "--help":
		printUsage()
	case "version", "-v", "--version":
		fmt.Printf("sftp-config-encrypt v%s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`sftp-config-encrypt - SFTP 管理系统配置加密工具

用法:
  sftp-config-encrypt <命令> [参数]

命令:
  keygen [选项]                   生成 RSA 密钥对
  encrypt -pub <公钥文件> [密码]  加密密码（输出 ENC[...] 格式）
  decrypt -key <私钥文件> <密文>  解密 ENC[...] 格式的密文
  help                            显示帮助信息
  version                         显示版本号

示例:
  # 1. 生成 2048 位密钥对
  sftp-config-encrypt keygen
  sftp-config-encrypt keygen -bits 4096 -o mykey

  # 2. 加密密码（交互式输入，密码不回显）
  sftp-config-encrypt encrypt -pub PublicKey.pem

  # 3. 加密密码（命令行直接传入）
  sftp-config-encrypt encrypt -pub PublicKey.pem "my_password"

  # 4. 解密验证
  sftp-config-encrypt decrypt -key PrivateKey.pem "ENC[xxxxx]"

输出格式:
  加密后的字符串格式为 ENC[base64密文]，可直接放入 config.yml 中。
  后端程序启动时会自动检测并解密 ENC[...] 格式的字段。`)
}

// cmdKeygen 生成 RSA 密钥对
func cmdKeygen() {
	bits := 2048
	output := "config-key"

	// 解析参数
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-bits":
			if i+1 >= len(args) {
				exitError("请指定密钥位数，如 -bits 2048")
			}
			fmt.Sscanf(args[i+1], "%d", &bits)
			i++
		case "-o":
			if i+1 >= len(args) {
				exitError("请指定输出文件名前缀，如 -o mykey")
			}
			output = args[i+1]
			i++
		default:
			exitError("未知参数: " + args[i])
		}
	}

	if bits < 1024 {
		exitError("密钥位数不能小于 1024")
	}

	privateKeyPath := output + ".pem"
	publicKeyPath := output + ".pub.pem"

	// 检查文件是否已存在
	for _, p := range []string{privateKeyPath, publicKeyPath} {
		if _, err := os.Stat(p); err == nil {
			exitError("文件已存在: %s，请先删除或使用 -o 指定其他名称", p)
		}
	}

	fmt.Printf("正在生成 %d 位 RSA 密钥对...\n", bits)
	if err := crypto.GenerateKeyPair(bits, privateKeyPath, publicKeyPath); err != nil {
		exitError("%v", err)
	}

	absPrivate, _ := filepath.Abs(privateKeyPath)
	absPublic, _ := filepath.Abs(publicKeyPath)

	fmt.Println("密钥对生成成功！")
	fmt.Printf("  私钥: %s （请妥善保管，不要提交到 Git）\n", absPrivate)
	fmt.Printf("  公钥: %s （用于加密配置文件中的密码）\n", absPublic)
	fmt.Println()
	fmt.Println("下一步: 将私钥复制到后端 key/ 目录，使用公钥加密 config.yml 中的敏感字段。")
}

// cmdEncrypt 加密密码
func cmdEncrypt() {
	pubKeyPath := ""
	password := ""

	// 解析参数
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-pub":
			if i+1 >= len(args) {
				exitError("请指定公钥文件路径，如 -pub PublicKey.pem")
			}
			pubKeyPath = args[i+1]
			i++
		default:
			// 非 flag 参数视为密码
			if password == "" {
				password = args[i]
			}
		}
	}

	if pubKeyPath == "" {
		exitError("请指定公钥文件路径，如 -pub PublicKey.pem")
	}

	// 如果没有传入密码，则交互式输入（不回显）
	if password == "" {
		password = readPasswordInteractive("请输入要加密的密码: ")
		if password == "" {
			exitError("密码不能为空")
		}
	}

	// 加载公钥
	pubKey, err := crypto.LoadPublicKey(pubKeyPath)
	if err != nil {
		exitError("加载公钥失败: %v", err)
	}

	// 加密
	ciphertext, err := crypto.Encrypt(pubKey, password)
	if err != nil {
		exitError("加密失败: %v", err)
	}

	// 输出结果
	result := crypto.Wrap(ciphertext)
	fmt.Println()
	fmt.Println("加密成功！请将以下内容复制到 config.yml 中：")
	fmt.Println()
	fmt.Printf("  %s\n", result)
	fmt.Println()
}

// cmdDecrypt 解密密码
func cmdDecrypt() {
	privateKeyPath := ""
	ciphertext := ""

	// 解析参数
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-key":
			if i+1 >= len(args) {
				exitError("请指定私钥文件路径，如 -key PrivateKey.pem")
			}
			privateKeyPath = args[i+1]
			i++
		default:
			if ciphertext == "" {
				ciphertext = args[i]
			}
		}
	}

	if privateKeyPath == "" {
		exitError("请指定私钥文件路径，如 -key PrivateKey.pem")
	}
	if ciphertext == "" {
		exitError("请提供要解密的密文字符串（ENC[...] 格式）")
	}

	// 提取密文
	inner, ok := crypto.Unwrap(ciphertext)
	if !ok {
		exitError("密文格式不正确，应为 ENC[base64] 格式")
	}

	// 加载私钥
	privKey, err := crypto.LoadPrivateKey(privateKeyPath)
	if err != nil {
		exitError("加载私钥失败: %v", err)
	}

	// 解密
	plaintext, err := crypto.Decrypt(privKey, inner)
	if err != nil {
		exitError("解密失败: %v", err)
	}

	fmt.Println()
	fmt.Printf("解密成功: %s\n", plaintext)
	fmt.Println()
}

// readPasswordInteractive 交互式读取密码（不回显）
func readPasswordInteractive(prompt string) string {
	fmt.Print(prompt)

	// 尝试使用 term 库禁用回显
	if term.IsTerminal(int(syscall.Stdin)) {
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println()
			exitError("读取密码失败: %v", err)
		}
		fmt.Println() // 换行
		return strings.TrimSpace(string(password))
	}

	// 非终端环境，直接读取一行
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// exitError 打印错误并退出
func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "错误: "+format+"\n", args...)
	os.Exit(1)
}
