package tools

import (
	"errors"
	"github.com/tidwall/gjson"
	"testing"
)

const json = `{"name": {"first":"Bob", "last":"James"}}`

func TestGet(t *testing.T) {
	// 判断是否是json
	if !gjson.Valid(json) {
		err := errors.New("invalid json")
		if err != nil {
			return
		}
	}
	value := gjson.Get(json, "name.first")
	println(value.String())

	m, ok := gjson.Parse(json).Value().(map[string]interface{})
	if !ok {
		// not a map
		return
	}
	println(m)
}
