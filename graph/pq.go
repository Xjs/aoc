package graph

import (
	"container/heap"
)

type id int

type pq struct {
	queue    []id
	priority map[id]float64
	indexes  map[id]int
}

func newPQ() *pq {
	return &pq{
		priority: make(map[id]float64),
		indexes:  make(map[id]int),
	}
}

// Len implements sort.Interface
func (pq *pq) Len() int { return len(pq.queue) }

// Less implements sort.Interface
func (pq *pq) Less(i, j int) bool {
	return pq.priority[pq.queue[i]] < pq.priority[pq.queue[j]]
}

// Swap implements sort.Interface
func (pq *pq) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
	pq.indexes[pq.queue[i]] = i
	pq.indexes[pq.queue[j]] = j
}

// Push implements heap.Interface
func (pq *pq) Push(x interface{}) {
	n := len(pq.queue)
	item := x.(id)
	pq.indexes[item] = n
	pq.queue = append(pq.queue, item)
}

// Pop implements heap.Interface
func (pq *pq) Pop() interface{} {
	old := pq.queue
	n := len(old)
	item := old[n-1]
	delete(pq.indexes, item)
	pq.queue = old[0 : n-1]
	return item
}

// pop returns the point with lowest priority and its priority.
func (pq *pq) pop() (id, float64) {
	p := heap.Pop(pq).(id)
	prio := pq.priority[p]
	delete(pq.priority, p)
	return p, prio
}

// decrease modifies the priority and value of a point in the queue only if the priority is lower than the current priority.
// It returns true if a modification was made.
func (pq *pq) decrease(item id, priority float64) (bool, bool) {
	oldPrio := pq.priority[item]
	if oldPrio < priority {
		return false, false
	}
	pq.update(item, priority)
	return true, oldPrio == priority
}

// update modifies the priority and value of an Item in the queue.
func (pq *pq) update(item id, priority float64) {
	pq.priority[item] = priority
	heap.Fix(pq, pq.indexes[item])
}
