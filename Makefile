APP_NAME := when-my-meeting
APP_ID := com.makaksel.when-my-meeting
VERSION := 0.1.0
BUILD := 1

DIST := dist
BIN := $(DIST)/$(APP_NAME)

CMD := ./cmd/when-my-meeting

ICON_SRC := internal/assets/icon.svg
ICON_PNG := internal/assets/icon.png
WINDOWS_ICON := internal/assets/icon_windows.ico

MACOS_DIR := cmd/when-my-meeting
MACOS_APP := $(MACOS_DIR)/When My Meeting.app
MACOS_DIST := $(MACOS_DIR)/dist
MACOS_DMG := $(MACOS_DIST)/When My Meeting.dmg
MACOS_DMG_SRC := $(MACOS_DIR)/dmg

.PHONY: \
	run \
	icons \
	build-linux \
	build-windows \
	build-macos \
	build-deb \
	clean


# Assets_____________________
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

icons: $(ICON_PNG) $(WINDOWS_ICON)

# Development_____________________
run: $(ICON_PNG)
	go run ./cmd/when-my-meeting/main.go


# Linux_____________________
build-linux: $(ICON_PNG)
	mkdir -p $(DIST)

	go build \
		-ldflags="-s -w" \
		-o $(BIN) \
		./cmd/when-my-meeting

build-deb: build-linux
	nfpm package \
		--packager deb \
		--config packaging/linux/nfpm.yaml \
		--target $(DIST)/$(APP_NAME)_$(VERSION)_amd64.deb

	rm -f $(BIN)


# Windows_____________________
build-windows: $(WINDOWS_ICON)
	mkdir -p $(DIST)
	go build -v \
        -ldflags="-H=windowsgui -s -w" \
        -o dist/when-my-meeting.exe \
        ./cmd/when-my-meeting

	cp \
        $(WINDOWS_ICON) \
        dist/icon.ico


# MacOS_____________________
build-macos:
	rm -rf "$(MACOS_APP)" "$(MACOS_DIST)" "$(MACOS_DMG_SRC)"

	mkdir -p "$(MACOS_DIST)" "$(MACOS_DMG_SRC)"

	cd "$(MACOS_DIR)" && \
	fyne package \
		-os darwin \
		-icon "../../$(ICON_PNG)" \
		-name "When My Meeting" \
		-app-id "$(APP_ID)" \
		-app-version "$(VERSION)" \
		-app-build "$(BUILD)"

	/usr/libexec/PlistBuddy \
		-c "Add :LSUIElement bool true" \
		"$(MACOS_APP)/Contents/Info.plist"

	xattr -cr "$(MACOS_APP)"

	codesign \
		--force \
		--sign - \
		"$(MACOS_APP)"

	codesign \
		--verify \
		--deep \
		--strict \
		--verbose=4 \
		"$(MACOS_APP)"

	plutil -p \
		"$(MACOS_APP)/Contents/Info.plist"

	ditto \
		"$(MACOS_APP)" \
		"$(MACOS_DMG_SRC)/When My Meeting.app"

	ln -s /Applications \
		"$(MACOS_DMG_SRC)/Applications"

	hdiutil create \
		-volname "When My Meeting" \
		-srcfolder "$(MACOS_DMG_SRC)" \
		-format UDZO \
		-imagekey zlib-level=9 \
		-ov \
		"$(MACOS_DMG)"

	rm -rf "$(MACOS_DMG_SRC)"


# Clean_____________________
clean:
	rm -rf "$(DIST)"
	rm -rf "$(MACOS_DIST)"
	rm -rf "$(MACOS_DMG_SRC)"
	rm -rf "$(MACOS_APP)"
