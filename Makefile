BINARY      = wclaude-setup
BUILD_DIR   = bin
INSTALL_DIR = $(HOME)/.local/bin
VERSION     = 0.1.0
LDFLAGS     = -s -w -X main.version=$(VERSION)

.PHONY: build install clean test fmt lint

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/wclaude-setup/

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "$(BINARY) instalado em $(INSTALL_DIR)/$(BINARY)"
	@echo "Verifique que $(INSTALL_DIR) esta no PATH"

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./... -race -count=1

fmt:
	go fmt ./...

lint:
	go vet ./...
