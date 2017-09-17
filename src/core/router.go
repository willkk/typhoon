package core

import "net/http"

var router *Router = Router{ routes:make(map[string]http.Handler)}

// Router acts the same way as the http.ServerMux
type Router struct {
	routes map[string]http.Handler
}

func (r *Router)AddRoute(path string, task Task) {
	r.routes[path] = task
}

func (r *Router)route(path string) *task {
	return r.routes[path]
}

func (r *Router)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	task := r.route(path)

	handler := NewHandler(task)
	if handler != nil {
		handler.ServeHTTP(w, req)
	}
}

func Router()*Router {
	return router
}


