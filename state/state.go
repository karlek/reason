package state

type State uint

const (
	Init = iota
	Intro
	Wilderness
	Inventory
	Look
	Drop
	Open
	Close
)

var Stack = new(stack)

type stack []State

func (ss *stack) Push(s State) {
	*ss = append(*ss, s)
}

func (ss *stack) Pop() (s State) {
	if len(*ss) == 0 {
		panic("state stack is empty")
	}
	s = (*ss)[len(*ss)-1]
	*ss = (*ss)[:len(*ss)-1]
	return s
}
