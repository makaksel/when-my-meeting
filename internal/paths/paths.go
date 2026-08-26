package paths

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

type Paths struct {
	ConfigDir  string
	ConfigPath string
	DataDir    string
	CacheDir   string
}

func New(appName string) (*Paths, error) {

	return &Paths{
		ConfigDir:  filepath.Join(xdg.ConfigHome, appName),
		ConfigPath: filepath.Join(xdg.ConfigHome, appName, "config.yaml"),
		DataDir:    filepath.Join(xdg.DataHome, appName),
		CacheDir:   filepath.Join(xdg.CacheHome, appName),
	}, nil
}
