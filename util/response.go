package util

import (
	"encoding/json"
)

func GetResponseCode(response []byte, field string) (res string) {
	var m map[string]interface{}
	if err := json.Unmarshal(response, &m); err != nil {
		return
	}
	if v, ok := m[field]; ok {
		res = v.(string)
	}
	return
}
