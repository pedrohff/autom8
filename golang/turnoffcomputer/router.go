package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

type Router struct {
	ChiRouter chi.Router
}

func (r *Router) DeclareRoutes() *Router {
	r.ChiRouter = chi.NewRouter()
	r.ChiRouter.Use(CheckAuthentication)
	r.ChiRouter.Post("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("shutting down")
		//if err := exec.Command("cmd", "/C", "shutdown", "/s").Run(); err != nil {
		//	fmt.Println("Failed to initiate shutdown:", err)
		//}
		_, _ = w.Write([]byte("shutting down"))
	})
	return r
}

func (r *Router) StartServer() {
	_ = http.ListenAndServe(":3000", r.ChiRouter)
}
