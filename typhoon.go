package typhoon

import (
	"net/http"
	"typhoon/core/task"
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
func (tp *Typhoon)AddCommandRoute(pattern string, task task.Task) {
	tp.router.AddCommandRoute(pattern, task)
}

func (tp *Typhoon)AddServiceRoute(task task.Task) {
	tp.router.AddServiceRoute(task)
}

func (tp *Typhoon)Run(addr string) error {
	tp.server.Addr = addr

	task.StartAllServices()
	return tp.server.ListenAndServe()
}



