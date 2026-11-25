package main

import (
	"fmt"
)

func main() {
	slice1 := []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}
	slice2 := []string{"banana", "date", "fig"}

	result := GetSlicesDiff(slice1, slice2)
	fmt.Println(result)
}

// GetSlicesDiff возвращает элементы из первого слайса (s1), которых нет во втором (s2).
func GetSlicesDiff[T comparable](s1, s2 []T) []T {
	result := make([]T, 0)
	set := make(map[T]struct{})

	for _, v := range s2 {
		set[v] = struct{}{}
	}

	for _, v := range s1 {
		if _, ok := set[v]; !ok {
			result = append(result, v)
		}
	}

	return result
}
