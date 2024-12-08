package testdata

import (
	"fmt"
)

func ErrorTestFunc(p *int, q *int) {
	fmt.Println("Hello, world!")
}

func ValidSinglePointerTestFunc(p *int) {
	if p == nil {
		fmt.Println("Pointer is nil")
	}
}

func ValidMultiplePointerTestFunc(p *int, q *int) {
	if p == nil || q == nil {
		fmt.Println("Pointer is nil")
	}
}
