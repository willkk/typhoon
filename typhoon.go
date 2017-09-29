package typhoon

import (
	"net/http"
	"typhoon/core/task"
	"typhoon/core"
)

type Typhoon struct {
	router *core.Router
	Server http.Server

	// services. Every Typhoon instance has its own services.
	services []task.Task
}

// Build new Typhoon instance
func New() *Typhoon {
	return &Typhoon{
		router: core.MainRouter(),
		Server: http.Server{
			Handler: core.MainRouter(),
		},
	}
}

// pattern matches req.URL.Path
func (tp *Typhoon)AddRoute(pattern string, task task.CommandTask) {
	tp.router.AddRoute(pattern, task)
}

// Add normal tasks that will be executed in separate go-routine.
func (tp *Typhoon)AddTask(handler func(ctx *task.Context)) {
	tp.services = append(tp.services, task.TaskFunc(handler))
}

func (tp *Typhoon)StartTasks() {
	for _, st := range tp.services {
		go st.Do(task.NewContext())
	}
}

func (tp *Typhoon)Run(addr string) error {
	tp.Server.Addr = addr

	return tp.Server.ListenAndServe()
}

// manually start a task
func StartTask(taskfunc func(ctx *task.Context)) {
	go task.TaskFunc(taskfunc).Do(task.NewContext())
}

