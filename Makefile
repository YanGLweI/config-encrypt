BINARY_NAME=config-encrypt
VERSION=1.0.0

# 输出目录
DIST_DIR=dist

.PHONY: build build-all build-linux build-windows build-darwin build-darwin-arm64 clean \
       build-gui build-gui-windows build-gui-darwin build-gui-all

# 编译当前平台
build:
	go build -ldflags "-s -w" -o $(BINARY_NAME) .

# 编译所有平台
build-all: build-linux build-windows build-darwin build-darwin-arm64

# Linux amd64
build-linux:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

# Windows amd64
build-windows:
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# macOS Intel
build-darwin:
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

# macOS Apple Silicon
build-darwin-arm64:
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

# 清理
clean:
	rm -rf $(DIST_DIR) $(BINARY_NAME) $(BINARY_NAME)-gui

# ==================== GUI 编译 ====================
# 注意: Fyne 依赖 CGO/OpenGL，跨平台编译需使用 fyne-cross (基于 Docker)
# 安装: go install github.com/fyne-io/fyne-cross@latest

# 编译 GUI（当前平台，需本地有 C 编译器）
build-gui:
	go build -ldflags "-s -w" -o $(BINARY_NAME)-gui ./cmd/gui

# GUI Windows amd64（需要 fyne-cross）
build-gui-windows:
	@mkdir -p $(DIST_DIR)
	fyne-cross windows -app-id com.config-encrypt.gui -name $(BINARY_NAME)-gui ./cmd/gui
	@cp fyne-cross/bin/windows-amd64/$(BINARY_NAME)-gui.exe $(DIST_DIR)/ 2>/dev/null || true

# GUI macOS（需要 fyne-cross）
build-gui-darwin:
	@mkdir -p $(DIST_DIR)
	fyne-cross darwin -app-id com.config-encrypt.gui -name $(BINARY_NAME)-gui ./cmd/gui
	@cp fyne-cross/bin/darwin-amd64/$(BINARY_NAME)-gui.app $(DIST_DIR)/ 2>/dev/null || true
	@cp fyne-cross/bin/darwin-arm64/$(BINARY_NAME)-gui.app $(DIST_DIR)/ 2>/dev/null || true

# GUI Linux amd64（需要 fyne-cross）
build-gui-linux:
	@mkdir -p $(DIST_DIR)
	fyne-cross linux -app-id com.config-encrypt.gui -name $(BINARY_NAME)-gui ./cmd/gui
	@cp fyne-cross/bin/linux-amd64/$(BINARY_NAME)-gui $(DIST_DIR)/ 2>/dev/null || true

# GUI 编译所有平台
build-gui-all: build-gui-windows build-gui-darwin build-gui-linux
