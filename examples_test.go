package jsonassert_test

import (
	"fmt"
	"testing"

	"github.com/545080079/jsonassert"
)

type printer struct{}

func (p *printer) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// using the varible name 't' to mimic a *testing.T variable
var t *printer

func ExampleNew() {
	ja := jsonassert.New(t)
	ja.Assertf(true, `{"hello":"world"}`, `
		{
			"hello": "world"
		}`)
}

func ExampleAsserter_Assertf_formatArguments() {
	ja := jsonassert.New(t)
	ja.Assertf(true, `{"hello":"世界"}`, `
		{
			"hello": "%s"
		}`, "world")
	//output:
	//expected string at '$.hello' to be 'world' but was '世界'
}

func ExampleAsserter_Assertf_presenceOnly() {
	ja := jsonassert.New(t)
	ja.Assertf(true, `{"asdf":"not the right key name"}`, `
		{
			"asdf": "<<PRESENCE>>"
		}`)
	//output:
	//unexpected object key(s) ["hi"] found at '$'
	//expected object key(s) ["hello"] missing at '$'
}


func TestX(t *testing.T) {

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