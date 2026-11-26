package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCustomWaitGroupBasic(t *testing.T) {
	wg := NewCustomWaitGroup()

	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	wg.Wait() // должен вернуться после Done
}

func TestCustomWaitGroupWaitBlocksUntilDone(t *testing.T) {
	wg := NewCustomWaitGroup()
	wg.Add(1)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait не должен вернуться, пока счётчик > 0
	select {
	case <-done:
		t.Fatal("Wait вернулся до вызова Done")
	case <-time.After(100 * time.Millisecond):
		// ожидаемо: Wait всё ещё заблокирован
	}

	wg.Done()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("таймаут: Wait не разблокировался после Done")
	}
}

func TestCustomWaitGroupZeroCounter(t *testing.T) {
	wg := NewCustomWaitGroup()

	// Wait на нулевом счётчике должен вернуться сразу
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait на нулевом счётчике заблокировался")
	}
}

func TestCustomWaitGroupConcurrentWorkers(t *testing.T) {
	wg := NewCustomWaitGroup()

	const n = 100
	var counter atomic.Int64

	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()

	if got := counter.Load(); got != n {
		t.Errorf("counter = %d; ожидалось %d — Wait вернулся до завершения всех горутин", got, n)
	}
}

func TestCustomWaitGroupMultipleWaiters(t *testing.T) {
	// Все ждущие горутины должны проснуться при закрытии канала
	wg := NewCustomWaitGroup()
	wg.Add(1)

	const waiters = 10
	awoke := make(chan struct{}, waiters)
	for range waiters {
		go func() {
			wg.Wait()
			awoke <- struct{}{}
		}()
	}

	// Даём ждущим горутинам время заблокироваться
	time.Sleep(100 * time.Millisecond)

	select {
	case <-awoke:
		t.Fatal("Waiter проснулся до Done")
	default:
	}

	wg.Done()

	for i := 0; i < waiters; i++ {
		select {
		case <-awoke:
		case <-time.After(2 * time.Second):
			t.Fatalf("таймаут: проснулись только %d из %d ждущих", i, waiters)
		}
	}
}

func TestCustomWaitGroupReuse(t *testing.T) {
	// После достижения нуля канал пересоздаётся — WaitGroup работоспособна снова.
	// Проверяем несколько циклов подряд.
	wg := NewCustomWaitGroup()

	for round := range 5 {
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
			}()
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("цикл %d: таймаут Wait", round)
		}
	}
}

func TestCustomWaitGroupAddWhileWaiting(t *testing.T) {
	// Сценарий: Wait заблокирован, затем Add(+1) и Done() —
	// счётчик не должен «мимоходом» разбудить Wait дважды.
	wg := NewCustomWaitGroup()
	wg.Add(1)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	time.Sleep(50 * time.Millisecond) // Wait заблокировался

	wg.Add(1) // счётчик: 2
	wg.Done() // счётчик: 1 — Wait НЕ должен проснуться
	wg.Done() // счётчик: 0 — Wait просыпается

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("таймаут: Wait не завершился")
	}
}

func TestCustomWaitGroupNegativeCounterPanics(t *testing.T) {
	wg := NewCustomWaitGroup()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ожидалась паника при отрицательном счётчике")
		}
	}()

	wg.Done() // счётчик станет -1 → паника
}

func TestCustomWaitGroupStress(t *testing.T) {
	wg := NewCustomWaitGroup()

	const n = 1000
	var counter atomic.Int64

	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()

	if got := counter.Load(); got != n {
		t.Errorf("counter = %d; ожидалось %d", got, n)
	}
}
