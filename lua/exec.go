package lua

import (
	"context"
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"reflect"
)

type ExecOptions struct {
	Client    *libmangal.Client
	Anilist   *libmangal.Anilist
	Variables map[string]string
}

func Exec(
	ctx context.Context,
	script string,
	options ExecOptions,
) error {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	state.SetContext(ctx)

	config := luar.GetConfig(state)
	config.FieldNames = func(s reflect.Type, f reflect.StructField) []string {
		return []string{caseCamelToSnake(f.Name)}
	}
	config.MethodNames = func(t reflect.Type, m reflect.Method) []string {
		return []string{caseCamelToSnake(m.Name)}
	}

	for _, injectLib := range []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
	} {
		injectLib(state)
	}

	varsTable := state.NewTable()
	for key, value := range options.Variables {
		varsTable.RawSetString(key, lua.LString(value))
	}

	state.SetGlobal("Vars", varsTable)
	state.Register(
		"SearchMangas",
		newClientDependentFunction(options.Client, newSearchMangas),
	)

	state.Register(
		"MangaVolumes",
		newClientDependentFunction(options.Client, newMangaVolumes),
	)

	state.Register(
		"VolumeChapters",
		newClientDependentFunction(options.Client, newVolumeChapters),
	)

	state.Register(
		"ChapterPages",
		newClientDependentFunction(options.Client, newChapterPages),
	)

	state.Register(
		"DownloadChapter",
		newClientDependentFunction(options.Client, newDownloadChapter),
	)

	state.RegisterModule("json", moduleJSON)
	state.RegisterModule("fzf", moduleFZF)

	state.RegisterModule("anilist", map[string]lua.LGFunction{
		"find_closest_manga": newAnilistFindClosestManga(options.Anilist),
		"search_mangas":      newAnilistSearchMangas(options.Anilist),
		"get_manga_by_id":    newAnilistGetByID(options.Anilist),
		"bind_title_to_id":   newAnilistBind(options.Anilist),
	})

	lFunction, err := state.LoadString(script)
	if err != nil {
		return err
	}

	return state.CallByParam(lua.P{
		Fn:      lFunction,
		NRet:    1,
		Protect: true,
	})
}

func newClientDependentFunction(
	client *libmangal.Client,
	fn func(*libmangal.Client) lua.LGFunction,
) lua.LGFunction {
	if client == nil {
		return func(state *lua.LState) int {
			state.RaiseError("provider not loaded, try --provider <path>")
			return 0
		}
	}

	return fn(client)
}

func newSearchMangas(client *libmangal.Client) lua.LGFunction {
	return func(state *lua.LState) int {
		query := state.CheckString(1)

		mangas, err := client.SearchMangas(state.Context(), query)

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, manga := range mangas {
			table.Append(luar.New(state, manga))
		}

		state.Push(table)
		return 1
	}
}

func newMangaVolumes(client *libmangal.Client) lua.LGFunction {
	return func(state *lua.LState) int {
		userdata := state.CheckUserData(1)
		manga, ok := userdata.Value.(libmangal.Manga)
		if !ok {
			state.ArgError(1, fmt.Sprintf("manga expected, got: %T", userdata.Value))
		}

		volumes, err := client.MangaVolumes(state.Context(), manga)

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, volume := range volumes {
			table.Append(luar.New(state, volume))
		}

		state.Push(table)
		return 1
	}
}

func newVolumeChapters(client *libmangal.Client) lua.LGFunction {
	return func(state *lua.LState) int {
		userdata := state.CheckUserData(1)
		volume, ok := userdata.Value.(libmangal.Volume)
		if !ok {
			state.ArgError(1, fmt.Sprintf("volume expected, got: %T", userdata.Value))
		}

		chapters, err := client.VolumeChapters(state.Context(), volume)

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, chapter := range chapters {
			table.Append(luar.New(state, chapter))
		}

		state.Push(table)
		return 1
	}
}

func newChapterPages(client *libmangal.Client) lua.LGFunction {
	return func(state *lua.LState) int {
		userdata := state.CheckUserData(1)
		chapter, ok := userdata.Value.(libmangal.Chapter)
		if !ok {
			state.ArgError(1, fmt.Sprintf("chapter expected, got: %T", userdata.Value))
		}

		pages, err := client.ChapterPages(state.Context(), chapter)

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, page := range pages {
			table.Append(luar.New(state, page))
		}

		state.Push(table)
		return 1
	}
}

