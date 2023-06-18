package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"io"
	"os"
)

var cmd struct {
	Run      runCmd      `cmd:"" help:"Run given script as string"`
	Cache    cacheCmd    `cmd:"" help:"Cache manipulation"`
	Download downloadCmd `cmd:"" help:"Download lua providers from GitHub"`
	Log      string      `enum:"stdout,stderr,none" default:"stderr" help:"Logging output. Possible values: stdout, stderr, none"`
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

	log.SetDefault(log.New(logWriter))

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
