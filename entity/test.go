package entity

import (
	"fmt"
)

type Test interface {
	M1()
	M2()
}

func M1() {
	fmt.Println("test")
}
