package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errorf("Add(2,3) = %d want %d", got, want)
	}
}

func TestDividePositiveByPositive(t *testing.T) {
	got := Divide(10, 2)
	want := 5

	if got != want {
		t.Errorf("Divide(10,2) = %d want %d", got, want)
	}
}

func TestDivideNegativeByPositive(t *testing.T) {
	got := Divide(-12, 3)
	want := -4

	if got != want {
		t.Errorf("Divide(-12, 3) = %d expected %d", got, want)
	}
}

func TestDivideNegativeByNegative(t *testing.T) {
	got := Divide(-26, -2)
	want := 13

	if got != want {
		t.Errorf("Divide(-26, -2) = %d expected %d", got, want)
	}
}

func TestDivideByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Divide did not panic on division by zero")
		}
	}()

	Divide(10, 0)
}

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive plus zero", 2, 0, 2},
		{"negative plus positive", -3, 4, 1},
		{"both negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("Add(%d,%d) = %d; expected %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive - larger positive ", 2, 3, -1},
		{"positive - smaller positive ", 12, 3, 9},
		{"zero substraction with positive", 67, 0, 67},
		{"zero substraction with negative", -69, 0, -69},
		{"zero substraction", 420, 0, 420},
		{"negative minus positive", -67, 69, -136},
		{"negative minus negative", -420, -69, -351},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("Subtract(%d,%d) = %d; expected %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
