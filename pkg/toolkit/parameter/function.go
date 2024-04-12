package parameter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-chocolate/chocolate/pkg/toolkit/expression"
)

type function func(args string) (string, error)

var exp = expression.NewExpression()

func init() {
	exp.SetOperators(
		&expression.Operator{
			Symbol: "+",
			Calculator: expression.CalculateFunc(func(args ...string) (string, error) {
				switch len(args) {
				case 1:
					return args[0], nil
				case 2:
					a, b := args[0], args[1]
					if strings.Count(a, ".") == 1 || strings.Count(b, ".") == 1 {
						a1, err1 := strconv.ParseFloat(a, 64)
						a2, err2 := strconv.ParseFloat(b, 64)
						if err1 == nil && err2 == nil {
							return fmt.Sprintf("%f", a1+a2), nil
						}
					}
					a1, err1 := strconv.ParseInt(a, 10, 64)
					a2, err2 := strconv.ParseInt(b, 10, 64)
					if err1 == nil && err2 == nil {
						return fmt.Sprintf("%d", a1+a2), nil
					}
					return a + b, nil
				default:
					return "", fmt.Errorf("invalid args count for operator '+', expected 1 or 2 args,but got %d", len(args))
				}
			}),
			Level: 1,
		},
		&expression.Operator{
			Symbol: "-",
			Calculator: expression.CalculateFunc(func(args ...string) (string, error) {
				switch len(args) {
				case 1:
					return "-" + args[0], nil
				case 2:
					a, b := args[0], args[1]
					if strings.Count(a, ".") == 1 || strings.Count(b, ".") == 1 {
						a1, err1 := strconv.ParseFloat(a, 64)
						a2, err2 := strconv.ParseFloat(b, 64)
						if err1 == nil && err2 == nil {
							return fmt.Sprintf("%f", a1-a2), nil
						}
					}
					a1, err1 := strconv.ParseInt(a, 10, 64)
					a2, err2 := strconv.ParseInt(b, 10, 64)
					if err1 == nil && err2 == nil {
						return fmt.Sprintf("%d", a1-a2), nil
					}
					return "", fmt.Errorf("args for operator '-' must be number: '%s - %s'", a, b)
				default:
					return "", fmt.Errorf("invalid args count for operator '-', expected 1 or 2 args,but got %d", len(args))
				}
			}),
			Level: 1,
		},
		&expression.Operator{
			Symbol: "*",
			Calculator: expression.CalculateFunc(func(args ...string) (string, error) {
				switch len(args) {
				case 2:
					a, b := args[0], args[1]
					if strings.Count(a, ".") == 1 || strings.Count(b, ".") == 1 {
						a1, err1 := strconv.ParseFloat(a, 64)
						a2, err2 := strconv.ParseFloat(b, 64)
						if err1 == nil && err2 == nil {
							return fmt.Sprintf("%f", a1*a2), nil
						}
					}
					a1, err1 := strconv.ParseInt(a, 10, 64)
					a2, err2 := strconv.ParseInt(b, 10, 64)
					if err1 == nil && err2 == nil {
						return fmt.Sprintf("%d", a1*a2), nil
					}
					return "", fmt.Errorf("args for operator '*' must be number: '%s * %s'", a, b)
				default:
					return "", fmt.Errorf("invalid args count for operator '*', expected 2 args,but got %d", len(args))
				}
			}),
			Level: 2,
		},
		&expression.Operator{
			Symbol: "/",
			Calculator: expression.CalculateFunc(func(args ...string) (string, error) {
				switch len(args) {
				case 2:
					a, b := args[0], args[1]
					a1, err1 := strconv.ParseFloat(a, 64)
					a2, err2 := strconv.ParseFloat(b, 64)
					if err1 == nil && err2 == nil {
						return fmt.Sprintf("%f", a1/a2), nil
					}
					return "", fmt.Errorf("args for operator '/' must be number: '%s * %s'", a, b)
				default:
					return "", fmt.Errorf("invalid args count for operator '/', expected 2 args,but got %d", len(args))
				}
			}),
			Level: 2,
		},
	)
}

var _functions_ = map[string]function{
	"_exp_": func(args string) (string, error) {
		if args != "" {
			return "", nil
		}
		return exp.Evaluate(args)
	},
	"_now_": func(args string) (string, error) {
		now := time.Now()
		if args != "" {
			return now.Format(args), nil
		}
		return now.Format(time.DateTime), nil
	},
	"_datetime_": func(args string) (string, error) {
		now := time.Now()
		if args == "" {
			return now.Format(time.DateTime), nil
		}
		layout := time.DateTime
		duration := time.Duration(0)
		params := splitArgs(args)
		if len(params) > 0 && params[0] != "" {
			var err error
			lastChar := params[0][len(params[0])-1]
			if lastChar >= '0' && lastChar <= '9' {
				v, err := strconv.ParseInt(params[0], 10, 64)
				if err != nil {
					return "", err
				}
				duration = time.Duration(v) * time.Second
			} else {
				duration, err = time.ParseDuration(params[0])
				if err != nil {
					return "", err
				}
			}
		}
		if len(params) > 1 && params[1] != "" {
			layout = params[1]
		}
		return now.Add(duration).Format(layout), nil
	},
}

func splitArgs(args string) []string {
	var result []string
	for _, v := range strings.Split(args, ",") {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}

func SetFunction(name string, f function) {
	_functions_[name] = f
}
