package calc

import "fmt"

func Add(a, b float64) float64 {
	fmt.Println("Add", a, b)
	return a + b
}

func Subtract(a, b float64) float64 {
	fmt.Println("Subtract", a, b)
	return a - b
}

func Multiply(a, b float64) float64 {
	fmt.Println("Multiply", a, b)
	return a * b
}

func Divide(a, b float64) float64 {
	fmt.Println("Divide", a, b)
	return a / b
}
