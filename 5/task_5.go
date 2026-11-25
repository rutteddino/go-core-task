package main

import "fmt"

func main() {
	a := []int{65, 3, 58, 678, 64}
	b := []int{64, 2, 3, 43}

	hasIntersection, common := FindIntersections(a, b)
	fmt.Printf("%v, %#v\n", hasIntersection, common)
}

// FindIntersections ищет общие элементы в слайсах s1 и s2.
// Возвращает true, если найдено хотя бы одно совпадение, и слайс общих элементов.
//
// Каждый общий элемент входит в результат только один раз (дубликаты убираются).
// Порядок элементов в результате соответствует порядку их первого появления
// в более длинном исходном слайсе (или в s1, если длины слайсов равны).
func FindIntersections[T comparable](s1, s2 []T) (bool, []T) {
	result := make([]T, 0)

	// Оптимизация по памяти: в map помещаем элементы меньшего среза.
	// Считаем, что меньший — это s2.
	if len(s2) > len(s1) {
		s1, s2 = s2, s1
	}

	set := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		set[v] = struct{}{}
	}

	for _, v := range s1 {
		if _, ok := set[v]; ok {
			result = append(result, v)
			delete(set, v) // дедупликация: элемент попадает в результат один раз
		}
	}

	return len(result) > 0, result
}
