package gogh_test

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type jsonMatcher struct {
	x interface{}
}

func (p jsonMatcher) jsonValue(x interface{}) interface{} {
	j, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	var v interface{}
	if err := json.Unmarshal(j, &v); err != nil {
		panic(err)
	}
	return v
}

func (p jsonMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(p.jsonValue(p.x), p.jsonValue(x))
}

func (p jsonMatcher) String() string {
	return fmt.Sprintf("is equal to %v", p.x)
}
