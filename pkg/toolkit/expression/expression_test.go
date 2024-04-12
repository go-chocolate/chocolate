package expression

import "testing"

func TestExpression(t *testing.T) {
	item := Item{
		"name":    "zhangsan",
		"age":     36,
		"enabled": true,
		"gender":  "male",
	}

	exp := NewExpression()
	exp.SetOperators(item.Operators()...)
	t.Log(exp.Evaluate("$name == zhangsan && $gender == male"))
	t.Log(exp.Evaluate("$age > 20 && $gender == 'female'"))
	t.Log(exp.Evaluate("$age > 24 && $gender == 'male'"))
	t.Log(exp.Evaluate("$gender == $name"))
	t.Log(exp.Evaluate("$gender < $name"))
}

func TestExpressionParse(t *testing.T) {
	{
		exp := &expression{op: map[string]*Operator{}}
		exp.SetOperators((Item{}).Operators()...)
		t.Log(exp.parse("(a == A && b > B)||(c == C && d < D)"))
	}

	//{
	//	exp := &expression{op: map[string]*Operator{}}
	//	exp.SetOperator("+", 10, CalculateFunc(func(args ...string) (string, error) {
	//		return "", nil
	//	}))
	//	exp.SetOperator("-", 10, CalculateFunc(func(args ...string) (string, error) {
	//		return "", nil
	//	}))
	//	exp.SetOperator("*", 20, CalculateFunc(func(args ...string) (string, error) {
	//		return "", nil
	//	}))
	//	exp.SetOperator("/", 20, CalculateFunc(func(args ...string) (string, error) {
	//		return "", nil
	//	}))
	//	t.Log(exp.parse("10 + 20 * 30 / 40 - 50"))
	//	// 10 20 30 * 40 / + 50 -
	//	t.Log(exp.parse("10 + 20 * 30 / (40 - 50)"))
	//	// 10 20 30 * 40 50 - / +
	//}
}
