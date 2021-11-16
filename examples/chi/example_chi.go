package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/igknot/apmzap/logger"
	"go.elastic.co/apm/module/apmzap"
	"go.elastic.co/apm/module/apmchi"
	"go.uber.org/zap"
)

func main() {
	r := chi.NewRouter()
	r.Use(apmchi.Middleware())
	log := zap.NewExample(zap.WrapCore((&apmzap.Core{}).WrapCore))
	r.Use(logger.NewZapLogger("my-zap-logger",log))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.ListenAndServe(":3666", r)
}
