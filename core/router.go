package core

import (
	"net/http"
	. "typhoon/core/task"
)

// router is used globally in this framework.
var router *Router = &Router{mux: http.NewServeMux()}

// Router acts the same way as the http.ServerMux
type Router struct {
	mux *http.ServeMux
}

func (r *Router)AddRoute(path string, task CommandTask) {
	r.mux.Handle(path, NewHandler(task))
}

func (r *Router)route(req *http.Request) http.Handler {
	handler, _ := r.mux.Handler(req)
	return handler
}

func (r *Router)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := r.route(req)
	if handler != nil {
		handler.ServeHTTP(w, req)
	}
}

func MainRouter()*Router {
	return router
}


