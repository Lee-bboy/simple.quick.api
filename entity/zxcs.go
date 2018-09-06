package entity

import (
	"fmt"
)

type Zxcs struct {
	value string
}

func (zxcs *Zxcs) SetValue(value string) {
	zxcs.value = value
}

func (zxcs *Zxcs) Print() {
	fmt.Println("zxncd")
}

func (zxcs *Zxcs) Test() string {
	return "333"
}
