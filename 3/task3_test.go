package main

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	m := NewStringIntMap(0)

	m.Add("one", 1)
	m.Add("two", 2)

	if got, ok := m.Get("one"); !ok || got != 1 {
		t.Errorf("Get(%q) = %d, %v; ожидалось 1, true", "one", got, ok)
	}
	if got, ok := m.Get("two"); !ok || got != 2 {
		t.Errorf("Get(%q) = %d, %v; ожидалось 2, true", "two", got, ok)
	}
}

func TestAddOverwrite(t *testing.T) {
	m := NewStringIntMap(0)

	m.Add("key", 1)
	m.Add("key", 42)

	got, ok := m.Get("key")
	if !ok || got != 42 {
		t.Errorf("после перезаписи Get(%q) = %d, %v; ожидалось 42, true", "key", got, ok)
	}
}

func TestRemove(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("key", 10)

	m.Remove("key")

	if m.Exists("key") {
		t.Error("после Remove ключ всё ещё существует")
	}
	if _, ok := m.Get("key"); ok {
		t.Error("после Remove Get вернул ok=true")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("a", 1)

	// Удаление несуществующего ключа не должно вызывать панику
	// и не должно затрагивать другие элементы.
	m.Remove("missing")

	if !m.Exists("a") {
		t.Error("удаление несуществующего ключа затронуло существующий")
	}
}

func TestCopy(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("a", 1)
	m.Add("b", 2)

	copied := m.Copy()

	want := map[string]int{"a": 1, "b": 2}
	if !reflect.DeepEqual(copied, want) {
		t.Errorf("Copy() = %v; ожидалось %v", copied, want)
	}
}

func TestCopyIsIndependent(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("a", 1)

	copied := m.Copy()
	copied["a"] = 999
	copied["new"] = 100

	if got, _ := m.Get("a"); got != 1 {
		t.Errorf("изменение копии затронуло оригинал: Get(%q) = %d; ожидалось 1", "a", got)
	}
	if m.Exists("new") {
		t.Error("добавление в копию затронуло оригинал")
	}
}

func TestExists(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("present", 5)

	if !m.Exists("present") {
		t.Errorf("Exists(%q) = false; ожидалось true", "present")
	}
	if m.Exists("absent") {
		t.Errorf("Exists(%q) = true; ожидалось false", "absent")
	}
}

func TestGet(t *testing.T) {
	m := NewStringIntMap(0)
	m.Add("x", 7)

	t.Run("существующий ключ", func(t *testing.T) {
		got, ok := m.Get("x")
		if !ok || got != 7 {
			t.Errorf("Get(%q) = %d, %v; ожидалось 7, true", "x", got, ok)
		}
	})

	t.Run("несуществующий ключ", func(t *testing.T) {
		got, ok := m.Get("y")
		if ok {
			t.Errorf("Get(%q) вернул ok=true для отсутствующего ключа", "y")
		}
		if got != 0 {
			t.Errorf("Get(%q) = %d для отсутствующего ключа; ожидалось нулевое значение 0", "y", got)
		}
	})
}

func TestEmptyMap(t *testing.T) {
	m := NewStringIntMap(0)

	if len(m.Copy()) != 0 {
		t.Error("копия пустой карты должна быть пустой")
	}
	if m.Exists("any") {
		t.Error("в пустой карте не должно быть ключей")
	}
}
