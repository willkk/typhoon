package task

import (
	"net/http"
)

// Task is the basic service interface that every "task" should implement.
// We define two kinds of Tasks, that is Task and CommandTask. Task represents the
// normal tasks like one-time or cyclic execution go-routines, and CommandTask represents
// web rpc calling tasks.
type Task interface {
	// Do executes task. err is used for ServiceTask and resp is used for CommandTask.
	Do(ctx *Context)
}

// commandTask does "Prepare"/"Response" before/after Do function.
// Clone method returns a new copy of commandTask.
type CommandTask interface {
	Do(ctx *Context)(resp []byte, err error)

	// Clone clones a copy of self
	Clone() CommandTask
	// Prepare does the preparation before calling Do.
	Prepare(ctx *Context) ([]byte, error)
	// Response replies result to client.
	Response(ctx *Context, resp []byte)
}

func NewHandler(task CommandTask) http.Handler {
	return &taskHandler{task}
}

type taskHandler struct {
	task CommandTask
}

func (th *taskHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmdtask, ok := th.task.(CommandTask)
	if !ok {
		panic("type assertion failed.")
	}

	task := cmdtask.Clone()
	if task == nil {
		panic("Clone returns nil.")
	}

	ctx := NewContext(w, r)

	resp_err, err := cmdtask.Prepare(ctx)
	if err != nil {
		cmdtask.Response(ctx, resp_err)
		return
	}

	// err is ignored if task is typeof CommandTask
	resp, _ := task.Do(ctx)
	cmdtask.Response(ctx, resp)
}
