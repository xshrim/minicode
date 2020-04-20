package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GetString(jsonStr, key string) string {
	return gjson.Get(jsonStr, key).String()
}

func GetInt(jsonStr, key string) int64 {
	return gjson.Get(jsonStr, key).Int()
}
