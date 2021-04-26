package jsonassert_test

import (
	"github.com/545080079/jsonassert"
	"testing"
)


func TestJsonValuate1(t *testing.T) {

	ja := jsonassert.New(t)
	ja.Assertf(false, `{"b": 25}`, `
		{
			"a": "@notExists()",
			"b": "@ >= 25"
		}`)
}

func TestEmpty(t *testing.T) {

	ja := jsonassert.New(t)
	ja.Assertf(true, `{"a":12, "b": "1"}`, `
		{
			"a": "@notEmpty()",
			"b": "@notEmpty()"
		}`)
}

func TestLen(t *testing.T) {

	ja := jsonassert.New(t)
	ja.Assertf(true, `{"a":12, "b": "25"}`, `
		{
			"a": "@len() >= 1",
			"b": "@len() < 3"
		}`)
}

func TestExist(t *testing.T) {

	ja := jsonassert.New(t)
	ja.Assertf(true, `{"a": 111, "b": 25}`, `
		{
			"a": "@exists()",
			"b": "@ >= 25"
		}`)
}

func TestNotExists(t *testing.T) {

	ja := jsonassert.New(t)
	ja.Assertf(false, `{"b": 25}`, `
		{
			"a": "@notExists()",
			"b": "@ >= 25"
		}`)
}