package utils

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func IsToday(target *time.Time) bool {
	return target.Format(time.DateOnly) == time.Now().Format(time.DateOnly)
}

func IsFolowing(target *time.Time) bool {
	return IsToday(target) && target.After(time.Now().Local())
}

func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()

	case "darwin":
		return exec.Command("open", url).Start()

	case "linux":
		return exec.Command("xdg-open", url).Start()

	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func ParseURL(location string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(location))
	if err != nil {
		return nil
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil
	}

	if u.Host == "" {
		return nil
	}

	return u
}

func LimitString(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
