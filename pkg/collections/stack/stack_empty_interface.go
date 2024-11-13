package stack

type StackEmptyInterface struct {
  items []interface{}
}

//type stack StackEmptyInterface

func (s StackEmptyInterface) Size() int {
  return len(s.items)
}

func (s StackEmptyInterface) IsEmpty() bool {
  return s.Size() == 0
}

func (s *StackEmptyInterface) Push(item interface{}) {
  s.items = append(s.items, item)
}

func (s *StackEmptyInterface) Pop() interface{} {
  if s.IsEmpty() {
	return nil
  }

  pos := len(s.items) - 1
  val := s.items[pos]
  s.items = s.items[:pos]

  return val
}

func (s *StackEmptyInterface) Peek() interface{} {
  if s.IsEmpty() {
	return nil
  }

  return s.items[s.Size() - 1]
}

