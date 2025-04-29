package utils

import "encoding/json"

func ToJSON(data any) []byte {
	d, _ := json.Marshal(data)
	return d
}
