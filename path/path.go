package path

import (
	"github.com/mangalorg/mangalcli/meta"
	"os"
	"path/filepath"
)

func Cache() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join("."+meta.AppName, "cache")
	}

	return filepath.Join(cacheDir, meta.AppName)
}
