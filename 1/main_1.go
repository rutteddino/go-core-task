package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

var numDecimal int = 42           // Десятичная система
var numOctal int = 052            // Восьмеричная система
var numHexadecimal int = 0x2A     // Шестнадцатиричная система
var pi float64 = 3.14             // Тип float64
var name string = "Golang"        // Тип string
var isActive bool = true          // Тип bool
var complexNum complex64 = 1 + 2i // Тип complex64

func main() {
	// Определение и вывод типов переменных
	printType(numDecimal)
	fmt.Printf("Значение: %#o(%v), тип: %T\n", numOctal, numOctal, numOctal)
	fmt.Printf("Значение: %#x(%v), тип: %T\n", numHexadecimal, numHexadecimal, numHexadecimal)
	printType(pi)
	printType(name)
	printType(isActive)
	printType(complexNum)

	// Объединение в строку
	var str = joinAsStrings(numDecimal, numOctal, numHexadecimal, pi, name, isActive, complexNum)
	fmt.Println("Объединённая строка:", str)
	hasedStr := hashRunesWithSalt(stringToRunes(str), "salt")
	fmt.Println("Хэш строки с солью:", hasedStr)
}

// printType выводит значения и типы переданных переменных
func printType(val any) {
	fmt.Printf("Значение: %v, тип: %T\n", val, val)
}

// joinAsStrings преобразует все переданные значения в строки и объединяет их
func joinAsStrings(vals ...any) string {
	var sb strings.Builder
	for i, v := range vals {
		sb.WriteString(fmt.Sprint(v))
		if i != len(vals)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

// stringToRunes преобразует строку в срез рун
func stringToRunes(s string) []rune {
	return []rune(s)
}

// hashRunesWithSalt вставляет соль в середину среза рун и возвращает SHA256-хэш в виде hex-строки
func hashRunesWithSalt(runes []rune, salt string) string {
	mid := len(runes) / 2
	saltRunes := []rune(salt)

	salted := make([]rune, 0, len(runes)+len(saltRunes))
	salted = append(salted, runes[:mid]...)
	salted = append(salted, saltRunes...)
	salted = append(salted, runes[mid:]...)

	hash := sha256.Sum256([]byte(string(salted)))
	return fmt.Sprintf("%x", hash)
}
