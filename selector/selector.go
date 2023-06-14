package selector

import (
	"context"
	"fmt"
	"github.com/mangalorg/libmangal"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func Select(
	client *libmangal.Client,
	vars map[string]string,
	query string,
) (any, error) {
	state := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	state.SetContext(context.Background())

	for _, injectLib := range []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
	} {
		injectLib(state)
	}

	varsTable := state.NewTable()

	for key, value := range vars {
		varsTable.RawSetString(key, lua.LString(value))
	}

	state.SetGlobal("Vars", varsTable)
	state.Register("SearchMangas", newSearchMangas(client))
	state.Register("MangaVolumes", newMangaVolumes(client))
	state.Register("VolumeChapters", newVolumeChapters(client))
	state.Register("ChapterPages", newChapterPages(client))

	lFunction, err := state.LoadString(query)
	if err != nil {
		return nil, err
	}

	err = state.CallByParam(lua.P{
		Fn:      lFunction,
		NRet:    1,
		Protect: true,
	})

	if err != nil {
		return nil, err
	}

	selected := state.Get(-1)

	switch selected.Type() {
	case lua.LTUserData:
		return selected.(*lua.LUserData).Value, nil
	case lua.LTTable:
		table := selected.(*lua.LTable)

		var values []any
		table.ForEach(func(_, value lua.LValue) {
			var v any
			if value.Type() == lua.LTUserData {
				ud := value.(*lua.LUserData)
				v = ud.Value
			} else {
				v = value
			}

			values = append(values, v)
		})

		return values, nil
	}

	if selected.Type() == lua.LTUserData {
	}

	return selected, nil
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
