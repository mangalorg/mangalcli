package cmd

import (
	"context"
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangalcli/cache"
	"github.com/mangalorg/mangalcli/fs"
	"github.com/mangalorg/mangalcli/lua"
	"net/http"
	"os"
)

type runCmd struct {
	Path string            `arg:"" help:"path to the lua script" type:"existingfile"`
	Vars map[string]string `help:"variables to pass to the selection query"`
	Exec string            `help:"lua script to execute. See wiki for more" required:""`
}

func (r *runCmd) Run(ctx *Context) error {
	contents, err := os.ReadFile(r.Path)
	if err != nil {
		return err
	}

	loader, err := luaprovider.NewLoader(contents, luaprovider.Options{
		HTTPClient: &http.Client{},
		HTTPStore:  cache.New("lua-http"),
	})
	if err != nil {
		return err
	}

	clientOptions := libmangal.DefaultClientOptions()
	anilist := newAnilist(ctx)
	clientOptions.Anilist = &anilist
	clientOptions.FS = fs.FS
	clientOptions.Log = func(msg string) {
		fmt.Fprintln(ctx.LogWriter, msg)
	}

	client, err := libmangal.NewClient(context.Background(), loader, clientOptions)
	if err != nil {
		return err
	}

	return lua.Exec(context.Background(), &client, r.Vars, r.Exec)
}
