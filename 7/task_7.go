package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		ch1 <- 1
		ch1 <- 2
	}()

	go func() {
		defer close(ch2)
		ch2 <- 10
		ch2 <- 20
	}()

	merged := Merge(ch1, ch2)

	for v := range merged {
		fmt.Println(v)
	}
}

// Merge сливает N входных каналов в один выходной.
// Выходной канал закрывается, когда закрыты все входные каналы.
// Порядок элементов в выходном канале не гарантирован.
func Merge[T any](channels ...<-chan T) <-chan T {
	output := make(chan T)

	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Go(func() {
			for v := range ch {
				output <- v
			}
		})
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
