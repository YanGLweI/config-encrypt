BINARY_NAME=sftp-config-encrypt
VERSION=1.0.0

# 输出目录
DIST_DIR=dist

.PHONY: build build-all build-linux build-windows build-darwin build-darwin-arm64 clean

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
	rm -rf $(DIST_DIR) $(BINARY_NAME)
