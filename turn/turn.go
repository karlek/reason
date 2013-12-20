package turn

import (
	"container/heap"
	"log"

	"github.com/karlek/reason/action"
	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/save"

	"github.com/karlek/worc/area"
)

var turnQueue *PriorityQueue

// An Turn is something we manage in a priority queue.
type Turn struct {
	c        *creature.Creature // The c of the turn; arbitrary.
	priority int                // The priority of the turn in the queue.

	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the turn in the heap.
}

// A PriorityQueue implements heap.Interface and holds Turns.
type PriorityQueue []*Turn

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
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
	old := *pq
	n := len(old)
	turn := old[n-1]
	turn.index = -1 // for safety
	*pq = old[0 : n-1]
	return turn
}

// update modifies the priority and c of an Turn in the queue.
func (pq *PriorityQueue) update(turn *Turn, c *creature.Creature, priority int) {
	heap.Remove(pq, turn.index)
	turn.c = c
	turn.priority = priority
	heap.Push(pq, turn)
}

func Proccess(sav *save.Save, a *area.Area, hero *creature.Creature) {
	t := heap.Pop(turnQueue).(*Turn)
	if t.c.Hp <= 0 {
		return
	}
	var timeTaken int
	if t.c.Name() == hero.Name() {
		timeTaken = action.HeroTurn(sav, a, hero)
	} else {
		timeTaken = t.c.Action(a, hero)
	}
	if timeTaken == 0 {
		heap.Push(turnQueue, t)
		return
	}

	for k, _ := range *turnQueue {
		(*turnQueue)[k].priority -= t.priority
	}
	t.priority = timeTaken
	heap.Push(turnQueue, t)
	log.Println("asdf")
}

func Init(a *area.Area) {
	// Create a priority queue and put the turns in it.
	turnQueue = &PriorityQueue{}
	heap.Init(turnQueue)

	for _, m := range a.Monsters {
		c, ok := m.(*creature.Creature)
		if !ok {
			continue
		}
		turn := &Turn{
			c:        c,
			priority: c.Speed,
		}
		heap.Push(turnQueue, turn)
	}
}
