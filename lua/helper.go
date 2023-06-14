package lua

import (
	json "github.com/json-iterator/go"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	lua "github.com/yuin/gopher-lua"
)

func marshal(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func luaValueToGo(value lua.LValue) any {
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTNumber:
		return float64(value.(lua.LNumber))
	case lua.LTString:
		return string(value.(lua.LString))
	case lua.LTTable:
		table := value.(*lua.LTable)
		om := orderedmap.New[any, any]()

		var asMap = make(map[any]any)

		table.ForEach(func(key lua.LValue, value lua.LValue) {
			k, v := luaValueToGo(key), luaValueToGo(value)
			asMap[k] = v
			om.Set(k, v)
		})

		// check if we can convert table to slice.
		// if not, return as map.
		var (
			prev    float64 = 0
			asSlice []any
		)
		for pair := om.Oldest(); pair != nil; pair = pair.Next() {
			asNum, ok := pair.Key.(float64)
			if !ok || asNum != prev+1 {
				return asMap
			}

			prev = asNum
			asSlice = append(asSlice, pair.Value)
		}

		return asSlice
	case lua.LTUserData:
		return value.(*lua.LUserData).Value
	default:
		return nil
	}
}
