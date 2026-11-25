package main

import (
	"fmt"
	"math/rand"
)

func main() {
	done := make(chan struct{})
	gen := RandomGenerator(done)

	for range 10 {
		fmt.Println(<-gen)
	}

	close(done)
}

func RandomGenerator(done <-chan struct{}) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)

		for {
			select {
			case output <- rand.Intn(100):
			case <-done:
				return
			}
		}
	}()

	return output
}
