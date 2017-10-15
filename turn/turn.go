package turn

import (
	"container/heap"

	"github.com/karlek/reason/action"
	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/state"

	"github.com/karlek/reason/area"
)

var turnQueue *PriorityQueue

// An Turn is something we manage in a priority queue.
type Turn struct {
	c        *creature.Creature // The c of the turn; arbitrary.
	priority int                // The priority of the turn in the queue.

	// The index is needed by update and is maintained by the heap.Interface
	// methods.
	index int // The index of the turn in the heap.
}

func deferWrap(sta *state.State) {
	state.Stack.Push(*sta)
}

func Proccess(sav *save.Save, a *area.Area) {
	sta := new(state.State)
	*sta = state.Wilderness

	defer deferWrap(sta)
	// Pop the next turn.
	t := heap.Pop(turnQueue).(*Turn)

	// Remove dead creatures from the queue.
	if t.c.Hp <= 0 {
		return
	}

	var timeTaken int
	if t.c.IsHero() {
		timeTaken, *sta = action.HeroTurn(sav, a)
	} else {
		timeTaken = action.Action(t.c, a)
	}
	// If no action was taken, re-insert the turn with the same priority so it
	// will be popped again.
	if timeTaken == 0 {
		heap.Push(turnQueue, t)
		return
	}
	// If some action was taken, reg unit.
	t.c.Reg()

	for k, _ := range *turnQueue {
		(*turnQueue)[k].priority -= t.priority
	}

	// Update new priority.
	t.priority = timeTaken

	// Re-add to queue.
	heap.Push(turnQueue, t)
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
