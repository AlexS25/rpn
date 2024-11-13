package rpn

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	stack "github.com/AlexS25/rpn/pkg/collections/stack"
)

var (
	// errEmptyStack     = errors.New("invalid operation, stack is empty")
	errDivisionByZero = errors.New("division by zero is not allowed")
)

func IsNumber(val string) bool {
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		return true
	}
	return false
}

func isBrackets(val string) bool {
	if val == "(" || val == ")" {
		return true
	}
	return false
}

func isOperator(val string) bool {
	if val == "+" || val == "-" || val == "*" || val == "/" {
		return true
	}
	return false
}

func ParseExpr(expr string) ([]string, error) {
	var res []string
	var lastVal, tmp string

	for _, val := range expr {
		tmp = strings.Trim(string(val), " ")

		if len(tmp) > 0 {
			// runeVal = []rune(tmp)[0]
			switch {
			case isBrackets(tmp) || isOperator(tmp):
				if len(lastVal) > 0 {
					res = append(res, lastVal)
				}
				lastVal = tmp
			case !IsNumber(lastVal) && IsNumber(tmp):
				if len(lastVal) > 0 {
					res = append(res, lastVal)
				}
				lastVal = tmp
			case IsNumber(tmp) && IsNumber(lastVal):
				lastVal += tmp
			default:
				return nil, errors.New("unknown value in expression:" + tmp)
			}
		}
	}
	if IsNumber(tmp) && IsNumber(lastVal) || !IsNumber(tmp) {
		res = append(res, lastVal)
	}

	return res, nil
}

func CheckBrackets(vals []string) (bool, error) {
	var brackets stack.StackString
	for _, val := range vals {
		if val == "(" {
			brackets.Push(val)
		}
		if val == ")" {
			if !brackets.IsEmpty() {
				_ = brackets.Pop()
			} else {
				return false, errors.New("not found open bracket in expression: " + strings.Join(vals, ""))
			}
		}
	}

	if !brackets.IsEmpty() {
		return false, errors.New("not found close bracket in expression: " + strings.Join(vals, ""))
	}

	return true, nil
}

func CheckSyntax(expression []string) (bool, error) {
	var stackValid stack.StackString
	for _, val := range expression {
		switch {
		case isBrackets(val):
			continue
		case isOperator(val):
			if !stackValid.IsEmpty() {
				_ = stackValid.Pop()
			} else {
				return false, errors.New("invalid expression: " + strings.Join(expression, ""))
			}
		default:
			if IsNumber(val) {
				stackValid.Push(val)
			}
		}
	}

	if stackValid.Size() != 1 {
		return false, errors.New("invalid expression: " + strings.Join(expression, ""))
	}

	return true, nil
}

func getPriority(val string) int {
	switch val {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}

	return 0
}

func convertToPostfix(list []string) []string {
	var postfixExpresion []string
	var operationStack stack.StackString

	for _, val := range list {
		switch {
		case IsNumber(val):
			postfixExpresion = append(postfixExpresion, val)
		case val == "(":
			operationStack.Push(val)
		case val == ")":
			for !operationStack.IsEmpty() && operationStack.Peek() != "(" {
				postfixExpresion = append(postfixExpresion, operationStack.Pop())
			}
			_ = operationStack.Pop()
		case isOperator(val):
			for !operationStack.IsEmpty() && getPriority(operationStack.Peek()) >= getPriority(val) {
				postfixExpresion = append(postfixExpresion, operationStack.Pop())
			}
			operationStack.Push(val)
		}
	}

	for !operationStack.IsEmpty() {
		postfixExpresion = append(postfixExpresion, operationStack.Pop())
	}

	return postfixExpresion
}

func evaluate(postfixExpresion []string) (float64, error) {
	var valueStack stack.StackString
	var res float64

	for _, val := range postfixExpresion {

		switch {
		case IsNumber(val):
			valueStack.Push(val)
		case isOperator(val):
			right, _ := strconv.ParseFloat(valueStack.Pop(), 64)
			left, _ := strconv.ParseFloat(valueStack.Pop(), 64)

			switch {
			case val == "+":
				res = left + right
			case val == "-":
				res = left - right
			case val == "*":
				res = left * right
			case val == "/":
				if right == 0 {
					return 0, errDivisionByZero
				}
				res = left / right
			}

			valueStack.Push(fmt.Sprintf("%f", res))
		}
	}
	res, _ = strconv.ParseFloat(valueStack.Pop(), 64)

	return res, nil
}

func CalcNew(expression string) (float64, error) {

	parsedVals, err := ParseExpr(expression)
	if err != nil {
		return 0, err
	}
	// fmt.Println(parsedVals)

	if res, err := CheckBrackets(parsedVals); !res {
		return 0, err
	}

	if res, err := CheckSyntax(parsedVals); !res {
		return 0, err
	}

	convertedVals := convertToPostfix(parsedVals)
	// fmt.Println(convertedVals)

	res, err := evaluate(convertedVals)
	if err != nil {
		return 0, err
	}
	return res, nil

}
