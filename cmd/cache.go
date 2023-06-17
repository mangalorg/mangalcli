package cmd

import (
	"fmt"
	"github.com/mangalorg/mangalcli/fs"
	"github.com/mangalorg/mangalcli/path"
)

type cacheCmd struct {
	Path  cachePathCmd  `cmd:"" help:"Show cache directory path"`
	Clear cacheClearCmd `cmd:"" help:"Remove cache directory"`
}

type cachePathCmd struct{}

func (c *cachePathCmd) Run(*Context) {
	fmt.Println(path.Cache())
}

type cacheClearCmd struct{}

func (c *cacheClearCmd) Run(*Context) {
	fs.FS.RemoveAll(path.Cache())
}
