package a

import (
	"fmt"
)

func InvalidSinglePointerTestFunc(p *int) { // want "Function InvalidSinglePointerTestFunc has pointer arguments but does not check for nil"

	fmt.Println("Hello, world!")
}

func InvalidMultiplePointerTestFunc(p *int, q *int) { // want "Function InvalidMultiplePointerTestFunc has pointer arguments but does not check for nil"
	if p == nil {
		fmt.Println("Pointer is nil")
	}
	// No check for q
}
