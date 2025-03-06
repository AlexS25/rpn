package calculation

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	stack "github.com/AlexS25/rpn/pkg/collections/stack"
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

func IsOperator(val string) bool {
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
			case isBrackets(tmp) || IsOperator(tmp):
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
				return nil, fmt.Errorf("unknown value %q in expression: %q",
					tmp, expr)
				//return nil, errors.New(errUnknownVale.Error() + tmp)
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
		case IsOperator(val):
			if !stackValid.IsEmpty() {
				_ = stackValid.Pop()
			} else {
				return false, errors.New(errInvalidExpression.Error() +
					"extra operator in expression: " +
					strings.Join(expression, ""))
			}
		default:
			if IsNumber(val) {
				stackValid.Push(val)
			}
		}
	}

	if stackValid.Size() != 1 {
		return false, errors.New(errInvalidExpression.Error() +
			strings.Join(expression, ""))
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

func ConvertToPostfix(list []string) []string {
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
		case IsOperator(val):
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

func EvalSimpleExpr(arg1, arg2, operation string) (string, error) {
	var res float64

	fmt.Printf("Arg1: %q, arg2: %q, operation: %q\n", arg1, arg2, operation)

	left, err := strconv.ParseFloat(arg1, 64)
	if err != nil {
		return "", fmt.Errorf("invalid parse to float value: %q. %v", arg1, err)
	}

	right, err := strconv.ParseFloat(arg2, 64)
	if err != nil {
		return "", fmt.Errorf("invalid parse to float value: %q. %v", arg1, err)
	}

	if !IsOperator(operation) {
		return "", fmt.Errorf("value %q does not apply to operation symbol", operation)
	}

	val := os.Getenv("TIME_ADDITION_MS")
	TimeAddition, _ := strconv.Atoi(val)
	if TimeAddition == 0 {
		TimeAddition = 100
	}

	val = os.Getenv("TIME_SUBTRACTION_MS")
	TimeSubtraction, _ := strconv.Atoi(val)
	if TimeSubtraction == 0 {
		TimeSubtraction = 100
	}

	val = os.Getenv("TIME_MULTIPLICATIONS_MS")
	TimeMultiply, _ := strconv.Atoi(val)
	if TimeMultiply == 0 {
		TimeMultiply = 100
	}

	val = os.Getenv("TIME_DIVISIONS_MS")
	TimeDivision, _ := strconv.Atoi(val)
	if TimeDivision == 0 {
		TimeDivision = 100
	}

	switch operation {
	case "+":
		res = left + right
		time.Sleep(time.Duration(TimeAddition) * time.Microsecond)
	case "-":
		res = left - right
		time.Sleep(time.Duration(TimeSubtraction) * time.Microsecond)
	case "*":
		res = left * right
		time.Sleep(time.Duration(TimeMultiply) * time.Microsecond)
	case "/":
		if right == 0 {
			return "", errDivisionByZero
		}
		res = left / right
		time.Sleep(time.Duration(TimeDivision) * time.Microsecond)
	}

	return fmt.Sprintf("%f", res), nil
}

func evaluate(postfixExpresion []string) (float64, error) {
	var valueStack stack.StackString
	var res float64

	for _, val := range postfixExpresion {
		//fmt.Printf("==> `evaluate` pos %v, val is: %q\n", i, val)

		switch {
		case IsNumber(val):
			valueStack.Push(val)
		case IsOperator(val):
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

func GetPostfixExpr(expression string) ([]string, error) {
	parsedVals, err := ParseExpr(expression)
	if err != nil {
		return nil, err
	}
	// fmt.Println(parsedVals)

	if res, err := CheckBrackets(parsedVals); !res {
		if err != nil {
			return nil, errors.New(errInvalidExpression.Error() + err.Error())
		}
		return nil, errors.New(errInvalidExpression.Error() + expression)
	}

	if res, err := CheckSyntax(parsedVals); !res {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errInvalidExpression.Error() + expression)
	}

	return ConvertToPostfix(parsedVals), nil
}

func Calc(expression string) (float64, error) {

	convertedVals, err := GetPostfixExpr(expression)
	if err != nil {
		return 0, err
	}
	// fmt.Println("Converted data: ", convertedVals)

	res, err := evaluate(convertedVals)
	if err != nil {
		return 0, err
	}
	return res, nil

}
