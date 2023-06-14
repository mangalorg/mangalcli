package path

import (
	"github.com/mangalorg/mangalcli/meta"
	"log"
	"os"
	"path/filepath"
)

func Config() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(configDir, meta.AppName)
}
