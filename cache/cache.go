package cache

import (
	"github.com/mangalorg/mangalcli/fs"
	"github.com/mangalorg/mangalcli/path"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bbolt"
	"github.com/philippgille/gokv/encoding"
	"log"
	"path/filepath"
)

func New(name string) gokv.Store {
	cacheDir := path.Cache()
	err := fs.FS.MkdirAll(cacheDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	store, err := bbolt.NewStore(bbolt.Options{
		BucketName: name,
		Path:       filepath.Join(cacheDir, name+".db"),
		Codec:      encoding.Gob,
	})

	if err != nil {
		log.Fatal(err)
	}

	return store
}
