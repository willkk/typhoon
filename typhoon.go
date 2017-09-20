package typhoon

import (
	"net/http"
	"typhoon/core"
)

type Typhoon struct {
	router *core.Router

	server http.Server
}

// Build new Typhoon instance
func New() *Typhoon {
	return &Typhoon{
		router: core.MainRouter(),
		server: http.Server{
			Handler: core.MainRouter(),
		},
	}
}

// pattern matches req.URL.Path
func (tp *Typhoon)AddCommandRoute(pattern string, task core.Task) {
	tp.router.AddCommandRoute(pattern, task)
}

func (tp *Typhoon)AddServiceRoute(pattern string, task core.Task) {
	tp.router.AddServiceRoute(pattern, task)
}

func (tp *Typhoon)Run(addr string) error {
	tp.server.Addr = addr

	return tp.server.ListenAndServe()
}

