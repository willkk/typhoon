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

func (tp *Typhoon)AddCommandRoute(path string, task core.Task) {
	tp.router.AddCommandRoute(path, task)
}

func (tp *Typhoon)AddServiceRoute(path string, task core.Task) {
	tp.router.AddServiceRoute(path, task)
}

func (tp *Typhoon)Run(addr string) error {
	tp.server.Addr = addr

	return tp.server.ListenAndServe()
}

