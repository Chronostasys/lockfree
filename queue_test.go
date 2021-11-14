package lockfree

import (
	"log"
	"runtime"
	"sync"
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	queue := Make()
	wg := sync.WaitGroup{}
	m := map[int]bool{}
	iter := 1000000
	wg.Add(iter)
	for i := 0; i < iter; i++ {
		m[i] = true
		go func(i int) {
			queue.Enqueue(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	wg.Add(iter)
	mu := sync.Mutex{}
	for i := 0; i < iter; i++ {
		go func(i int) {
			el := queue.Dequeue()
			if el == nil {
				log.Fatalf("queue len not right. expected %v, got %v", iter, i)
			}
			mu.Lock()
			delete(m, el.(int))
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := range m {
		t.Fatalf("queue does not contain enqueued value %v. len m %v", i, len(m))
	}
}

func BenchmarkLockfree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench(b)
		runtime.GC()
	}
}
func bench(b *testing.B) {
	b.StartTimer()
	iter := 1000000
	queue := Make()
	doneCh := make(chan struct{})
	dewg := sync.WaitGroup{}
	dewg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer dewg.Done()
			for {
				select {
				case <-doneCh:
					return
				default:
				}
				el := queue.Dequeue()
				if el == 0 {
					close(doneCh)
				}
			}
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(iter * 10)
	for j := 0; j < 10; j++ {
		go func() {
			for i := 0; i < iter; i++ {
				queue.Enqueue(1)
				wg.Done()
			}
		}()
	}
	wg.Wait()
	queue.Enqueue(0)
	<-doneCh
	b.StopTimer()
	dewg.Wait()

}

func BenchmarkLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchLock(b)
		runtime.GC()
	}
}
func benchLock(b *testing.B) {
	b.StartTimer()
	iter := 1000000
	dewg := sync.WaitGroup{}
	dewg.Add(10)
	queue := MakeLK()
	doneCh := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			defer dewg.Done()
			for {
				select {
				case <-doneCh:
					return
				default:
				}
				el := queue.Dequeue()
				if el == 0 {
					close(doneCh)
				}
			}
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(iter * 10)
	for j := 0; j < 10; j++ {
		go func() {
			for i := 0; i < iter; i++ {
				queue.Enqueue(1)
				wg.Done()
			}
		}()
	}
	wg.Wait()
	queue.Enqueue(0)
	<-doneCh
	b.StopTimer()
	dewg.Wait()
}
