package dbtypes

import "encoding/json"

type JSONString struct {
	Data interface{}
}

func (js *JSONString) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &js.Data)
}
