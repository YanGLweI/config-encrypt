//go:build !darwin

package gui

// nativeOpenFile 非 macOS 平台暂不支持原生选择器，返回空
func nativeOpenFile(title string, defaultDir string, extensions []string) string {
	return ""
}

// nativeOpenFolder 非 macOS 平台暂不支持原生选择器，返回空
func nativeOpenFolder(title string, defaultDir string) string {
	return ""
}
