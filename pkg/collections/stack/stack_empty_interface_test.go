package stack

import (
  "testing"
)

/*
go test -v ## запустить текущие тест
go test -v -run=^Test[a-zA-Z]+ ## запустить тест по регулярке
go test -v ./...  ## запустить тесты всех модулей (т.е. вложенные модули тоже будут тетсироваться)
go test -v -count=3  ## прогнать тесты три раза
go test -v mod_name  ## тестирование определенного модуля
go test -cover  ## тестовое покрытие (test coverage), процент покрытия тестами всех условий
*/

func TestSizeEInter(t *testing.T) {
  var s StackEmptyInterface

  got := s.Size()
  expected := 0
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }

  s.items = append(s.items, "string one")
  got = s.Size()
  expected = 1
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }

  s.items = append(s.items, "string two")
  got = s.Size()
  expected = 2
  if got != expected {
	t.Errorf("Size() = %d, want %d", got, expected)
  }
}

func TestIsEmptyEInter(t *testing.T) {
  var s StackEmptyInterface

  got := s.IsEmpty()
  expected := true
  if got != expected {
	t.Errorf("IsEmpty() = %v, want %v", got, expected)
  }


  s.items = append(s.items, "string value")
  got = s.IsEmpty()
  expected = false
  if got != expected {
	t.Errorf("IsEmpty() = %v, want %v", got, expected)
  }
}

func TestPushEInter(t *testing.T) {
  var s StackEmptyInterface

  val := "string one"
  s.Push(val)
  got := s.items[len(s.items) - 1].(string)
  expected := val
  if got != expected {
	t.Errorf("Push(%q) = %q, want %q", val, got, expected)
  }

  val = "string two"
  s.Push(val)
  got = s.items[len(s.items) - 1].(string)
  expected = val
  if got != expected {
	t.Errorf("Push(%q) = %q, want %q", val, got, expected)
  }
}

func TestPopEInter(t *testing.T) {
  var s StackEmptyInterface
  s.items = append(s.items, "string three")
  s.items = append(s.items, "string two")
  s.items = append(s.items, "string one")

  cases := []struct {
	name string
	want interface{}
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
	  want : nil,  
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

func TestPeekEInter(t *testing.T) {
  var s StackEmptyInterface
  s.items = append(s.items, "string three")
  s.items = append(s.items, "string two")
  s.items = append(s.items, "string one")

  cases := []struct {
	name string
	want interface{}
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
	  want : nil, 
	},
  }

  for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
	  length := len(s.items)
	  got := s.Peek()

	  if got != tc.want {
		t.Errorf("Pop() = %q, want %q", got.(string), tc.want.(string))
	  }

	  if length != len(s.items) {
		t.Errorf("After running peek() length = %d, but should be = %d", len(s.items),  length)
	  }

	  if length > 0 {
		s.items = s.items[:length - 1]
	  }

	})
  }
}

