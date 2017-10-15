package turn

import (
	"container/heap"

	"github.com/karlek/reason/creature"
)

// A PriorityQueue implements heap.Interface and holds Turns.
type PriorityQueue []*Turn

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use lesser
	// than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	t := x.(*Turn)
	t.index = n
	*pq = append(*pq, t)
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	turn := (*pq)[n-1]
	turn.index = -1 // for safety
	*pq = (*pq)[:n-1]
	return turn
}

// update modifies the priority and c of an Turn in the queue.
func (pq *PriorityQueue) update(turn *Turn, c *creature.Creature, priority int) {
	heap.Remove(pq, turn.index)
	turn.c = c
	turn.priority = priority
	heap.Push(pq, turn)
}
