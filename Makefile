APP_NAME := when-my-meeting
APP_ID := com.makaksel.when-my-meeting
VERSION := 0.1.0
BUILD := 1

DIST := dist
BIN := $(DIST)/$(APP_NAME)

ICON_SRC := internal/assets/icon.svg
ICON_PNG := internal/assets/icon.png
LINUX_ICON := packaging/linux/icon.png
WINDOWS_ICON := packaging/windows/icon.ico

$(ICON_PNG): $(ICON_SRC)
	mkdir -p $(dir $@)
	magick \
		-background none \
		$< \
		-resize 256x256 \
		$@

$(WINDOWS_ICON): $(ICON_SRC)
	mkdir -p $(dir $@)
	magick \
		-background none \
		$< \
		-define icon:auto-resize=16,32,48,64,128,256 \
		$@

.PHONY: build deb clean

run: $(ICON_PNG)
	go run ./cmd/when-my-meeting/main.go

build-linux: $(ICON_PNG)
	mkdir -p $(DIST)

	go build \
		-ldflags="-s -w" \
		-o $(BIN) \
		./cmd/when-my-meeting

deb: build
	nfpm package \
		--packager deb \
		--config packaging/linux/nfpm.yaml \
		--target $(DIST)/$(APP_NAME)_$(VERSION)_amd64.deb

	rm -rf $(BIN)



build-windows: $(WINDOWS_ICON)
	mkdir -p $(DIST)

windows-icons: $(WINDOWS_ICON)

build-macos:
	mkdir -p $(DIST)


clean:
	rm -rf $(DIST)
