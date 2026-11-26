package main

import (
	"sort"
	"testing"
	"time"
)

// makeChan отправляет значения в канал и закрывает его.
func makeChan[T any](values ...T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

// collect вычитывает канал до закрытия с таймаутом.
func collect[T any](t *testing.T, ch <-chan T) []T {
	t.Helper()
	result := make([]T, 0)
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, v)
		case <-timeout:
			t.Fatal("таймаут: канал не закрылся")
		}
	}
}

func TestMergeAllValues(t *testing.T) {
	ch1 := makeChan(1, 2, 3)
	ch2 := makeChan(4, 5, 6)
	ch3 := makeChan(7, 8, 9)

	got := collect(t, Merge(ch1, ch2, ch3))

	// Порядок не гарантирован — сравниваем после сортировки
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	if len(got) != len(want) {
		t.Fatalf("получено %d значений %v; ожидалось %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Merge = %v; ожидалось %v", got, want)
		}
	}
}

func TestMergeSingleChannel(t *testing.T) {
	got := collect(t, Merge(makeChan("a", "b")))

	sort.Strings(got)
	want := []string{"a", "b"}

	if len(got) != len(want) {
		t.Fatalf("получено %v; ожидалось %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Merge = %v; ожидалось %v", got, want)
		}
	}
}

func TestMergeNoChannels(t *testing.T) {
	// Без входных каналов выходной должен закрыться сразу
	got := collect(t, Merge[int]())

	if len(got) != 0 {
		t.Errorf("Merge() без входных каналов вернул значения: %v", got)
	}
}

func TestMergeEmptyChannels(t *testing.T) {
	// Каналы без значений: выход закрывается, значений нет
	got := collect(t, Merge(makeChan[int](), makeChan[int]()))

	if len(got) != 0 {
		t.Errorf("Merge от пустых каналов вернул значения: %v", got)
	}
}

func TestMergeManyChannels(t *testing.T) {
	const n = 100
	chs := make([]<-chan int, n)
	for i := range chs {
		chs[i] = makeChan(i)
	}

	got := collect(t, Merge(chs...))

	if len(got) != n {
		t.Fatalf("получено %d значений; ожидалось %d", len(got), n)
	}
	seen := make(map[int]bool, n)
	for _, v := range got {
		if seen[v] {
			t.Errorf("значение %d получено дважды", v)
		}
		seen[v] = true
	}
}

func TestMergeOutputClosesOnlyAfterAllInputs(t *testing.T) {
	// Медленный канал: закроется позже остальных
	slow := make(chan int, 1)
	slow <- 42

	fast := makeChan(1, 2, 3)
	out := Merge(slow, fast)

	// Вычитываем быстрые значения, канал ещё не должен быть закрыт
	// (slow ещё не закрыт)
	received := 0
	timeout := time.After(2 * time.Second)
loop:
	for received < 3 {
		select {
		case v, ok := <-out:
			if !ok {
				t.Fatal("выходной канал закрылся до закрытия всех входных")
			}
			if v == 42 {
				continue // значение из slow тоже считается, но нам важны быстрые
			}
			received++
		case <-timeout:
			break loop
		}
	}

	// Теперь закрываем slow — выход должен закрыться после вычитки
	close(slow)
	got := collect(t, out)
	t.Logf("дочитано после закрытия slow: %v", got)
}
