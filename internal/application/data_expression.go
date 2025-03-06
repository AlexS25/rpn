package application

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/AlexS25/rpn/pkg/calculation"
	"github.com/AlexS25/rpn/pkg/collections/stack"
)

type dataExpression struct {
	srcExpr     string   // На всякий сохраняю всё выражение
	postfixExpr []string // Выражение приведенное к Postfix стандарту (из него будем брать данные)
	stackExpr   stack.StackString
	state       int // 0-free, 1-busy, 2-solved, 5-selected
}

type safeExpression struct {
	de  map[int]*dataExpression
	mtx *sync.RWMutex
}

func NewSafeExpression() *safeExpression {
	return &safeExpression{
		de:  make(map[int]*dataExpression),
		mtx: &sync.RWMutex{},
	}
}

func (se *safeExpression) Add(id int, expr *dataExpression) {
	se.mtx.Lock()
	defer se.mtx.Unlock()

	se.de[id] = expr
}

var counter = 0

func (se *safeExpression) AddExpr(expression string) (int, error) {
	postfixExpr, err := calculation.GetPostfixExpr(expression)
	if err != nil {
		return -1, err
	}

	de := &dataExpression{
		srcExpr:     expression,
		postfixExpr: postfixExpr,
		stackExpr:   stack.StackString{},
		state:       0,
	}

	counter++
	se.Add(counter, de)
	return counter, nil
}

func (se *safeExpression) Get(id int) (*dataExpression, bool) {
	se.mtx.RLock()
	defer se.mtx.RUnlock()

	expr, exists := se.de[id]
	return expr, exists
}

// func (se *safeExpression) SetStat(id int, state int) bool {
// 	se.mtx.Lock()
// 	defer se.mtx.Unlock()

// 	expr, exists := se.de[id]
// 	if !exists {
// 		return false
// 	}
// 	expr.state = state
// 	return true
// }

// func exampleSafeExpression() {
// 	safeExpr := NewSafeExpression()
// 	expr := &dataExpression{
// 		srcExpr:     "2 + 3",
// 		postfixExpr: []string{},
// 		stackExpr:   stack.StackString{},
// 		state:       0,
// 	}
// 	safeExpr.Add(1, expr)
// 	result, _ := calculation.Calc(expr.srcExpr)
// 	fmt.Println("==> Result is:", result)

// 	fmt.Println(safeExpr.Get(0))
// }

// Получаем значения для передачи таске
func (se *safeExpression) getExprForTask(id int) ([]string, error) {
	se.mtx.Lock()
	defer se.mtx.Unlock()

	var res []string
	de := se.de[id]

	if de.state != 5 {
		return nil, fmt.Errorf("This expression %q is not selected", de.srcExpr)
	}

	for _, val := range de.postfixExpr {
		// fmt.Printf("==> `evaluate` pos %v, val is: %q\n", i, val)

		de.postfixExpr = de.postfixExpr[1:]
		switch {
		case calculation.IsNumber(val):
			de.stackExpr.Push(val)
		case calculation.IsOperator(val):
			right := de.stackExpr.Pop()
			left := de.stackExpr.Pop()
			operation := val
			res = append(res, left, right, operation)
			de.state = 1 // Говорим, что заняли данный элемент
			se.de[id] = de

			return res, nil
			// valueStack.Push(fmt.Sprintf("%f", res))
		default:
			return nil, fmt.Errorf("Incorrect value: %q", val)
		}
	}
	// res, _ = strconv.ParseFloat(valueStack.Pop(), 64)
	de.state = 2 // Фиксируем итоговое решение
	se.de[id] = de

	return res, nil
}

// Записываем решение от таски обратно
func (se *safeExpression) pushValFromTask(id int, val string) error {
	if de, ok := se.Get(id); ok {
		if de.state != 1 { // Проверяем, что данны элемент был помечен, как занятый
			return fmt.Errorf("Expression with id %d, does not wait data.", id)
		}
		de.stackExpr.Push(val)
		de.state = 0 // Освобождаем наш элемент

		se.mtx.Lock()
		se.de[id] = de
		se.mtx.Unlock()

		return nil
	} else {
		return fmt.Errorf("Expression with id %d not found.", id)
	}
}

func (se *safeExpression) getFreeExprId() int {
	se.mtx.Lock()
	defer se.mtx.Unlock()

	for id, de := range se.de {
		if de.state == 0 {
			// Задаем стату, что выбрали данный элемент
			de.state = 5
			se.de[id] = de
			return id
		}
	}
	return -1 // Возвращаем отрицательное число, если нет значений для обработки

}

func (se *safeExpression) getSolution(id int) (float64, bool) {
	// se.mtx.RLock()
	// defer se.mtx.RUnlock()
	// if de, ok := [id]; ok {

	if de, ok := se.Get(id); ok {
		if de.state == 2 {
			val_s := de.stackExpr.Pop()
			if len(val_s) > 0 {
				val, _ := strconv.ParseFloat(val_s, 64)
				return val, true
			}
		}
	}

	return 0, false
}
