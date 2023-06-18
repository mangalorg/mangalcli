package cmd

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangalcli/cache"
	"github.com/mangalorg/mangalcli/fs"
	"github.com/mangalorg/mangalcli/lua"
	"net/http"
	"os"
)

func newAnilist() libmangal.Anilist {
	options := libmangal.DefaultAnilistOptions()
	options.Log = func(msg string) {
		log.Info(msg)
	}

	options.QueryToIDsStore = cache.New("query-to-id")
	options.TitleToIDStore = cache.New("title-to-id")
	options.IDToMangaStore = cache.New("id-to-manga")
	options.AccessTokenStore = cache.New("access-token")

	return libmangal.NewAnilist(options)
}

type runCmd struct {
	Script   string            `arg:"" help:"Lua script string to execute. See wiki for more" required:""`
	Vars     map[string]string `help:"Variables to pass to the exec script"`
	Provider string            `help:"Path to the lua provider" optional:"" type:"existingfile"`
}

func (r *runCmd) Run() error {
	anilist := newAnilist()

	var client *libmangal.Client

	if r.Provider != "" {
		contents, err := os.ReadFile(r.Provider)
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
		clientOptions.Anilist = &anilist
		clientOptions.FS = fs.FS
		clientOptions.Log = func(msg string) {
			log.Info(msg)
		}

		c, err := libmangal.NewClient(context.Background(), loader, clientOptions)
		if err != nil {
			return err
		}

		client = &c
	}

	return lua.Exec(context.Background(), r.Script, lua.ExecOptions{
		Client:    client,
		Anilist:   &anilist,
		Variables: r.Vars,
	})
}
