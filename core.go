package jsonassert

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (a *Asserter) pathassertf(strictMode bool, path, act, exp string) {
	a.tt.Helper()
	if act == exp {
		return
	}
	actType, err := findType(act)
	if err != nil {
		a.tt.Errorf("'actual' JSON is not valid JSON: " + err.Error())
		return
	}
	expType, err := findType(exp)
	if err != nil {
		a.tt.Errorf("'expected' JSON is not valid JSON: " + err.Error())
		return
	}

	// If we're only caring about the presence of the key, then don't bother checking any further
	if expPresence, _ := extractString(exp); expPresence == "<<PRESENCE>>" {
		if actType == jsonNull {
			a.tt.Errorf(`expected the presence of any value at '%s', but was absent`, path)
		}
		return
	}

	//新增valuate类型，只要exp类型为valuate即满足
	if expType == jsonValuate {
		actValuate, _ := extractValuate(act, false)
		expValuate, _ := extractValuate(exp, true)
		a.checkValuate(strictMode, path, actValuate, expValuate)
		return
	}
	if actType != expType {
		a.tt.Errorf("actual JSON (%s) and expected JSON (%s) were of different types at '%s'", actType, expType, path)
		return
	}
	switch actType {
	case jsonBoolean:
		actBool, _ := extractBoolean(act)
		expBool, _ := extractBoolean(exp)
		a.checkBoolean(path, actBool, expBool)
	case jsonNumber:
		actNumber, _ := extractNumber(act)
		expNumber, _ := extractNumber(exp)
		a.checkNumber(path, actNumber, expNumber)
	case jsonString:
		actString, _ := extractString(act)
		expString, _ := extractString(exp)
		a.checkString(path, actString, expString)
	case jsonObject:
		actObject, _ := extractObject(act)
		expObject, _ := extractObject(exp)
		a.checkObject(strictMode, path, actObject, expObject)
	case jsonArray:
		actArray, _ := extractArray(act)
		expArray, _ := extractArray(exp)
		a.checkArray(strictMode, path, actArray, expArray)
	}
}

func serialize(a interface{}) string {
	bytes, err := json.Marshal(a)
	if err != nil {
		// Really don't want to panic here, but I can't see a reasonable solution.
		// If this line *does* get executed then we should really investigate what kind of input was given
		panic(errors.New("unexpected failure to re-serialize nested JSON. Please raise an issue including this error message and both the expected and actual JSON strings you used to trigger this panic" + err.Error()))
	}
	return string(bytes)
}

type jsonType string

const (
	jsonString      jsonType = "string"
	jsonNumber      jsonType = "number"
	jsonBoolean     jsonType = "boolean"
	jsonNull        jsonType = "null"
	jsonObject      jsonType = "object"
	jsonArray       jsonType = "array"
	jsonValuate		jsonType = "valuate"
	jsonTypeUnknown jsonType = "unknown"
)

func findType(j string) (jsonType, error) {
	j = strings.TrimSpace(j)
	if _, err := extractValuate(j, true); err == nil && j != "{}" {
		return jsonValuate, nil
	}
	if _, err := extractString(j); err == nil {
		return jsonString, nil
	}
	if _, err := extractNumber(j); err == nil {
		return jsonNumber, nil
	}
	if j == "null" {
		return jsonNull, nil
	}
	if _, err := extractObject(j); err == nil {
		return jsonObject, nil
	}
	if _, err := extractBoolean(j); err == nil {
		return jsonBoolean, nil
	}
	if _, err := extractArray(j); err == nil {
		return jsonArray, nil
	}
	return jsonTypeUnknown, fmt.Errorf(`unable to identify JSON type of "%s"`, j)
}

// *testing.T has a Helper() func that allow testing tools like this package to
// ignore their own frames when calling Errorf on *testing.T instances.
// This interface is here to avoid breaking backwards compatibility in terms of
// the interface we expect in New.
type tt interface {
	Printer
	Helper()
}
