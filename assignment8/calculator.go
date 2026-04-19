package main

func Add(a, b int) int {
	return a + b
}

func Divide(a, b int) int {
	return a / b
}

func Subtract(a, b int) int {
	return a - b
}

// Substract is kept as a backward-compatible alias.
func Substract(a, b int) int {
	return Subtract(a, b)
}
