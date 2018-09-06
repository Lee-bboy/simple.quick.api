package entity

import (
	"fmt"
)

type Shop struct {
	value string
}

func (shop *Shop) SetValue(value string) {
	shop.value = value
}

func (shop *Shop) Print() {
	fmt.Println("sadfadf")
}
