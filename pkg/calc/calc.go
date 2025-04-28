package calc

import "fmt"

func Add(a, b float64) float64 {
	fmt.Println("This is a test for a demo. Should run only add test.")
	return a + b
}

func Subtract(a, b float64) float64 {
	fmt.Println("Subtract", a, b)
	return a - b
}

func Multiply(a, b float64) float64 {
	if a > b {
		fmt.Println("A is greater than b")
	} else {
		fmt.Println("b is greater than a")
	}
	fmt.Println("test me")
	return a * b
}

func Divide(a, b float64) float64 {
	fmt.Println("Divide", a, b)
	return a / b
}
