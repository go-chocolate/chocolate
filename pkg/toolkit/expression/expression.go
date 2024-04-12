package expression

import (
	"errors"
)

// Expression 表达式运算器
// 利用后缀表达式算法进行组合运算处理，支持自定义运算符。
// eg:
// (1 + 2) * 3
// (age > 22 && gender == 'male') || (age > 20 && gender == 'female')
// 如果表达式中使用了变量，需要在自定义操作符中自行处理；
// 变量或常量与操作符之间，以及操作符与操作符之间都必须使用空格分隔；
// 表达式中括号，单引号，双引号以及空格会做特殊处理，如需要使用这类字符，请使用单引号或双引号包裹。
// 解析结果中所有空格，单引号，双引号都会被忽略，括号为限定运算优先级的特殊字符；
// 特别的，单引号需要使用双引号包裹，双引号需要使用单引号包裹。
type Expression interface {
	Evaluate(s string) (string, error)
	SetOperator(symbol string, level Level, calculator Calculator) Expression
	SetOperators(operators ...*Operator) Expression
}

type expression struct {
	op map[string]*Operator
}

func NewExpression() Expression {
	return &expression{op: map[string]*Operator{}}
}

func (e *expression) SetOperator(symbol string, level Level, calculator Calculator) Expression {
	e.op[symbol] = &Operator{
		Level:      level,
		Calculator: calculator,
		Symbol:     symbol,
	}
	return e
}

func (e *expression) SetOperators(operators ...*Operator) Expression {
	for _, op := range operators {
		e.op[op.Symbol] = op
	}
	return e
}

func (e *expression) parse(text string) ([]string, error) {
	var s1 []string //操作符栈
	var s2 []string //结果栈
	for _, element := range e.split(text) {
		switch {
		case element == "(":
			s1 = append(s1, element) //遇到左括号直接写入操作符栈
		case element == ")": // 遇到右括号弹出操作符栈写入结果栈，直到遇到左括号为止
			for i := len(s1) - 1; i >= 0; i-- {
				if s1[i] == "(" {
					s1 = s1[:i]
					break
				}
				s2 = append(s2, s1[i])
				s1 = s1[:i]
			}
		case e.op[element] != nil: //遇到操作符，判断优先级
			if len(s1) == 0 || s1[len(s1)-1] == "(" { //操作符栈为空，或者操作符栈顶元素为左括号，直接写入操作符栈
				s1 = append(s1, element)
				continue
			}
			o1 := e.op[element]
			if o1.Level > e.op[s1[len(s1)-1]].Level { //当前操作符优先级大于栈顶元素，直接写入操作符栈
				s1 = append(s1, element)
			} else {
				//当前操作符优先级小于等于栈顶元素，弹出栈顶元素写入结果栈
				for i := len(s1) - 1; i >= 0; i-- {
					//当前操作符优先级大于栈顶元素，直接写入操作符栈
					if s1[i] == "(" || o1.Level > e.op[s1[i]].Level {
						s1 = append(s1, element)
						break
					}
					//弹出操作符栈并写入结果栈
					s2 = append(s2, s1[i])
					s1 = s1[:i]
				}
				if len(s1) == 0 { //遍历s1栈仍未找到比当前操作符优先级高的元素，直接写入操作符栈
					s1 = append(s1, element)
				}
			}
		default:
			s2 = append(s2, element)
		}
	}

	//将操作符栈剩余元素写入结果栈
	for i := len(s1) - 1; i >= 0; i-- {
		s2 = append(s2, s1[i])
	}
	return s2, nil
}

func (e *expression) Evaluate(s string) (string, error) {
	exp, err := e.parse(s)
	if err != nil {
		return "", err
	}
	var s1 []string
	for _, v := range exp {
		if op := e.op[v]; op != nil {
			var er error
			if len(s1) == 0 {
				return "", errors.New("invalid expression")
			} else if len(s1) == 1 { //单目运算
				s1[0], er = op.Calculator.Calculate(s1[0])
			} else if len(s1) > 1 { //双目运算
				a := s1[len(s1)-2]
				b := s1[len(s1)-1]
				s1[len(s1)-2], er = op.Calculator.Calculate(a, b)
				s1 = s1[:len(s1)-1]
			}
			if er != nil {
				return "", er
			}
		} else {
			s1 = append(s1, v)
		}
	}
	if len(s1) != 1 {
		return "", errors.New("invalid expression")
	} else {
		return s1[0], nil
	}
}

func (e *expression) split(str string) []string {
	var results []string
	var singleQuote, doubleQuote bool
	var current []byte
	for _, v := range str {
		if singleQuote {
			if v == '\'' {
				singleQuote = false
				continue
			} else {
				current = append(current, byte(v))
				continue
			}
		}
		if doubleQuote {
			if v == '"' {
				doubleQuote = false
				continue
			} else {
				current = append(current, byte(v))
				continue
			}
		}

		switch v {
		case '"':
			doubleQuote = true
		case '\'':
			singleQuote = true
		case ' ':
			if len(current) > 0 {
				results = append(results, string(current))
			}
			current = []byte{}
		case '(':
			if len(current) > 0 {
				results = append(results, string(current))
			}
			results = append(results, "(")
			current = []byte{}
		case ')':
			if len(current) > 0 {
				results = append(results, string(current))
			}
			results = append(results, ")")
			current = []byte{}
		default:
			current = append(current, byte(v))
		}
	}
	if len(current) > 0 {
		results = append(results, string(current))
	}
	return results
}
