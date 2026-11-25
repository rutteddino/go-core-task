package main

import (
	"runtime"
	"testing"
	"time"
)

func TestRandomGeneratorRange(t *testing.T) {
	done := make(chan struct{})
	gen := RandomGenerator(done)
	defer close(done) // остановка генератора

	const n = 1000
	for i := 0; i < n; i++ {
		v := <-gen
		if v < 0 || v >= 100 {
			t.Fatalf("значение %d вне диапазона [0, 100)", v)
		}
	}
}

func TestRandomGeneratorProducesDistinctValues(t *testing.T) {
	done := make(chan struct{})
	gen := RandomGenerator(done)
	defer close(done)

	// Санити-проверка: за 1000 выборок из 100 возможных значений
	// должно встретиться заметно больше одного различного значения.
	const n = 1000
	seen := make(map[int]struct{})
	for i := 0; i < n; i++ {
		seen[<-gen] = struct{}{}
	}

	if len(seen) < 50 {
		t.Errorf("за %d выборок получено всего %d различных значений; генератор подозрительно неслучаен", n, len(seen))
	}
}

// waitClosed вычитывает канал, пока он не закроется.
// ВАЖНО: после close(done) генератор может успеть отправить ещё одно значение -
// когда select видит готовыми обе ветки (читатель ждёт output, done закрыт),
// выбор между ними случаен. Поэтому одно значение в полёте допустимо,
// и проверять закрытие нужно вычитыванием до ok=false, а не одиночным чтением.
func waitClosed(t *testing.T, gen <-chan int) {
	t.Helper()
	timeout := time.After(2 * time.Second)
	for {
		select {
		case _, ok := <-gen:
			if !ok {
				return // канал закрыт
			}
			// значение в полёте - допустимо, продолжаем вычитывать
		case <-timeout:
			t.Fatal("таймаут: канал output не был закрыт после остановки генератора")
		}
	}
}

func TestRandomGeneratorStopsAndClosesOutput(t *testing.T) {
	done := make(chan struct{})
	gen := RandomGenerator(done)

	<-gen       // читаем хотя бы одно значение
	close(done) // сигналим остановку

	// После остановки генератор выходит из горутины и закрывает output
	// (defer close(output)).
	waitClosed(t, gen)
}

func TestRandomGeneratorNoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	done := make(chan struct{})
	gen := RandomGenerator(done)
	<-gen
	close(done)

	// Ждём завершения горутины (с запасом на планировщик).
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			return // горутина завершилась
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("утечка горутины: было %d, стало %d", before, runtime.NumGoroutine())
}

func TestRandomGeneratorOnDemand(t *testing.T) {
	// Небуферизированный канал: чтение N значений не должно блокироваться
	// и зависать - каждое чтение получает ровно одно свежее значение.
	done := make(chan struct{})
	gen := RandomGenerator(done)
	defer close(done)

	const n = 100
	timeout := time.After(2 * time.Second)
	for i := 0; i < n; i++ {
		select {
		case <-gen:
		case <-timeout:
			t.Fatalf("таймаут на чтении значения №%d: генератор не выдаёт значения по требованию", i)
		}
	}
}

func TestRandomGeneratorStopWithoutReading(t *testing.T) {
	// Генератор можно остановить, ни разу не прочитав из него:
	// close(done) разблокирует select даже без потребителя.
	done := make(chan struct{})
	gen := RandomGenerator(done)

	close(done)

	waitClosed(t, gen)
}

func TestRandomGeneratorMultipleGeneratorsStoppable(t *testing.T) {
	// Бонус классического паттерна: одним close(done) можно остановить
	// сразу несколько генераторов.
	done := make(chan struct{})
	gen1 := RandomGenerator(done)
	gen2 := RandomGenerator(done)

	<-gen1
	<-gen2
	close(done)

	waitClosed(t, gen1)
	waitClosed(t, gen2)
}
