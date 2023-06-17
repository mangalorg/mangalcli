package cmd

import (
	"github.com/alecthomas/kong"
	"io"
	"os"
)

type Context struct {
	LogWriter io.Writer
}

var cmd struct {
	Run   runCmd   `cmd:"" help:"Run given script as string"`
	Cache cacheCmd `cmd:"" help:"Cache manipulation"`
	Log   string   `enum:"stdout,stderr,none" default:"none" help:"Logging output. Possible values: stdout, stderr, none"`
}

func Run() {
	ctx := kong.Parse(&cmd)

	var logWriter io.Writer

	switch cmd.Log {
	case "none":
		logWriter = io.Discard
	case "stdout":
		logWriter = os.Stdout
	case "stderr":
		logWriter = os.Stderr
	default:
		panic("unknown log")
	}

	err := ctx.Run(&Context{LogWriter: logWriter})
	ctx.FatalIfErrorf(err)
}
