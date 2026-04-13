SHELL := /bin/bash

VERSION ?= $(shell sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' version.json | head -1)
ifeq ($(strip $(VERSION)),)
VERSION := dev
endif
BUILD_DATE ?= $(shell date +"%Y%m%dT%H%M")
VERSION_CONFIG := version.json
VERSION_PKG := github.com/loula/pic2video/internal/app/version
LDFLAGS := -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).BuildDate=$(BUILD_DATE)

define WRITE_VERSION_FILE
	@printf '{\n  "version": "%s",\n  "build_date": "%s"\n}\n' "$(VERSION)" "$(BUILD_DATE)" > $(1)
endef

HOST_OS := $(shell uname -s)
IS_WSL := $(shell grep -qiE '(microsoft|wsl)' /proc/version 2>/dev/null && echo 1 || echo 0)
HAS_APT := $(shell command -v apt-get >/dev/null 2>&1 && echo 1 || echo 0)
HAS_BREW := $(shell command -v brew >/dev/null 2>&1 && echo 1 || echo 0)

LINUX_GUI_DEPS := build-essential pkg-config libgl1-mesa-dev xorg-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libasound2-dev
LINUX_WIN_CROSS_DEPS := mingw-w64 gcc-mingw-w64-x86-64

.PHONY: \
	build build-gui run-gui \
	build-linux build-macos build-windows \
	build-gui-linux build-gui-macos build-gui-windows \
	build-all-cli build-all-gui build-all-full build-all \
	install-build-deps bootstrap-build-deps build-smart build-smart-cli build-smart-gui \
	test test-unit test-e2e fmt

build:
	mkdir -p bin/cli
	go build -ldflags "$(LDFLAGS)" -o bin/cli/pic2video ./cmd/pic2video
	$(call WRITE_VERSION_FILE,bin/cli/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-gui:
	mkdir -p bin/gui
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o bin/gui/pic2video-gui ./cmd/pic2video-gui
	$(call WRITE_VERSION_FILE,bin/gui/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

run-gui:
	go run ./cmd/pic2video-gui

build-linux:
	mkdir -p bin/cli
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/cli/pic2video-linux-amd64 ./cmd/pic2video
	$(call WRITE_VERSION_FILE,bin/cli/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-macos:
	mkdir -p bin/cli
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/cli/pic2video-darwin-amd64 ./cmd/pic2video
	$(call WRITE_VERSION_FILE,bin/cli/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-windows:
	mkdir -p bin/cli
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/cli/pic2video-windows-amd64.exe ./cmd/pic2video
	$(call WRITE_VERSION_FILE,bin/cli/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-gui-linux:
	mkdir -p bin/gui
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/gui/pic2video-gui-linux-amd64 ./cmd/pic2video-gui
	$(call WRITE_VERSION_FILE,bin/gui/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-gui-macos:
	mkdir -p bin/gui
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/gui/pic2video-gui-darwin-amd64 ./cmd/pic2video-gui
	$(call WRITE_VERSION_FILE,bin/gui/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-gui-windows:
	mkdir -p bin/gui
	CGO_ENABLED=1 CC=$${CC_WINDOWS:-x86_64-w64-mingw32-gcc} GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/gui/pic2video-gui-windows-amd64.exe ./cmd/pic2video-gui
	$(call WRITE_VERSION_FILE,bin/gui/$(VERSION_CONFIG))
	$(call WRITE_VERSION_FILE,$(VERSION_CONFIG))

build-all-cli: build-linux build-macos build-windows

build-all-gui: build-gui-linux build-gui-macos build-gui-windows

build-all-full: build-all-cli build-all-gui

install-build-deps:
	@echo "[deps] host=$(HOST_OS) wsl=$(IS_WSL)"
	@if [[ "$(HOST_OS)" == "Linux" ]]; then \
		if [[ "$(HAS_APT)" == "1" ]]; then \
			echo "[deps] installing Linux GUI + Windows cross deps via apt-get"; \
			sudo apt-get update; \
			sudo apt-get install -y $(LINUX_GUI_DEPS) $(LINUX_WIN_CROSS_DEPS); \
		else \
			echo "[deps] apt-get not found; install manually: $(LINUX_GUI_DEPS) $(LINUX_WIN_CROSS_DEPS)"; \
			exit 1; \
		fi; \
	elif [[ "$(HOST_OS)" == "Darwin" ]]; then \
		if [[ "$(HAS_BREW)" == "1" ]]; then \
			echo "[deps] installing macOS deps via brew"; \
			brew install pkg-config glfw molten-vk mingw-w64 || true; \
		else \
			echo "[deps] brew not found; install Homebrew and required libs manually"; \
			exit 1; \
		fi; \
	else \
		echo "[deps] unsupported host for automated dependency installation: $(HOST_OS)"; \
		exit 1; \
	fi

bootstrap-build-deps: install-build-deps

build-smart-cli: build-all-cli

build-smart-gui:
	@echo "[smart-gui] host=$(HOST_OS) wsl=$(IS_WSL)"
	@if [[ "$(HOST_OS)" == "Linux" ]]; then \
		$(MAKE) build-gui-linux; \
		$(MAKE) build-gui-windows; \
		echo "[smart-gui] skipping macOS GUI build on Linux/WSL (requires macOS SDK/toolchain; build on macOS runner)"; \
	elif [[ "$(HOST_OS)" == "Darwin" ]]; then \
		$(MAKE) build-gui-macos; \
		echo "[smart-gui] skipping Linux/Windows GUI cross-builds by default on macOS"; \
	else \
		echo "[smart-gui] unsupported host for GUI smart-build: $(HOST_OS)"; \
		exit 1; \
	fi

build-smart: build-smart-cli build-smart-gui

build-all: build-smart

test: test-unit test-e2e

test-unit:
	go test ./... -run Test -count=1

test-e2e:
	go test ./tests/e2e -count=1

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './.git/*')
