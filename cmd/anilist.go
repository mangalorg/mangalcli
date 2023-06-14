package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/mangalorg/libmangal"
	"os"
)

func newAnilist() libmangal.Anilist {
	options := libmangal.DefaultAnilistOptions()
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

func (a *anilistGetCmd) Run() error {
	anilist := newAnilist()

	manga, ok, err := anilist.GetByID(context.Background(), a.Id)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("not found")
	}

	marshalled, err := marshal(manga)
	if err != nil {
		return err
	}

	fmt.Println(marshalled)
	return nil
}

type anilistSearchCmd struct {
	Query string `arg:"" help:"search query"`
}

func (a *anilistSearchCmd) Run() error {
	anilist := newAnilist()

	mangas, err := anilist.SearchMangas(context.Background(), a.Query)
	if err != nil {
		return err
	}

	marshalled, err := marshal(mangas)
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

func (a *anilistBindCmd) Run() error {
	anilist := newAnilist()

	err := anilist.BindTitleWithID(a.Title, a.Id)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "successfully binded %q => %d\n", a.Title, a.Id)
	return nil
}
