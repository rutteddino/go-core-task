package main

import (
	"fmt"
	"maps"
)

// StringIntMap - обёртка над map[string]int со своим API.
type StringIntMap struct {
	data map[string]int
}

// NewStringIntMap создаёт новую пустую map.
func NewStringIntMap(cap int) *StringIntMap {
	validCap := max(0, cap)
	return &StringIntMap{data: make(map[string]int, validCap)}
}

// Add добавляет новую пару ключ-значение в map.
// Если ключ уже существует, значение перезаписывается.
func (m *StringIntMap) Add(key string, value int) {
	m.data[key] = value
}

// Remove удаляет элемент по ключу из map.
// Если ключа нет, операция ничего не делает.
func (m *StringIntMap) Remove(key string) {
	delete(m.data, key)
}

// Copy возвращает новую map[string]int, содержащую все элементы текущей карты.
func (m *StringIntMap) Copy() map[string]int {
	copied := make(map[string]int, len(m.data))
	maps.Copy(copied, m.data)
	return copied
}

// Exists проверяет, существует ли ключ в карте.
func (m *StringIntMap) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

// Get возвращает значение по ключу и булевый флаг,
// указывающий на успешность операции (найден ли ключ).
func (m *StringIntMap) Get(key string) (int, bool) {
	val, ok := m.data[key]
	return val, ok
}

// String возвращает строковое предстовление map для отладки
func (m *StringIntMap) string() string {
	return fmt.Sprintf("%v", m.data)
}

func main() {
	m := NewStringIntMap(0)

	// Добавление элементов
	m.Add("яблоки", 10)
	m.Add("бананы", 5)
	m.Add("апельсины", 3)
	fmt.Println(m.string())

	// Перезапись существующего ключа
	m.Add("бананы", 7)
	if v, ok := m.Get("бананы"); ok {
		fmt.Printf("После перезаписи бананы = %d\n", v)
	}

	// Проверка наличия ключа
	fmt.Printf("Есть ли яблоки? %v\n", m.Exists("яблоки"))
	fmt.Printf("Есть ли груши?  %v\n", m.Exists("груши"))

	// Получение значения
	if v, ok := m.Get("яблоки"); ok {
		fmt.Printf("Значение для яблоки: %d\n", v)
	}
	if _, ok := m.Get("груши"); !ok {
		fmt.Println("Ключ груши не найден")
	}

	// Копирование карты (копия независима от оригинала)
	copied := m.Copy()
	copied["яблоки"] = 999 // меняем копию — оригинал не пострадает
	origApples, _ := m.Get("яблоки")
	fmt.Printf("Оригинал: яблоки = %d, копия: яблоки = %d\n", origApples, copied["яблоки"])

	// Удаление элемента
	m.Remove("апельсины")
	fmt.Printf("После Remove: есть ли апельсины? %v\n", m.Exists("апельсины"))

	// Вывод всех элементов через копию
	fmt.Println("Все элементы карты:")
	for k, v := range m.Copy() {
		fmt.Printf("  %s: %d\n", k, v)
	}
}
