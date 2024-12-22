package calculation

import "errors"

var (
	// errEmptyStack      = errors.New("invalid operation, stack is empty")
	errDivisionByZero    = errors.New("division by zero is not allowed")
	errInvalidExpression = errors.New("invalid expression: ")
	errUnknownVale       = errors.New("unknown value in expression: ")
)