func parseDownloadOptions(table *lua.LTable) (libmangal.DownloadOptions, error) {
	options := libmangal.DefaultDownloadOptions()

	mapping := map[string]struct {
		need  lua.LValueType
		apply func(lua.LValue) error
	}{
		"format": {
			need: lua.LTString,
			apply: func(value lua.LValue) error {
				format, err := libmangal.FormatString(string(value.(lua.LString)))
				if err != nil {
					return err
				}

				options.Format = format
				return nil
			},
		},
		"directory": {
			need: lua.LTString,
			apply: func(value lua.LValue) error {
				options.Directory = string(value.(lua.LString))
				return nil
			},
		},
		"create_manga_dir": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.CreateMangaDir = bool(value.(lua.LBool))
				return nil
			},
		},
		"create_volume_dir": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.CreateVolumeDir = bool(value.(lua.LBool))
				return nil
			},
		},
		"strict": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.Strict = bool(value.(lua.LBool))
				return nil
			},
		},
		"skip_if_exists": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.SkipIfExists = bool(value.(lua.LBool))
				return nil
			},
		},
		"download_manga_cover": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.DownloadMangaCover = bool(value.(lua.LBool))
				return nil
			},
		},
		"download_manga_banner": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.DownloadMangaBanner = bool(value.(lua.LBool))
				return nil
			},
		},
		"write_series_json": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.WriteSeriesJson = bool(value.(lua.LBool))
				return nil
			},
		},
		"write_comic_info_xml": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.WriteComicInfoXml = bool(value.(lua.LBool))
				return nil
			},
		},
		"read_after": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.ReadAfter = bool(value.(lua.LBool))
				return nil
			},
		},
		"read_incognito": {
			need: lua.LTBool,
			apply: func(value lua.LValue) error {
				options.ReadIncognito = bool(value.(lua.LBool))
				return nil
			},
		},
		"comic_info_options": {
			need: lua.LTTable,
			apply: func(value lua.LValue) error {
				table := value.(*lua.LTable)
				err := gluamapper.Map(table, &options.ComicInfoOptions)
				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	var values = make(map[string]lua.LValue)

	table.ForEach(func(key, value lua.LValue) {
		keyString, ok := key.(lua.LString)
		if !ok {
			return
		}

		values[string(keyString)] = value
	})

	for key, value := range values {
		mapper, ok := mapping[key]
		if !ok {
			return options, fmt.Errorf("unknown option: %s", key)
		}

		if value.Type() != mapper.need {
			return options, fmt.Errorf("expected %s, got %s", mapper.need, value.Type())
		}

		err := mapper.apply(value)
		if err != nil {
			return options, err
		}
	}

	return options, nil
}

func newDownloadChapter(client *libmangal.Client) lua.LGFunction {
	return func(state *lua.LState) int {
		userdata := state.CheckUserData(1)
		chapter, ok := userdata.Value.(libmangal.Chapter)
		if !ok {
			state.ArgError(1, fmt.Sprintf("chapter expected, got: %T", userdata.Value))
		}

		optionsTable := state.OptTable(2, state.NewTable())
		options, err := parseDownloadOptions(optionsTable)
		if err != nil {
			state.RaiseError(err.Error())
		}

		path, err := client.DownloadChapter(state.Context(), chapter, options)
		if err != nil {
			state.RaiseError(err.Error())
		}

		state.Push(lua.LString(path))
		return 1
	}
}

func newAnilistFindClosestManga(anilist *libmangal.Anilist) lua.LGFunction {
	return func(state *lua.LState) int {
		query := state.CheckString(1)

		manga, ok, err := anilist.FindClosestManga(state.Context(), query)

		if err != nil {
			state.RaiseError(err.Error())
		}

		if !ok {
			state.Push(lua.LNil)
			return 1
		}

		state.Push(luar.New(state, manga))
		return 1
	}
}

func newAnilistSearchMangas(anilist *libmangal.Anilist) lua.LGFunction {
	return func(state *lua.LState) int {
		query := state.CheckString(1)

		mangas, err := anilist.SearchMangas(state.Context(), query)

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, manga := range mangas {
			table.Append(luar.New(state, manga))
		}

		state.Push(table)
		return 1
	}
}

func newAnilistGetByID(anilist *libmangal.Anilist) lua.LGFunction {
	return func(state *lua.LState) int {
		id := state.CheckInt(1)

		manga, ok, err := anilist.GetByID(state.Context(), id)

		if err != nil {
			state.RaiseError(err.Error())
		}

		if !ok {
			state.Push(lua.LNil)
			return 1
		}

		state.Push(luar.New(state, manga))
		return 1
	}
}

func newAnilistBind(anilist *libmangal.Anilist) lua.LGFunction {
	return func(state *lua.LState) int {
		title := state.CheckString(1)
		id := state.CheckNumber(2)

		err := anilist.BindTitleWithID(title, int(id))
		if err != nil {
			state.RaiseError(err.Error())
		}

		return 0
	}
}
