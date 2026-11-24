package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
)

// captureStdout перехватывает вывод в stdout на время выполнения f.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read captured output: %v", err)
	}
	return buf.String()
}

func TestPrintType(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"int", 42, "Значение: 42, тип: int\n"},
		{"float64", 3.14, "Значение: 3.14, тип: float64\n"},
		{"string", "Golang", "Значение: Golang, тип: string\n"},
		{"bool", true, "Значение: true, тип: bool\n"},
		{"complex64", complex64(1 + 2i), "Значение: (1+2i), тип: complex64\n"},
		{"nil", nil, "Значение: <nil>, тип: <nil>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() { printType(tt.val) })
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJoinAsStrings(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		want string
	}{
		{"пустой ввод", nil, ""},
		{"один элемент", []any{42}, "42"},
		{"все переменные", []any{numDecimal, numOctal, numHexadecimal, pi, name, isActive, complexNum},
			"42 42 42 3.14 Golang true (1+2i)"},
		{"разделители-пробелы", []any{"a", "b", "c"}, "a b c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinAsStrings(tt.vals...); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringToRunes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []rune
	}{
		{"пустая строка", "", []rune{}},
		{"ascii", "abc", []rune{'a', 'b', 'c'}},
		// Ключевой кейс: кириллица — рун меньше, чем байт.
		{"кириллица", "при", []rune{'п', 'р', 'и'}},
		{"смешанная", "Go语言", []rune{'G', 'o', '语', '言'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringToRunes(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("rune[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// эталонная реализация вставки соли — независимая от тестируемого кода
func referenceSalted(s, salt string) string {
	r := []rune(s)
	mid := len(r) / 2
	return string(r[:mid]) + salt + string(r[mid:])
}

func TestHashRunesWithSalt(t *testing.T) {
	t.Run("golden: вся строка из main", func(t *testing.T) {
		str := joinAsStrings(numDecimal, numOctal, numHexadecimal, pi, name, isActive, complexNum)
		// SHA-256 от "42 42 42 3.14 Gosaltlang true (1+2i)"
		const want = "30bf9f580918f39f22a18e961f9f49eaebd573243e868cad601a1bd71f371e5c"
		if got := hashRunesWithSalt(stringToRunes(str), "salt"); got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("совпадает с эталонным sha256", func(t *testing.T) {
		cases := []struct{ in, salt string }{
			{"hello", "salt"},    // нечётная длина: mid = 2
			{"hello world", "s"}, // нечётная длина
			{"", "salt"},         // пустой вход
			{"привет", "соль"},   // юникод
			{"abc", ""},          // пустая соль = обычный хэш
		}
		for _, c := range cases {
			salted := referenceSalted(c.in, c.salt)
			want := fmt.Sprintf("%x", sha256.Sum256([]byte(salted)))
			if got := hashRunesWithSalt(stringToRunes(c.in), c.salt); got != want {
				t.Errorf("hashRunesWithSalt(%q, %q) = %s, want %s", c.in, c.salt, got, want)
			}
		}
	})

	t.Run("длина hex-строки", func(t *testing.T) {
		if got := hashRunesWithSalt([]rune("x"), "y"); len(got) != 64 {
			t.Errorf("len = %d, want 64", len(got))
		}
	})

	t.Run("не изменяет исходный срез", func(t *testing.T) {
		in := []rune("abcdef")
		before := string(in)
		_ = hashRunesWithSalt(in, "salt")
		if string(in) != before {
			t.Errorf("входной срез изменён: %q, было %q", string(in), before)
		}
	})
}
