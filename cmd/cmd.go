package cmd

import "github.com/alecthomas/kong"

var cmd struct {
	Anilist anilistCmd `cmd:""`
	Run     runCmd     `cmd:""`
}

func Run() {
	ctx := kong.Parse(&cmd)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
