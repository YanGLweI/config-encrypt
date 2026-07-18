//go:build !darwin && !windows

package gui

// nativeOpenFile 非 macOS/Windows 平台不支持原生选择器
func nativeOpenFile(title string, defaultDir string, extensions []string) (string, bool) {
	return "", false
}

// nativeOpenFolder 非 macOS/Windows 平台不支持原生选择器
func nativeOpenFolder(title string, defaultDir string) (string, bool) {
	return "", false
}
