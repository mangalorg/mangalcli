package lua

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	lua "github.com/yuin/gopher-lua"
)

var moduleJSON = map[string]lua.LGFunction{
	"print": func(state *lua.LState) int {
		value := state.CheckAny(1)
		json, err := marshal(luaValueToGo(value))
		if err != nil {
			state.RaiseError(err.Error())
		}

		fmt.Println(json)
		return 0
	},
}

var moduleFZF = map[string]lua.LGFunction{
	"select_one": func(state *lua.LState) int {
		var values []lua.LValue
		state.CheckTable(1).ForEach(func(_, value lua.LValue) {
			values = append(values, value)
		})

		showFn := state.CheckFunction(2)

		index, err := fuzzyfinder.Find(values, func(i int) string {
			state.Push(showFn)
			state.Push(values[i])

			if err := state.PCall(1, 1, nil); err != nil {
				state.RaiseError(err.Error())
			}

			return state.Get(-1).String()
		})

		if err != nil {
			state.RaiseError(err.Error())
		}

		state.Push(values[index])
		return 1
	},
	"select_multi": func(state *lua.LState) int {
		var values []lua.LValue
		state.CheckTable(1).ForEach(func(_, value lua.LValue) {
			values = append(values, value)
		})
		showFn := state.CheckFunction(2)

		indexes, err := fuzzyfinder.FindMulti(values, func(i int) string {
			state.Push(showFn)
			state.Push(values[i])

			if err := state.PCall(1, 1, nil); err != nil {
				state.RaiseError(err.Error())
			}

			return state.Get(-1).String()
		})

		if err != nil {
			state.RaiseError(err.Error())
		}

		table := state.NewTable()
		for _, index := range indexes {
			table.Append(values[index])
		}

		state.Push(table)
		return 1
	},
}
