package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"go.elastic.co/apm/module/apmgorilla"
	 
	"github.com/igknot/apmzap/logger"
	"go.elastic.co/apm/module/apmzap"
	 
	"go.uber.org/zap"
)

func main() {
	router := mux.NewRouter()
	apmgorilla.Instrument(router)
	log := zap.NewExample(zap.WrapCore((&apmzap.Core{}).WrapCore))
	router.Use(logger.NewZapLogger("my-zap-logger",log))
	
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.ListenAndServe(":3666", router)
}
