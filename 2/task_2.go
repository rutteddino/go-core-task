package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// 1. Слайс из 10 случайных чисел (при каждом запуске - новые значения)
	originalSlice := generateSlice(10, 0, 100)
	fmt.Println("Сгенерированый слайс: ", originalSlice)

	// 2. Только чётные числа
	fmt.Println("\nЧётные числа (случайный): ", sliceExample(originalSlice))

	// 3. Добавление элемента
	extended := addElements(originalSlice, 999)
	fmt.Println("\nПосле addElements(999): ", extended)
	fmt.Println("Оригинал не изменился: ", originalSlice)

	// 4. Копирование слайса
	copied := copySlice(originalSlice)
	originalSlice[0] = -1 // меняем оригинал
	fmt.Println("\nКопия слайса: ", copied)
	fmt.Println("Оригинал после изменения: ", originalSlice)
	fmt.Println("Копия не изменилась. Изменения оригинала не влияют на копию")

	// 5. Удаление элемента по индексу
	fmt.Println("\nУдаление элемента с индексом 2:", removeElement(originalSlice, 2))
	fmt.Println("Удаление первого элемента: ", removeElement(originalSlice, 0))
	fmt.Println("Удаление последнего элемента: ", removeElement(originalSlice, len(originalSlice)-1))
	fmt.Println("Некорректный индекс 100: ", removeElement(originalSlice, 100))
}

// generateSlice возвращает слайс из n случайных чисел в диапазоне [min, max).
func generateSlice(n, min, max int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rand.Intn(max-min) + min
	}
	return s
}

// sliceExample возвращает новый слайс, содержащий только чётные числа из исходного.
func sliceExample(s []int) []int {
	result := make([]int, 0)
	for _, v := range s {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return result
}

// addElements добавляет число в конец слайса и возвращает новый слайс.
func addElements(s []int, num int) []int {
	result := make([]int, len(s), len(s)+1)
	copy(result, s)
	return append(result, num)
}

// copySlice возвращает независимую копию слайса.
func copySlice(s []int) []int {
	result := make([]int, len(s))
	copy(result, s)
	return result
}

// removeElement возвращает новый слайс без элемента по указанному индексу.
// При некорректном индексе возвращает копию исходного слайса без изменений.
func removeElement(s []int, i int) []int {
	if i < 0 || i > len(s) {
		return copySlice(s)
	}

	result := make([]int, 0, len(s)-1)
	result = append(result, s[:i]...)
	result = append(result, s[i+1:]...)
	return result
}
