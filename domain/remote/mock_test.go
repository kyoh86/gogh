package remote_test

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type jsonMatcher struct {
	x any
}

func (p jsonMatcher) jsonValue(x any) any {
	j, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	var v any
	if err := json.Unmarshal(j, &v); err != nil {
		panic(err)
	}
	return v
}

func (p jsonMatcher) Matches(x any) bool {
	return reflect.DeepEqual(p.jsonValue(p.x), p.jsonValue(x))
}

func (p jsonMatcher) String() string {
	return fmt.Sprintf("is equal to %v", p.x)
}
