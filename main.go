package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"datacenter.analysis.api/common"
	"datacenter.analysis.api/router"
)

func main() {
	port := common.GetConf("port", "8080")

	router := router.NewRouter()

	log.Fatal(http.ListenAndServe(":"+port, router))
}
