package internal

import (
	"encoding/json"
)

type List[T any] []T

func (value *List[T]) UnmarshalJSON(data []byte) error {
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		*value = []T{single}
		return nil
	}

	var many []T
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}

	*value = many
	return nil
}
