package crypto

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncPrefix 加密字符串前缀标识
const EncPrefix = "ENC["
const EncSuffix = "]"

// Wrap 将密文包装为 ENC[base64] 格式
func Wrap(ciphertext string) string {
	return EncPrefix + ciphertext + EncSuffix
}

// Unwrap 从 ENC[base64] 格式中提取密文
func Unwrap(wrapped string) (string, bool) {
	if !IsEncrypted(wrapped) {
		return wrapped, false
	}
	inner := wrapped[len(EncPrefix) : len(wrapped)-len(EncSuffix)]
	return inner, true
}

// IsEncrypted 判断字符串是否为加密格式
func IsEncrypted(s string) bool {
	return strings.HasPrefix(s, EncPrefix) && strings.HasSuffix(s, EncSuffix)
}

// encode Base64 编码
func encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// decode Base64 解码
func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64 解码失败: %w", err)
	}
	return data, nil
}
