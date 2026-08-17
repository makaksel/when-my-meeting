.PHONY: run gnome-refresh

run:
	go run ./cmd/when-my-meeting/main.go

gnome-refresh:
	gnome-extensions disable ubuntu-appindicators@ubuntu.com
	sleep 1
	gnome-extensions enable ubuntu-appindicators@ubuntu.com

p-kill:
	pgrep -f "/cmd/when-my-meeting/main.go" | grep -v "$$$$" | xargs -r kill -9

app-kill:
	make p-kill
	sleep 1
	make gnome-refresh
