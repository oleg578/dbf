package dbf

import (
	"encoding/json"
)

func AnyToJson(T any) ([]byte, error) {
	return json.Marshal(T)
}
