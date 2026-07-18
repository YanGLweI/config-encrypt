package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateKeyPair 生成 RSA 密钥对并写入文件
// keySize: 密钥位数（推荐 2048 或 4096）
// privateKeyPath: 私钥输出路径
// publicKeyPath: 公钥输出路径
func GenerateKeyPair(keySize int, privateKeyPath, publicKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("生成密钥对失败: %w", err)
	}

	// 编码私钥为 PEM 格式（PKCS#1，与后端现有格式兼容）
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := os.WriteFile(privateKeyPath, pem.EncodeToMemory(privateKeyPEM), 0600); err != nil {
		return fmt.Errorf("写入私钥文件失败: %w", err)
	}

	// 导出并编码公钥为 PEM 格式（PKCS#1）
	publicKeyDER := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	if err := os.WriteFile(publicKeyPath, pem.EncodeToMemory(publicKeyPEM), 0644); err != nil {
		return fmt.Errorf("写入公钥文件失败: %w", err)
	}

	return nil
}

// LoadPublicKey 从 PEM 文件加载 RSA 公钥
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取公钥文件失败: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解码 PEM 块，请确认文件格式正确")
	}

	// 尝试 PKCS#1 格式
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return pub, nil
	}

	// 尝试 PKIX 格式
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析公钥，请确认格式为 PKCS#1 或 PKIX")
	}

	rsaPub, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥类型不是 RSA")
	}
	return rsaPub, nil
}

// LoadPrivateKey 从 PEM 文件加载 RSA 私钥
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解码 PEM 块，请确认文件格式正确")
	}

	// 尝试 PKCS#1 格式
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privKey, nil
	}

	// 尝试 PKCS#8 格式
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析私钥，请确认格式为 PKCS#1 或 PKCS#8")
	}

	rsaPriv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥类型不是 RSA")
	}
	return rsaPriv, nil
}

// Encrypt 使用 RSA-OAEP + SHA-256 加密明文
// 返回值: Base64 编码的密文字符串
func Encrypt(publicKey *rsa.PublicKey, plaintext string) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("加密失败: %w", err)
	}
	return encode(ciphertext), nil
}

// Decrypt 使用 RSA-OAEP + SHA-256 解密密文
// 输入: Base64 编码的密文字符串
func Decrypt(privateKey *rsa.PrivateKey, ciphertext string) (string, error) {
	data, err := decode(ciphertext)
	if err != nil {
		return "", fmt.Errorf("解码失败: %w", err)
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}
	return string(plaintext), nil
}
