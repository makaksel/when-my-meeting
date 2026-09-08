APP_NAME := when-my-meeting
APP_ID := com.makaksel.when-my-meeting
VERSION := 0.1.0
BUILD := 1

DIST := dist
BIN := $(DIST)/$(APP_NAME)

ICON_SRC := internal/assets/icon.svg
ICON_PNG := internal/assets/icon.png
WINDOWS_ICON := internal/assets/icon_windows.ico

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

icons: $(ICON_PNG) $(WINDOWS_ICON)

run: $(ICON_PNG)
	go run ./cmd/when-my-meeting/main.go

build-deb: $(ICON_PNG)
	mkdir -p $(DIST)

	go build \
		-ldflags="-s -w" \
		-o $(BIN) \
		./cmd/when-my-meeting

	nfpm package \
		--packager deb \
		--config packaging/linux/nfpm.yaml \
		--target $(DIST)/$(APP_NAME)_$(VERSION)_amd64.deb

	rm -rf $(BIN)


build-windows: $(WINDOWS_ICON)
	mkdir -p $(DIST)
	go build -v \
        -ldflags="-H=windowsgui -s -w" \
        -o dist/when-my-meeting.exe \
        ./cmd/when-my-meeting

	cp \
	  $(WINDOWS_ICON) \
	  dist/icon.ico


build-macos:
	cd cmd/when-my-meeting

	rm -rf "When My Meeting.app" dist

	fyne package \
	  -os darwin \
	    -icon $(ICON_PNG) \
	    -name "When My Meeting" \
	    -app-id com.makaksel.when-my-meeting \
	    -app-version $(VERSION) \
	    -app-build $(BUILD)

	mv "When My Meeting.app" $(DIST)/
	cd ~/$(DIST)
	echo "Adding LSUIElement..."

	/usr/libexec/PlistBuddy \
	  -c "Add :LSUIElement bool true" \
	    "When My Meeting.app/Contents/Info.plist"

	echo "Removing extended attributes..."

	xattr -cr "When My Meeting.app"

	echo "Ad-hoc signing..."

	codesign \
	    --force \
        --sign - \
        "When My Meeting.app"

	echo "Verifying signature..."

	codesign \
        --verify \
        --deep \
        --strict \
        --verbose=4 \
        "When My Meeting.app"

	echo "Checking Info.plist..."

	plutil -p \
	  "When My Meeting.app/Contents/Info.plist"

	rm -rf dmg dist
	mkdir -p dmg dist

	ditto \
	"When My Meeting.app" \
	"dmg/When My Meeting.app"

	ln -s /Applications dmg/Applications

	hdiutil create \
	-volname "When My Meeting" \
	-srcfolder dmg \
	-format UDZO \
	-imagekey zlib-level=9 \
	-ov \
	"dist/When My Meeting.dmg"


clean:
	rm -rf $(DIST)
