package entity

import (
	"fmt"
	"reflect"
)

//定义注册结构map
type RegisterStructMaps map[string]reflect.Type

//根据name初始化结构
//在这里根据结构的成员注解进行注入，这里没有实现，只是简单都初始化
func (rsm RegisterStructMaps) New(name string) (c interface{}, err error) {
	if v, ok := rsm[name]; ok {
		c = reflect.New(v).Interface()
	} else {
		err = fmt.Errorf("not found %s struct", name)
	}
	return
}

//根据名字注册实例
func (rsm RegisterStructMaps) Register(name string, c interface{}) {
	rsm[name] = reflect.TypeOf(c).Elem()

}

func GetEntity(name string) interface{} {
	rsm := RegisterStructMaps{}

	//注册接口结构体类
	//test1 := new(Zxcs)

	switch name {
	case "lingji":
		rsm.Register(name, &Lingji{})

	case "zxcs":
		rsm.Register(name, &Zxcs{})

	case "shop":
		rsm.Register(name, &Shop{})

	default:
		rsm.Register(name, &Lingji{})
	}

	test11, _ := rsm.New(name)

	return test11

}
