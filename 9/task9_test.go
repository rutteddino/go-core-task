package main

import (
	"testing"
	"time"
)

// collectFloats вычитывает канал float64 до закрытия с таймаутом.
func collectFloats(t *testing.T, ch <-chan float64) []float64 {
	t.Helper()
	result := make([]float64, 0)
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, v)
		case <-timeout:
			t.Fatal("таймаут: выходной канал не закрылся")
		}
	}
}

// uint8Chan отправляет значения в канал и закрывает его.
func uint8Chan(values ...uint8) chan uint8 {
	ch := make(chan uint8)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

func TestCubePipelineBasic(t *testing.T) {
	in := uint8Chan(1, 2, 3, 4, 5)

	got := collectFloats(t, CubePipeline(in))
	want := []float64{1, 8, 27, 64, 125}

	if len(got) != len(want) {
		t.Fatalf("получено %d значений %v; ожидалось %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("значение №%d = %v; ожидалось %v", i, got[i], want[i])
		}
	}
}

func TestCubePipelineZero(t *testing.T) {
	got := collectFloats(t, CubePipeline(uint8Chan(0)))

	if len(got) != 1 || got[0] != 0 {
		t.Errorf("куб 0: получено %v; ожидалось [0]", got)
	}
}

func TestCubePipelineMaxUint8(t *testing.T) {
	// 255^3 = 16581375 — помещается в float64 без потери точности
	got := collectFloats(t, CubePipeline(uint8Chan(255)))

	if len(got) != 1 || got[0] != 16581375 {
		t.Errorf("куб 255: получено %v; ожидалось [16581375]", got)
	}
}

func TestCubePipelineEmptyInput(t *testing.T) {
	// Пустой вход → выход сразу закрывается без значений
	got := collectFloats(t, CubePipeline(uint8Chan()))

	if len(got) != 0 {
		t.Errorf("пустой вход: получено %v; ожидался пустой результат", got)
	}
}

func TestCubePipelinePreservesOrder(t *testing.T) {
	// Конвейер однопоточный: порядок элементов сохраняется
	in := uint8Chan(9, 1, 7, 3, 3, 12)

	got := collectFloats(t, CubePipeline(in))
	want := []float64{729, 1, 343, 27, 27, 1728}

	if len(got) != len(want) {
		t.Fatalf("получено %v; ожидалось %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("позиция %d = %v; ожидалось %v (порядок нарушен)", i, got[i], want[i])
		}
	}
}

func TestCubePipelineOnDemand(t *testing.T) {
	// Числа обрабатываются ПО МЕРЕ ПОСТУПЛЕНИЯ: пишем во входной канал
	// постепенно и проверяем, что каждое значение появляется на выходе
	// до отправки следующего.
	in := make(chan uint8)
	out := CubePipeline(in)

	values := []uint8{2, 5, 10}
	want := []float64{8, 125, 1000}

	go func() {
		defer close(in)
		for _, v := range values {
			in <- v
		}
	}()

	for i, w := range want {
		select {
		case got := <-out:
			if got != w {
				t.Errorf("значение №%d = %v; ожидалось %v", i, got, w)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("таймаут: значение №%d не появилось на выходе по мере поступления", i)
		}
	}
}
