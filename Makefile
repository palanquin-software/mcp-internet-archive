APP_INSTALL_DIR = "$(HOME)/bin"
APP_NAME = "mcp-internet-archive"
APP_INSTALL_FILE = "$(APP_INSTALL_DIR)/$(APP_NAME)"

.PHONY: install
install:
	@ mkdir -p "$(APP_INSTALL_DIR)"
	@ go build -o "$(APP_INSTALL_FILE)" cmd/mcp/main.go
	@ chmod +x "$(APP_INSTALL_FILE)"
	@ echo "installed to $(APP_INSTALL_FILE)"