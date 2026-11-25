package main

import (
	"testing"
)

// equalAsMultiset сравнивает два среза без учёта порядка элементов,
// но с учётом количества вхождений (дубликаты значимы).
func equalAsMultiset[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[T]int, len(a))
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}

// testFindIntersections — обобщённый помощник для табличных тестов
// FindIntersections с произвольным comparable-типом.
// Результат сравнивается без привязки к порядку элементов.
func testFindIntersections[T comparable](t *testing.T, tests []struct {
	name    string
	s1      []T
	s2      []T
	wantHas bool
	want    []T
}) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has, got := FindIntersections(tt.s1, tt.s2)
			if has != tt.wantHas {
				t.Errorf("FindIntersections(%v, %v): флаг = %v; ожидалось %v", tt.s1, tt.s2, has, tt.wantHas)
			}
			if got == nil {
				t.Errorf("FindIntersections(%v, %v) вернул nil вместо пустого среза", tt.s1, tt.s2)
			}
			if !equalAsMultiset(got, tt.want) {
				t.Errorf("FindIntersections(%v, %v) = %v; ожидалось (без учёта порядка) %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestFindIntersectionsInt(t *testing.T) {
	tests := []struct {
		name    string
		s1      []int
		s2      []int
		wantHas bool
		want    []int
	}{
		{
			name:    "пример из задания",
			s1:      []int{65, 3, 58, 678, 64},
			s2:      []int{64, 2, 3, 43},
			wantHas: true,
			want:    []int{3, 64},
		},
		{
			name:    "первый срез меньше второго",
			s1:      []int{3, 64},
			s2:      []int{65, 3, 58, 678, 64},
			wantHas: true,
			want:    []int{64, 3},
		},
		{
			name:    "нет пересечений",
			s1:      []int{1, 2},
			s2:      []int{3, 4},
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "полное совпадение",
			s1:      []int{1, 2},
			s2:      []int{1, 2},
			wantHas: true,
			want:    []int{1, 2},
		},
		{
			name:    "один общий элемент",
			s1:      []int{42},
			s2:      []int{1, 42, 100},
			wantHas: true,
			want:    []int{42},
		},
		{
			name:    "пустой первый срез",
			s1:      []int{},
			s2:      []int{1, 2},
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "пустой второй срез",
			s1:      []int{1, 2},
			s2:      []int{},
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "оба среза пустые",
			s1:      []int{},
			s2:      []int{},
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "nil-срезы",
			s1:      nil,
			s2:      nil,
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "nil и непустой срез",
			s1:      nil,
			s2:      []int{1},
			wantHas: false,
			want:    []int{},
		},
		{
			name:    "дубликаты убираются из результата",
			s1:      []int{3},
			s2:      []int{3, 3, 3},
			wantHas: true,
			want:    []int{3}, // дедупликация: одна тройка вне зависимости от кратности
		},
		{
			name:    "дубликаты в большем срезе схлопываются",
			s1:      []int{3, 3},
			s2:      []int{3, 3, 3},
			wantHas: true,
			want:    []int{3},
		},
		{
			name:    "несколько общих элементов с дубликатами",
			s1:      []int{1, 2, 1, 2, 3},
			s2:      []int{2, 2, 1, 1},
			wantHas: true,
			want:    []int{1, 2},
		},
		{
			name:    "отрицательные числа и ноль",
			s1:      []int{-5, 0, 7},
			s2:      []int{0, -5},
			wantHas: true,
			want:    []int{-5, 0},
		},
	}

	testFindIntersections(t, tests)
}

func TestFindIntersectionsString(t *testing.T) {
	tests := []struct {
		name    string
		s1      []string
		s2      []string
		wantHas bool
		want    []string
	}{
		{
			name:    "базовый случай со строками",
			s1:      []string{"apple", "banana", "cherry"},
			s2:      []string{"cherry", "date"},
			wantHas: true,
			want:    []string{"cherry"},
		},
		{
			name:    "несколько общих строк",
			s1:      []string{"apple", "banana", "cherry"},
			s2:      []string{"cherry", "apple"},
			wantHas: true,
			want:    []string{"apple", "cherry"},
		},
		{
			name:    "нет пересечений",
			s1:      []string{"a"},
			s2:      []string{"b"},
			wantHas: false,
			want:    []string{},
		},
		{
			name:    "чувствительность к регистру",
			s1:      []string{"Apple"},
			s2:      []string{"apple"},
			wantHas: false,
			want:    []string{},
		},
	}

	testFindIntersections(t, tests)
}

// point — пользовательский comparable-тип для проверки работы со структурами.
type point struct {
	X, Y int
}

func TestFindIntersectionsStruct(t *testing.T) {
	tests := []struct {
		name    string
		s1      []point
		s2      []point
		wantHas bool
		want    []point
	}{
		{
			name:    "пересечение структур по значению",
			s1:      []point{{1, 2}, {3, 4}},
			s2:      []point{{3, 4}, {5, 6}},
			wantHas: true,
			want:    []point{{3, 4}},
		},
		{
			name:    "разные структуры не пересекаются",
			s1:      []point{{1, 2}},
			s2:      []point{{2, 1}},
			wantHas: false,
			want:    []point{},
		},
	}

	testFindIntersections(t, tests)
}

// TestFindIntersectionsCommutative проверяет, что результат не зависит
// от порядка аргументов (как мультимножество): функция симметрична.
func TestFindIntersectionsCommutative(t *testing.T) {
	tests := []struct {
		name string
		s1   []int
		s2   []int
	}{
		{"разная длина", []int{2, 1, 3}, []int{3, 2}},
		{"одинаковая длина", []int{1, 2, 3}, []int{3, 1, 2}},
		{"без пересечений", []int{1}, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has1, got1 := FindIntersections(tt.s1, tt.s2)
			has2, got2 := FindIntersections(tt.s2, tt.s1)

			if has1 != has2 {
				t.Errorf("флаги различаются: (%v, %v) -> %v, (%v, %v) -> %v", tt.s1, tt.s2, has1, tt.s2, tt.s1, has2)
			}
			if !equalAsMultiset(got1, got2) {
				t.Errorf("результаты различаются: (%v, %v) -> %v, (%v, %v) -> %v", tt.s1, tt.s2, got1, tt.s2, tt.s1, got2)
			}
		})
	}
}
