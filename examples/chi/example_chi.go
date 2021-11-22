package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/igknot/apmzap"
	"go.elastic.co/apm/module/apmchi"
	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
)

func main() {
	r := chi.NewRouter()
	r.Use(apmchi.Middleware())
	log, _ := zap.NewProduction(zap.WrapCore((&apmzap.Core{}).WrapCore))
	//figure out how to change time stamp format.
	r.Use(middleware.NewZapLogger("example-chi", log))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	port := "3666"
	log.Info("ListenAndServe :" + port )
	http.ListenAndServe(":"+port , r)
}
