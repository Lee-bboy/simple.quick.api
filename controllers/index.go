package controllers

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"datacenter.analysis.api/models"
	"datacenter.analysis.api/common"
	"datacenter.analysis.api/entity"
	//"github.com/gorilla/mux"

	//"runtime"

	//"reflect"
)

func Index(w http.ResponseWriter, r *http.Request) {
	b := []byte(`{"index":"this is json"}`)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.WriteHeader(http.StatusOK)
	//common.ERROR.Printf("This is the No[%v] INFO log .", "测试") //测试写LOG

	//models.TestGet() //测试读取数据库

	w.Write(b)
}

func Test(w http.ResponseWriter, r *http.Request) {
	b := []byte(`{"test":"this is json"}`)

	//	var MULTICORE int = runtime.NumCPU() //number of core

	//	runtime.GOMAXPROCS(MULTICORE) //running in multicore

	//	fmt.Printf("with %d core\n", MULTICORE)

	obj := entity.GetEntity("zxcs")

	switch t := obj.(type) {
	case *entity.Zxcs:
		t.Print()
	case *entity.Shop:
		fmt.Printf("nil 1111", 1111)
	case *entity.Lingji:
		fmt.Printf("nil 2222", 2222)
	case nil:
		fmt.Printf("nil value: nothing to check?\n")
	default:
		fmt.Printf("nil sss: nothing to check?\n")
	}

	w.Write(b)
}

func PostTest(w http.ResponseWriter, r *http.Request) {
	b := []byte(`{"post":"this is json"}`)

	fmt.Printf(r.PostFormValue("id"))

	w.Write(b)

}

func LongjieTest(w http.ResponseWriter, r *http.Request) {
	accessKey := "test123"
	accessSecret, err := common.GetAccessSecret(accessKey)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("%s的秘钥为：%s", accessKey, accessSecret)))
	}
}
