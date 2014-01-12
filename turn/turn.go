package turn

import (
	"container/heap"

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
