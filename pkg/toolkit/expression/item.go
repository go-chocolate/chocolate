package expression

import (
	"fmt"
)

const (
	Integer Type = iota
	Float
	String
	Bool
	Others
)

type Type int

type Item map[string]any

func (i Item) Operators() []*Operator {
	return []*Operator{
		{Symbol: "==", Level: 20, Calculator: CalculatorEQ(i)},
		{Symbol: "!=", Level: 20, Calculator: CalculatorNE(i)},
		{Symbol: "<", Level: 20, Calculator: CalculatorLT(i)},
		{Symbol: "<=", Level: 20, Calculator: CalculatorLTE(i)},
		{Symbol: ">", Level: 20, Calculator: CalculatorGT(i)},
		{Symbol: ">=", Level: 20, Calculator: CalculatorGTE(i)},
		{Symbol: "&&", Level: 15, Calculator: CalculatorAND()},
		{Symbol: "||", Level: 10, Calculator: CalculatorOR()},
	}
}

func (i Item) get(key string) (any, Type, error) {
	if key == "" {
		return nil, Others, fmt.Errorf("empty key")
	}
	if key[0] == '$' {
		key = key[1:]
	} else {
		return key, String, nil
	}

	v, ok := i[key]
	if !ok {
		return nil, Others, fmt.Errorf("field '%s' not found", key)
	}
	switch val := v.(type) {
	case int:
		return int(val), Integer, nil
	case int8:
		return int(val), Integer, nil
	case int16:
		return int(val), Integer, nil
	case int32:
		return int(val), Integer, nil
	case int64:
		return int(val), Integer, nil
	case uint:
		return int(val), Integer, nil
	case uint8:
		return int(val), Integer, nil
	case uint16:
		return int(val), Integer, nil
	case uint32:
		return int(val), Integer, nil
	case uint64:
		return int(val), Integer, nil
	case float64:
		return float64(val), Float, nil
	case float32:
		return float64(val), Float, nil
	case string:
		return val, String, nil
	case bool:
		return val, Bool, nil
	default:
		return v, Others, fmt.Errorf("unknown type for field '%s'", key)
	}
}

type Items interface {
	SetItems(items ...map[string]any) Items
	SetExpression(expression string) Items
	Evaluate() ([]string, error)
}

func NewItems() Items {
	return &items{}
}

type items struct {
	items      []Item
	expression string
}

func (e *items) SetItems(items ...map[string]any) Items {
	for _, v := range items {
		e.items = append(e.items, v)
	}
	return e
}
func (e *items) SetExpression(expression string) Items {
	e.expression = expression
	return e
}

func (e *items) Evaluate() ([]string, error) {
	var results = make([]string, 0, len(e.items))
	for _, item := range e.items {
		exp := NewExpression().SetOperators(item.Operators()...)
		result, err := exp.Evaluate(e.expression)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
