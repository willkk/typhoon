package core

import (
	"net/http"
	. "typhoon/core/task"
)

// router is used globally in this framework.
var router *Router = &Router{ routes:make(map[string]http.Handler)}

// Router acts the same way as the http.ServerMux
type Router struct {
	routes map[string]http.Handler
}

func (r *Router)AddRoute(path string, task CommandTask) {
	handler := NewHandler(task)
	r.routes[path] = handler
}

func (r *Router)route(path string) http.Handler {
	return r.routes[path]
}

func (r *Router)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	handler := r.route(path)
	if handler != nil {
		handler.ServeHTTP(w, req)
	}
}

func MainRouter()*Router {
	return router
}


