/**
 * @Author: yutaoluo@tencent.com
 * @Description: 支持JSON表达式、函数式断言
 * @File: valuate
 * @Date: 2021/4/25 19:01
 */

package jsonassert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"strconv"
	"strings"
)

/*
	标识该字段为表达式jsonValuate类型
	形如:
		exp: @>10
		act: @12>10
	支持@>, @<, @>=, @<=, @==,
	@len()			:该key对应的value长度必须满足后续给出的条件,如 @len()>10,
	@notEmpty()		:该key对应的value非空
	@notExists()	:该key必须不存在
	@exists()		:该key必须存在
 */
const (
	preFlag				= "@"
	preFlagLen 			= "@len()"
	preFlagNotEmpty 	= "@notEmpty()"
	preFlagNotExists	= "@notExists()"
	preFlagExists		= "@exists()"
)

func (a *Asserter) checkValuate(strictMode bool, path string, act, exp map[string]interface{}) {
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

	for k, _ := range exp {
		//获取用户输入的valuate表达式
		expStr := processInput(exp[k])
		if expStr == "" && exp[k] != "" {
			a.tt.Errorf(`expected act type is string or float64 at '%s', but not.`, path)
		}

		if strings.HasPrefix(expStr, preFlagNotExists) {
			if _, ok := act[k]; ok {
				a.tt.Errorf(`expected not exists at '%s', but contains in it.`, path)
			}
			//无需判断value内容, 提前返回
			return
		}
		if strings.HasPrefix(expStr, preFlagExists) {
			if _, ok := act[k]; !ok {
				a.tt.Errorf(`expected the presence of any value at '%s', but was absent`, path)
			}
			//无需判断value内容, 提前返回
			return
		}

		//处理实际输入
		actValue := processInput(act[k])
		if actValue == "" && act[k] != nil {
			a.tt.Errorf(`expected act type is string or float64 at '%s', but not.`, path)
		}

		//构造go valuate表达式输入
		if strings.HasPrefix(expStr, preFlagNotEmpty) {
			//	'@notEmpty()' --> len(actValue)>0
			expStr = strings.Replace(expStr, preFlagNotEmpty, strconv.Itoa(len(actValue)), -1) + ">0"

		} else if strings.HasPrefix(expStr, preFlagLen) {
			// '@len()' --> len(actValue)
			expStr = strings.Replace(expStr, preFlagLen, strconv.Itoa(len(actValue)), -1)

		} else if strings.HasPrefix(expStr, "@") {
			//	'@' --> act key
			expStr = strings.Replace(expStr, preFlag, k, -1)

		} else {
			expStr = expStr + "==" + actValue
		}

		res, err := evaluate(expStr, act)
		if err != nil || !res {
			a.tt.Errorf("expected valuate at '%s' to be %v but was %v", path, exp, act)
		}
	}
}

func extractValuate(s string, isExp bool) (map[string]interface{}, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot parse nothing as string")
	}
	if s[0] != '{' {
		return nil, fmt.Errorf("cannot parse '%s' as object", s)
	}

	var arr map[string]interface{}
	err := json.Unmarshal([]byte(s), &arr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse '%s' as object", s)
	}

	//判断key-value的value是否为@开头
	fmt.Println(arr)
	for _, v := range arr {
		value, valueType := v.(string)
		if isExp && !valueType {
			continue
		}
		if isExp && value[0] == preFlag[0] {
			return arr, err
		}
	}
	if isExp {
		return nil, fmt.Errorf("this field is not valuate")
	}
	return arr, err
}

//解析表达式
func evaluate(expStr string, act map[string]interface{}) (bool, error) {

	expBuf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(expBuf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(expStr)
	if err != nil {
		return false, errors.New(fmt.Sprintf("expected valuate at '%s', but is %v", expStr, act))
	}
	expStr = strings.TrimSpace(expBuf.String())
	expStr = strings.Trim(expStr, "\"")
	expression, _ := govaluate.NewEvaluableExpression(expStr)
	res, err := expression.Evaluate(act)
	if err != nil {
		return false, errors.New(fmt.Sprintf("expected valuate at '%s', but is %v", expStr, act))
	}
	if _, t := res.(bool); !t {
		return false, errors.New(fmt.Sprintf("expected valuate at '%s', but is %v", expStr, act))
	}
	return res.(bool), nil
}


//校验输入, 如果是float64, 转为string
func processInput(input interface{}) string {
	if input == nil {
		return ""
	}
	_, t1 := input.(float64)
	_, t2 := input.(string)
	if !t1 && !t2 {
		return ""
	}
	if t1 {
		return strconv.FormatFloat(input.(float64), 'f', -1, 64)
	}
	return input.(string)
}