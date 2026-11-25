package main

import (
	"reflect"
	"testing"
)

func TestGetSlicesDiff(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSlicesDiff(tt.s1, tt.s2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlicesDiff(%v, %v) = %v; ожидалось %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}
