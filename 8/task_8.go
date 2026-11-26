package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := NewCustomWaitGroup()

	wg.Add(5)

	for i := range 5 {
		go func() {
			defer wg.Done()

			fmt.Println(i)
		}()
	}

	wg.Wait()
}

type CustomWaitGroup struct {
	mu      sync.Mutex
	counter int64
	sem     chan struct{}
}

func NewCustomWaitGroup() *CustomWaitGroup {
	return &CustomWaitGroup{
		sem: make(chan struct{}),
	}
}

func (wg *CustomWaitGroup) Add(delta int64) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	wg.counter += delta

	if wg.counter < 0 {
		panic("CustomWaitGroup: отрицательный счётчик")
	}

	if wg.counter == 0 && delta < 0 {
		close(wg.sem)
		wg.sem = make(chan struct{})
	}
}

func (wg *CustomWaitGroup) Done() {
	wg.Add(-1)
}

func (wg *CustomWaitGroup) Wait() {
	wg.mu.Lock()

	if wg.counter == 0 {
		wg.mu.Unlock()
		return
	}

	ch := wg.sem
	wg.mu.Unlock()
	<-ch
}
