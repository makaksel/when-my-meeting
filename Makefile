APP_NAME := when-my-meeting
APP_ID := com.makaksel.when-my-meeting
VERSION ?= 0.1.0
BUILD ?= 1

DIST := dist
BIN := $(DIST)/$(APP_NAME)

CMD := ./cmd/when-my-meeting

ICON_PNG := internal/assets/icon.png

LINUX_ASSET := $(DIST)/$(APP_NAME)-linux-amd64-$(VERSION).deb
WINDOWS_ASSET := $(DIST)/$(APP_NAME)-windows-amd64-$(VERSION).exe
MACOS_ASSET := $(DIST)/$(APP_NAME)-macos-universal-$(VERSION).dmg

MACOS_DIR := cmd/when-my-meeting
MACOS_APP := $(MACOS_DIR)/When My Meeting.app
MACOS_DIST := $(MACOS_DIR)/dist
MACOS_DMG_SRC := $(MACOS_DIR)/dmg
MACOS_FYNE_DMG := $(MACOS_DIR)/When My Meeting.dmg
MACOS_ASSET := $(DIST)/$(APP_NAME)-macos-universal-$(VERSION).dmg


.PHONY: \
	run \
	icons \
	build-linux \
	build-windows \
	build-macos \
	build-deb \
	clean

# Development_____________________
run: $(ICON_PNG)
	go run $(CMD)/main.go


# Linux_____________________
build-linux: $(ICON_PNG)
	mkdir -p $(DIST)

	go build \
		-ldflags="-s -w" \
		-o $(BIN) \
		$(CMD)

build-deb: build-linux
	nfpm package \
		--packager deb \
		--config packaging/linux/nfpm.yaml \
		--target $(LINUX_ASSET)

	rm -f $(BIN)


# Windows_____________________
build-windows: $(ICON_PNG)
	mkdir -p $(DIST)

	cd "$(CMD)" && \
	fyne package \
		-os windows \
		-icon "../../$(ICON_PNG)" \
		-name "When My Meeting" \
		-app-id "$(APP_ID)" \
		-app-version "$(VERSION)" \
		-app-build "$(BUILD)"

	mv \
		"$(CMD)/When My Meeting.exe" \
		"$(WINDOWS_ASSET)"


# MacOS_____________________
build-macos:
	rm -rf \
		"$(MACOS_APP)" \
		"$(MACOS_DIST)" \
		"$(MACOS_DMG_SRC)" \
		"$(MACOS_FYNE_DMG)"

	mkdir -p \
		"$(MACOS_DIST)" \
		"$(MACOS_DMG_SRC)" \
		"$(DIST)"

	cd "$(MACOS_DIR)" && \
	fyne package \
		-os darwin \
		-icon "../../$(ICON_PNG)" \
		-name "When My Meeting" \
		-app-id "$(APP_ID)" \
		-app-version "$(VERSION)" \
		-app-build "$(BUILD)"

	/usr/libexec/PlistBuddy \
		-c "Delete :LSUIElement" \
		"$(MACOS_APP)/Contents/Info.plist" \
		|| true

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
		"$(MACOS_FYNE_DMG)"

	mv \
		"$(MACOS_FYNE_DMG)" \
		"$(MACOS_ASSET)"

	rm -rf \
		"$(MACOS_DMG_SRC)" \
		"$(MACOS_APP)"


# Clean_____________________
clean:
	rm -rf "$(DIST)"
	rm -rf "$(MACOS_DIST)"
	rm -rf "$(MACOS_DMG_SRC)"
	rm -rf "$(MACOS_APP)"
