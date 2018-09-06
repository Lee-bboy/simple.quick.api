package entity

import (
	"fmt"
)

type Lingji struct {
	value string
}

func (lingji *Lingji) SetValue(value string) {
	lingji.value = value
}

func (lingji *Lingji) Print() {
	fmt.Println("sadfadf")
}

func (lingji *Lingji) Test() string {
	return "22"
}
