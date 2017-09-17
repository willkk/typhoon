package typhoon

import (
	"net/http"
	"core"
)

type Typhoon struct {
	router *core.Router

	server http.Server
}

// Build new Typhoon instance
func New() *Typhoon {
	return &Typhoon{
		router: core.Router(),
		server: http.Server{
			Handler: core.Router(),
		},
	}
}

func (tp *Typhoon)AddRoute(path string, task core.Task) {
	tp.router.AddRoute(path, task)
}

func (tp *Typhoon)Run(addr string) error {
	tp.server.Addr = addr

	return tp.server.ListenAndServe()
}

