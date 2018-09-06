package router

import (
	//"net/http"

	"datacenter.analysis.api/common"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		//var handler http.Handler
		var handler = common.CorsHeader(route.HandlerFunc)
		handler = common.ContentType(handler, route.ContentType)
		handler = route.HandlerFunc
		handler = common.Logger(handler, route.Name)

		//API认证
		handler = common.Auth(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
