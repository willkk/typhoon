package core

import "net/http"

var router *Router = Router{ routes:make(map[string]*http.Handler)}

// Router acts the same way as the http.ServerMux
type Router struct {
	routes map[string]*http.Handler
}

func (r *Router)AddRoute(path string, handler *http.Handler) {
	r.routes[path] = handler
}

func (r *Router)Route(path string) *http.Handler {
	return r.routes[path]
}

func (r *Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func Router()*Router {
	return router
}

func AddRoute(path string, handler *http.Handler) {
	router.AddRoute(path, handler)
}

func Route(path string) *http.Handler {
	return router.Route(path)
}

