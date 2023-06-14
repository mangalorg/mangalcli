package path

import (
	"github.com/mangalorg/mangalcli/meta"
	"os"
	"path/filepath"
)

func Config() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join("."+meta.AppName, "config")
	}

	return filepath.Join(configDir, meta.AppName)
}

func Cache() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join("."+meta.AppName, "cache")
	}

	return filepath.Join(cacheDir, meta.AppName)
}
