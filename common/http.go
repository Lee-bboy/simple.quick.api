package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Handler func(http.ResponseWriter, *http.Request) *APPError
type response struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type APPError struct {
	Error   error
	Message string
	Code    string
	Status  int
}

func Response(w http.ResponseWriter, description string, code string, status int) {
	out := &response{status, description, code}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	//go log.Printf("response:\t%s", description)
	w.WriteHeader(status)
	w.Write(b)
}

func ContentType(inner http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		inner.ServeHTTP(w, r)
	})
}

func CorsHeader(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://api"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Method", "POST, OPTIONS, GET, HEAD, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-HTTP-Method-Override,accept-charset,accept-encoding , Content-Type, Accept, Cookie")
		inner.ServeHTTP(w, r)
	})
}

func Auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appKey := r.FormValue("appKey")
		timestamp := r.FormValue("timestamp")
		signature := r.FormValue("signature")

		//获取用户秘钥
		appSecret, err := GetAccessSecret(appKey)
		if err != nil {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		//验证签名是否正确
		sign := md5.Sum([]byte(appKey + "+" + appSecret + "+" + timestamp))
		if hex.EncodeToString(sign[:]) != signature {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		//请求是否超时，请求最大时间为5秒
		t, err := strconv.Atoi(timestamp)
		if err != nil || t+5 <= int(time.Now().Unix()) {
			w.Write([]byte(`{"errcode": 403, "errmsg": "API认证失败"}`))
			return
		}

		inner.ServeHTTP(w, r)
	})
}

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		Response(w, e.Message, e.Code, e.Status)
	}
}

func Config(w http.ResponseWriter, r *http.Request) *APPError {
	b := []byte(`{"description":"this is json"}`)
	w.Write(b)
	return nil
}
