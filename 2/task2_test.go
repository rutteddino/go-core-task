package main

import (
	"reflect"
	"testing"
)

func TestSliceExample(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"смешанные числа", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{2, 4, 6, 8, 10}},
		{"все чётные", []int{2, 4, 6}, []int{2, 4, 6}},
		{"все нечётные", []int{1, 3, 5}, []int{}},
		{"пустой слайс", []int{}, []int{}},
		{"отрицательные и ноль", []int{-4, -3, 0, 1}, []int{-4, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceExample(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceExample(%v) = %v, ожидалось %v", tt.input, got, tt.want)
			}
		})
	}
}

// Оригинальный слайс не должен изменяться.
func TestSliceExampleDoesNotModifyOriginal(t *testing.T) {
	input := []int{1, 2, 3, 4}
	snapshot := copySlice(input)
	sliceExample(input)
	if !reflect.DeepEqual(input, snapshot) {
		t.Errorf("оригинальный слайс изменился: %v", input)
	}
}

func TestAddElements(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		num   int
		want  []int
	}{
		{"добавление в конец", []int{1, 2, 3}, 4, []int{1, 2, 3, 4}},
		{"в пустой слайс", []int{}, 7, []int{7}},
		{"отрицательное число", []int{5}, -5, []int{5, -5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addElements(tt.input, tt.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addElements(%v, %d) = %v, ожидалось %v", tt.input, tt.num, got, tt.want)
			}
		})
	}
}

// Исходный слайс не должен измениться после добавления.
func TestAddElementsDoesNotModifyOriginal(t *testing.T) {
	input := []int{1, 2, 3}
	got := addElements(input, 4)
	input[0] = 100
	if got[0] != 1 {
		t.Errorf("изменение оригинала повлияло на результат: %v", got)
	}
}

func TestCopySlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := copySlice(input)
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("copySlice(%v) = %v, слайсы должны совпадать", input, got)
	}
	// Изменение оригинала не должно влиять на копию
	input[0] = 999
	if got[0] == 999 {
		t.Errorf("копия изменилась вместе с оригиналом: %v", got)
	}
	// И наоборот: изменение копии не должно влиять на оригинал
	got[1] = -1
	if input[1] == -1 {
		t.Errorf("оригинал изменился вместе с копией: %v", input)
	}
}

func TestCopySliceEmpty(t *testing.T) {
	if got := copySlice([]int{}); len(got) != 0 {
		t.Errorf("копия пустого слайса должна быть пустой, получено %v", got)
	}
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
		want  []int
	}{
		{"из середины", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}},
		{"первый элемент", []int{1, 2, 3}, 0, []int{2, 3}},
		{"последний элемент", []int{1, 2, 3}, 2, []int{1, 2}},
		{"единственный элемент", []int{42}, 0, []int{}},
		{"индекс больше длины", []int{1, 2, 3}, 100, []int{1, 2, 3}},
		{"отрицательный индекс", []int{1, 2, 3}, -1, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeElement(tt.input, tt.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeElement(%v, %d) = %v, ожидалось %v", tt.input, tt.index, got, tt.want)
			}
		})
	}
}

// Оригинальный слайс не должен изменяться при удалении.
func TestRemoveElementDoesNotModifyOriginal(t *testing.T) {
	input := []int{1, 2, 3, 4}
	snapshot := copySlice(input)
	removeElement(input, 1)
	if !reflect.DeepEqual(input, snapshot) {
		t.Errorf("оригинальный слайс изменился: %v", input)
	}
}
