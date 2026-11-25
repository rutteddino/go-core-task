package main

import (
	"reflect"
	"testing"
)

// testGetSlicesDiff — обобщённый помощник для табличных тестов GetSlicesDiff
// с произвольным comparable-типом.
func testGetSlicesDiff[T comparable](t *testing.T, tests []struct {
	name string
	s1   []T
	s2   []T
	want []T
}) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSlicesDiff(tt.s1, tt.s2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlicesDiff(%v, %v) = %v; ожидалось %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestGetSlicesDiffString(t *testing.T) {
	tests := []struct {
		name string
		s1   []string
		s2   []string
		want []string
	}{
		{
			name: "базовый случай",
			s1:   []string{"apple", "banana", "cherry", "date"},
			s2:   []string{"banana", "date"},
			want: []string{"apple", "cherry"},
		},
		{
			name: "все элементы совпадают",
			s1:   []string{"a", "b", "c"},
			s2:   []string{"a", "b", "c"},
			want: []string{},
		},
		{
			name: "нет совпадений",
			s1:   []string{"a", "b"},
			s2:   []string{"c", "d"},
			want: []string{"a", "b"},
		},
		{
			name: "пустой первый срез",
			s1:   []string{},
			s2:   []string{"a", "b"},
			want: []string{},
		},
		{
			name: "пустой второй срез",
			s1:   []string{"a", "b"},
			s2:   []string{},
			want: []string{"a", "b"},
		},
		{
			name: "оба среза пустые",
			s1:   []string{},
			s2:   []string{},
			want: []string{},
		},
		{
			name: "nil вместо первого среза",
			s1:   nil,
			s2:   []string{"a"},
			want: []string{},
		},
		{
			name: "nil вместо второго среза",
			s1:   []string{"a", "b"},
			s2:   nil,
			want: []string{"a", "b"},
		},
		{
			name: "дубликаты в первом срезе сохраняются",
			s1:   []string{"a", "a", "b", "b", "c"},
			s2:   []string{"b"},
			want: []string{"a", "a", "c"},
		},
		{
			name: "дубликаты во втором срезе",
			s1:   []string{"a", "b", "c"},
			s2:   []string{"b", "b", "b"},
			want: []string{"a", "c"},
		},
		{
			name: "второй срез больше первого",
			s1:   []string{"a"},
			s2:   []string{"a", "b", "c", "d"},
			want: []string{},
		},
		{
			name: "сохранение порядка элементов",
			s1:   []string{"d", "b", "a", "c"},
			s2:   []string{"b"},
			want: []string{"d", "a", "c"},
		},
		{
			name: "чувствительность к регистру",
			s1:   []string{"Apple", "apple"},
			s2:   []string{"apple"},
			want: []string{"Apple"},
		},
	}

	testGetSlicesDiff(t, tests)
}

func TestGetSlicesDiffInt(t *testing.T) {
	tests := []struct {
		name string
		s1   []int
		s2   []int
		want []int
	}{
		{
			name: "базовый случай",
			s1:   []int{1, 2, 3, 4, 5},
			s2:   []int{2, 4},
			want: []int{1, 3, 5},
		},
		{
			name: "отрицательные числа и ноль",
			s1:   []int{-1, 0, 1},
			s2:   []int{0},
			want: []int{-1, 1},
		},
		{
			name: "дубликаты сохраняются",
			s1:   []int{7, 7, 8},
			s2:   []int{8},
			want: []int{7, 7},
		},
		{
			name: "nil-срезы",
			s1:   nil,
			s2:   nil,
			want: []int{},
		},
	}

	testGetSlicesDiff(t, tests)
}

// point — пользовательский comparable-тип для проверки работы со структурами.
type point struct {
	X, Y int
}

func TestGetSlicesDiffStruct(t *testing.T) {
	tests := []struct {
		name string
		s1   []point
		s2   []point
		want []point
	}{
		{
			name: "сравнение структур по значению",
			s1:   []point{{1, 2}, {3, 4}, {5, 6}},
			s2:   []point{{3, 4}},
			want: []point{{1, 2}, {5, 6}},
		},
		{
			name: "нулевые структуры",
			s1:   []point{{}, {1, 1}},
			s2:   []point{{}},
			want: []point{{1, 1}},
		},
	}

	testGetSlicesDiff(t, tests)
}

func TestGetSlicesDiffBool(t *testing.T) {
	tests := []struct {
		name string
		s1   []bool
		s2   []bool
		want []bool
	}{
		{
			name: "булевы значения",
			s1:   []bool{true, false, true},
			s2:   []bool{false},
			want: []bool{true, true},
		},
	}

	testGetSlicesDiff(t, tests)
}
