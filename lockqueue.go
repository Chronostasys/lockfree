package lockfree

import (
	"sync"
)

type LockQueue struct {
	head *QueueItem
	tail *QueueItem
	mu   *sync.Mutex
}

func MakeLK() *LockQueue {
	return &LockQueue{
		mu: &sync.Mutex{},
	}
}

func (queue *LockQueue) Enqueue(elm interface{}) {
	newItm := &QueueItem{
		Item: elm,
	}
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.head == nil {
		queue.head = newItm
		queue.tail = newItm
	} else {
		queue.tail.Next = newItm
		queue.tail = newItm
	}
}

func (queue *LockQueue) Dequeue() interface{} {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.head != nil {
		v := queue.head.Item
		queue.head = queue.head.Next
		if queue.head == nil {
			queue.tail = nil
		}
		return v
	} else {
		return nil
	}
}
