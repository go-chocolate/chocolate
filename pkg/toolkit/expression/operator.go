package expression

import (
	"errors"
	"fmt"

	"github.com/go-chocolate/chocolate/pkg/toolkit/types"
)

// Level 运算符优先级
type Level float64

//const (
//	LevelAddSub Level = 10.0 // 加减法优先级
//	LevelMulDiv Level = 20.0 // 乘除法优先级
//)

const (
	EQ  = "=="
	NE  = "!="
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="
	AND = "&&"
	OR  = "||"
)

type Operator struct {
	Symbol     string
	Calculator Calculator
	Level      Level
}

type Calculator interface {
	Calculate(args ...string) (string, error)
}

type CalculateFunc func(args ...string) (string, error)

func (f CalculateFunc) Calculate(args ...string) (string, error) {
	return f(args...)
}

func CalculatorEQ(item Item) Calculator {
	return &calculator{item: item, op: EQ}
}
func CalculatorNE(item Item) Calculator {
	return &calculator{item: item, op: NE}
}
func CalculatorLT(item Item) Calculator {
	return &calculator{item: item, op: LT}
}
func CalculatorLTE(item Item) Calculator {
	return &calculator{item: item, op: LTE}
}
func CalculatorGT(item Item) Calculator {
	return &calculator{item: item, op: GT}
}
func CalculatorGTE(item Item) Calculator {
	return &calculator{item: item, op: GTE}
}

func CalculatorAND() CalculateFunc {
	return func(args ...string) (string, error) {
		if len(args) != 2 {
			return "", errors.New("invalid args")
		}
		left, right := args[0], args[1]
		if !isBool(left, right) {
			return "", fmt.Errorf("invalid expression: %s && %s", left, right)
		}
		return bool2string(left == "true" && right == "true"), nil
	}
}

func CalculatorOR() CalculateFunc {
	return func(args ...string) (string, error) {
		left, right := args[0], args[1]
		if !isBool(left, right) {
			return "", fmt.Errorf("invalid expression: %s || %s", left, right)
		}
		return bool2string(left == "true" || right == "true"), nil
	}
}

//func CalculatorADD() CalculateFunc {
//	return func(args ...string) (string, error) {
//		switch len(args) {
//		case 1:
//
//		case 2:
//
//		default:
//			return "", fmt.Errorf("invalid arguments for ADD")
//		}
//	}
//}
//
//func CalculatorSUB() CalculateFunc {
//	return func(args ...string) (string, error) {
//
//	}
//}
//func CalculatorMUL() CalculateFunc {
//	return func(args ...string) (string, error) {
//
//	}
//}
//func CalculatorDIV() CalculateFunc {
//	return func(args ...string) (string, error) {
//
//	}
//}

type calculator struct {
	op   string
	item Item
}

func (c *calculator) Calculate(args ...string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("invalid args")
	}

	left, t1, err := c.item.get(args[0])
	if err != nil {
		return "", err
	}

	right, _, err := c.item.get(args[1])
	if err != nil {
		return "", err
	}

	switch c.op {
	case EQ:
		return bool2string(eq(left, right, t1)), nil
	case NE:
		return bool2string(ne(left, right, t1)), nil
	case LT:
		return bool2string(lt(left, right, t1)), nil
	case LTE:
		return bool2string(lte(left, right, t1)), nil
	case GT:
		return bool2string(gt(left, right, t1)), nil
	case GTE:
		return bool2string(gte(left, right, t1)), nil
	default:
		return "", errors.New("invalid operator")
	}
}

func eq(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) == types.AnyToInt(v2)
	case Float:
		return v1.(float64) == types.AnyToFloat64(v2)
	case String:
		return v1.(string) == types.AnyToString(v2)
	case Bool:
		return v1.(bool) == types.AnyToBool(v2)
	default:
		return false
	}
}

func ne(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) != types.AnyToInt(v2)
	case Float:
		return v1.(float64) != types.AnyToFloat64(v2)
	case String:
		return v1.(string) != types.AnyToString(v2)
	case Bool:
		return v1.(bool) != types.AnyToBool(v2)
	default:
		return false
	}
}

func lt(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) < types.AnyToInt(v2)
	case Float:
		return v1.(float64) < types.AnyToFloat64(v2)
	case String:
		return v1.(string) < types.AnyToString(v2)
	case Bool:
		return false
	default:
		return false
	}
}

func lte(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) <= types.AnyToInt(v2)
	case Float:
		return v1.(float64) <= types.AnyToFloat64(v2)
	case String:
		return v1.(string) <= types.AnyToString(v2)
	case Bool:
		return false
	default:
		return false
	}
}

func gt(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) > types.AnyToInt(v2)
	case Float:
		return v1.(float64) > types.AnyToFloat64(v2)
	case String:
		return v1.(string) > types.AnyToString(v2)
	case Bool:
		return false
	default:
		return false
	}
}

func gte(v1, v2 any, t Type) bool {
	switch t {
	case Integer:
		return v1.(int) >= types.AnyToInt(v2)
	case Float:
		return v1.(float64) >= types.AnyToFloat64(v2)
	case String:
		return v1.(string) >= types.AnyToString(v2)
	case Bool:
		return false
	default:
		return false
	}
}

func bool2string(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func isBool(s ...string) bool {
	for _, v := range s {
		if v == "true" || v == "false" {
			continue
		}
		return false
	}
	return true
}
