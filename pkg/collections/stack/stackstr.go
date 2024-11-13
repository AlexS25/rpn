package stack

type StackString struct {
  items []string
}

//type stack StackString

func (s *StackString) Push(item string) {
  s.items = append(s.items, item)
}

func (s StackString) Size() int {
  return len(s.items)
}

func (s StackString) IsEmpty() bool {
  return s.Size() == 0
}

func (s *StackString) Pop() string {
  if s.IsEmpty() {
	return ""
  }

  pos := len(s.items) - 1
  val := s.items[pos]
  s.items = s.items[:pos]

  return val
}

func (s StackString) Peek() string {
  if s.IsEmpty() {
	return ""
  }

  return s.items[s.Size() - 1]
}

