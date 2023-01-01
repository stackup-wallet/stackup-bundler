package utils

type Stack[T any] struct {
	keys []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{nil}
}

func (stack *Stack[T]) Push(key T) {
	stack.keys = append(stack.keys, key)
}

func (stack *Stack[T]) Top() (T, bool) {
	var x T
	if len(stack.keys) > 0 {
		x = stack.keys[len(stack.keys)-1]
		return x, true
	}
	return x, false
}

func (stack *Stack[T]) Pop() (T, bool) {
	var x T
	if len(stack.keys) > 0 {
		x, stack.keys = stack.keys[len(stack.keys)-1], stack.keys[:len(stack.keys)-1]
		return x, true
	}
	return x, false
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.keys) == 0
}
