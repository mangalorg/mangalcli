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

func (c *cachePathCmd) Run() error {
	fmt.Println(path.Cache())
	return nil
}

type cacheClearCmd struct{}

func (c *cacheClearCmd) Run() error {
	return fs.FS.RemoveAll(path.Cache())
}
