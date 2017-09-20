package core

import "net/http"

var router *Router = &Router{ routes:make(map[string]http.Handler)}

// Router acts the same way as the http.ServerMux
type Router struct {
	routes map[string]http.Handler
}

func (r *Router)AddCommandRoute(path string, task Task) {
	handler := NewHandler(task, taskTypeCommand)
	r.routes[path] = handler
}

func (r *Router)AddServiceRoute(path string, task Task) {
	handler := NewHandler(task, taskTypeService)
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


