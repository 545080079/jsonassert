package jsonassert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (a *Asserter) checkObject(strictMode bool, path string, act, exp map[string]interface{}) {
	a.tt.Helper()

	//严格模式检查内容
	if strictMode {
		if len(act) != len(exp) {
			a.tt.Errorf("expected %d keys at '%s' but got %d keys", len(exp), path, len(act))
		}
		if unique := difference(act, exp); len(unique) != 0 {
			a.tt.Errorf("unexpected object key(s) %+v found at '%s'", serialize(unique), path)
		}
		if unique := difference(exp, act); len(unique) != 0 {
			a.tt.Errorf("expected object key(s) %+v missing at '%s'", serialize(unique), path)
		}
	}

	for key := range act {
		if contains(exp, key) {
			a.pathassertf(strictMode, path+"."+key, serialize(act[key]), serialize(exp[key]))
		}
	}
}

func difference(act, exp map[string]interface{}) []string {
	unique := []string{}
	for key := range act {
		if !contains(exp, key) {
			unique = append(unique, key)
		}
	}
	return unique
}

func contains(container map[string]interface{}, candidate string) bool {
	for key := range container {
		if key == candidate {
			return true
		}
	}
	return false
}

func extractObject(s string) (map[string]interface{}, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot parse empty string as object")
	}
	if s[0] != '{' {
		return nil, fmt.Errorf("cannot parse '%s' as object", s)
	}
	var arr map[string]interface{}
	err := json.Unmarshal([]byte(s), &arr)
	return arr, err
}
