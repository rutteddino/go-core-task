package main

import "fmt"

func main() {
	in := make(chan uint8)
	result := CubePipeline(in)

	go func() {
		defer close(in)
		for v := range uint8(10) {
			in <- v
		}
	}()

	for v := range result {
		fmt.Printf("Результат: %v\n", v)
	}

	fmt.Println("Конвейер завершен")
}

func CubePipeline(in <-chan uint8) <-chan float64 {
	out := make(chan float64)

	go func() {
		defer close(out)

		for v := range in {
			f := float64(v)
			out <- f * f * f
		}
	}()

	return out
}
