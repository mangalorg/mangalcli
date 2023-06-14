package cmd

import (
	"encoding/json"
)

func marshal(value any) (string, error) {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
