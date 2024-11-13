package stack_test

import (
  "testing"
  stack "github.com/AlexS25/rpn/pkg/collections/stack"
)

/*
go test -v ## запустить текущие тест
go test -v -run=^Test[a-zA-Z]+ ## запустить тест по регулярке
go test -v ./...  ## запустить тесты всех модулей (т.е. вложенные модули тоже будут тетсироваться)
go test -v -count=3  ## прогнать тесты три раза
go test -v mod_name  ## тестирование определенного модуля
go test -cover  ## тестовое покрытие (test coverage), процент покрытия тестами всех условий
*/

func TestSize(t *testing.T) {
  var s stack.StackString

  got := s.Size()
  expected := 0
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }

  //s.items = append(s.items, "string one")
  s.Push("string one")
  got = s.Size()
  expected = 1
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }

  //s.items = append(s.items, "string two")
  s.Push("string two")
  got = s.Size()
  expected = 2
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }
}

func TestIsEmpty(t *testing.T) {
  var s stack.StackString

  got := s.IsEmpty()
  expected := true
  if got != expected {
	t.Errorf("IsEmpty() = %v, want %v", got, expected)
  }


  //s.items = append(s.items, "string value")
  s.Push("string value")
  got = s.IsEmpty()
  expected = false
  if got != expected {
	t.Errorf("IsEmpty() = %v, want %v", got, expected)
  }
}

func TestPush(t *testing.T) {
  var s stack.StackString

  val := "string one"
  s.Push(val)
  //got := s.items[len(s.items) - 1]
  got := s.Peek()
  expected := val
  if got != expected {
	t.Errorf("Push(%q) = %q, want %q", val, got, expected)
  }

  val = "string two"
  s.Push(val)
  //got = s.items[len(s.items) - 1]
  got = s.Peek()
  expected = val
  if got != expected {
	t.Errorf("Push(%q) = %q, want %q", val, got, expected)
  }
}

func TestPop(t *testing.T) {
  var s stack.StackString
  //s.items = append(s.items, "string three")
  //s.items = append(s.items, "string two")
  //s.items = append(s.items, "string one")
  s.Push("string three")
  s.Push("string two")
  s.Push("string one")

  cases := []struct {
	name string
	want string
  }{
	{
	  name : "first run",
	  want : "string one",
	}, {
	  name : "second run",
	  want : "string two",
	}, {
	  name : "third run",
	  want : "string three",
	}, {
	  name : "fourth run", 
	  want : "",
	},
  }

  for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
	  got := s.Pop()
	  if got != tc.want {
		t.Errorf("Pop() = %q, want %q", got, tc.want)
	  }
	})
  }
}

func TestPeek(t *testing.T) {
  var s stack.StackString
  //s.items = append(s.items, "string three")
  //s.items = append(s.items, "string two")
  //s.items = append(s.items, "string one")
  s.Push("string three")
  s.Push("string two")
  s.Push("string one")

  cases := []struct {
	name string
	want string
  }{
	{
	  name : "first run",
	  want : "string one",
	}, {
	  name : "second run",
	  want : "string two",
	}, {
	  name : "third run",
	  want : "string three",
	}, {
	  name : "fourth run", 
	  want : "",
	},
  }

  for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
	  //length := len(s.items)
	  length := s.Size()
	  got := s.Peek()

	  if got != tc.want {
		t.Errorf("Pop() = %q, want %q", got, tc.want)
	  }

	  if length != s.Size() {
		t.Errorf("After running peek() length = %d, but should be = %d", s.Size(), length)
	  }

	  if length > 0 {
		//s.items = s.items[:length - 1]
		_ = s.Pop()
	  }

	})
  }
}

