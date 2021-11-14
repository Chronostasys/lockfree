package lockfree

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type QueueItem struct {
	Item interface{}
	Next *QueueItem
}

type LockFreeQueue struct {
	head *QueueItem
	tail *QueueItem
}

func Make() *LockFreeQueue {
	emp := new(QueueItem)
	return &LockFreeQueue{
		head: emp,
		tail: emp,
	}
}

func (queue *LockFreeQueue) Enqueue(elm interface{}) {
	newItm := &QueueItem{
		Item: elm,
	}
	for {
		p := (*unsafe.Pointer)(unsafe.Pointer(&queue.tail))
		tail := (*QueueItem)(atomic.LoadPointer(p))
		np := (*unsafe.Pointer)(unsafe.Pointer(&tail.Next))
		if atomic.CompareAndSwapPointer(np, nil, unsafe.Pointer(newItm)) {
			atomic.CompareAndSwapPointer(p, unsafe.Pointer(queue.tail), unsafe.Pointer(newItm))
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func (queue *LockFreeQueue) Dequeue() interface{} {
	var v interface{}
	for {
		p := (*unsafe.Pointer)(unsafe.Pointer(&queue.head))
		head := (*QueueItem)(atomic.LoadPointer(p))
		if head.Next == nil {
			return nil
		}
		v = head.Next.Item
		if atomic.CompareAndSwapPointer(p, unsafe.Pointer(head), unsafe.Pointer(head.Next)) {
			break
		}
		time.Sleep(time.Millisecond)
	}
	return v
}
