package typhoon

import (
	"net/http"
	"typhoon/core/task"
	"typhoon/core"
)

type Typhoon struct {
	router *core.Router

	server http.Server

	// services. Every Typhoon instance has its own services.
	services []task.Task
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
func (tp *Typhoon)AddRoute(pattern string, task task.CommandTask) {
	tp.router.AddRoute(pattern, task)
}

// Add normal tasks that will be executed in separate go-routine.
func (tp *Typhoon)AddTask(task task.Task) {
	tp.services = append(tp.services, task)
}

func (tp *Typhoon)StartTasks() {
	for _, st := range tp.services {
		ctx := task.NewContext(nil, nil)
		go st.Do(ctx)
	}
}

func (tp *Typhoon)Run(addr string) error {
	tp.server.Addr = addr

	return tp.server.ListenAndServe()
}



