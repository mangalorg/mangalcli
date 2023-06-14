package cmd

import (
	"context"
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangalcli/cache"
	"os"
)

func newAnilist(ctx *Context) libmangal.Anilist {
	options := libmangal.DefaultAnilistOptions()
	options.Log = func(msg string) {
		fmt.Fprintln(ctx.LogWriter, msg)
	}

	options.QueryToIDsStore = cache.New("query-to-id")
	options.TitleToIDStore = cache.New("title-to-id")
	options.IDToMangaStore = cache.New("id-to-manga")
	options.AccessTokenStore = cache.New("access-token")

	return libmangal.NewAnilist(options)
}

type anilistCmd struct {
	Bind   anilistBindCmd   `cmd:""`
	Search anilistSearchCmd `cmd:""`
	Get    anilistGetCmd    `cmd:""`
}

type anilistGetCmd struct {
	Id int `arg:"" help:"anilist manga id"`
}

func (a *anilistGetCmd) Run(ctx *Context) error {
	anilist := newAnilist(ctx)

	manga, ok, err := anilist.GetByID(context.Background(), a.Id)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("not found")
	}

	marshalled, err := json.MarshalToString(manga)
	if err != nil {
		return err
	}

	fmt.Println(marshalled)
	return nil
}

type anilistSearchCmd struct {
	Query string `arg:"" help:"search query"`
}

func (a *anilistSearchCmd) Run(ctx *Context) error {
	anilist := newAnilist(ctx)

	mangas, err := anilist.SearchMangas(context.Background(), a.Query)
	if err != nil {
		return err
	}

	marshalled, err := json.MarshalToString(mangas)
	if err != nil {
		return err
	}

	fmt.Println(marshalled)
	return nil
}

type anilistBindCmd struct {
	Title string `help:"title to bind"`
	Id    int    `help:"id to bind"`
}

func (a *anilistBindCmd) Run(ctx *Context) error {
	anilist := newAnilist(ctx)

	err := anilist.BindTitleWithID(a.Title, a.Id)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "successfully binded %q => %d\n", a.Title, a.Id)
	return nil
}
