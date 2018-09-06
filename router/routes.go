package router

import (
	"net/http"

	//. "simple.quick.api/common"
	. "simple.quick.api/controllers"
)

const (
	JSON  string = "application/json;charset=utf-8"
	HTML  string = "text/html"
	PLAIN string = "text/plain"
)

type Route struct {
	Name        string //日志中显示的路由名称
	Method      string //GET PUT POST DELETE ...
	Pattern     string //对用的访问路径
	HandlerFunc http.HandlerFunc
	ContentType string //返回的数据类型"application/json;charset=utf-8" 或者 "text/html" 等等 ...
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
		JSON,
	},
	Route{
		"Test",
		"GET",
		"/test/{todoId}",
		Test,
		JSON,
	},
	Route{
		"PostTest",
		"POST",
		"/postTest",
		PostTest,
		JSON,
	},
	Route{
		"LongjieTest",
		"GET",
		"/longjie/test",
		LongjieTest,
		HTML,
	},
	Route{
		"LongjieTest",
		"GET",
		"/longjie/test",
		LongjieTest,
		HTML,
	},
	Route{
		"Associate",
		"POST",
		"/market/associate",
		MDController.Associate,
		JSON,
	},
	Route{
		"Frequency",
		"POST",
		"/market/frequency",
		MDController.Frequency,
		JSON,
	},
	Route{
		"TimeInterval",
		"POST",
		"/market/time_interval",
		MDController.TimeInterval,
		JSON,
	},
	Route{
		"FirstTrade",
		"POST",
		"/market/first_trade",
		MDController.FirstTrade,
		JSON,
	},
	Route{
		"user_portrait_gender",
		"POST",
		"/user/portrait/gender",
		UPController.Gender,
		JSON,
	},
	Route{
		"user_portrait_age",
		"POST",
		"/user/portrait/age",
		UPController.Age,
		JSON,
	},
	Route{
		"user_portrait_hour",
		"POST",
		"/user/portrait/hour",
		UPController.Hour,
		JSON,
	},
}
