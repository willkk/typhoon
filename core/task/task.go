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

type TaskFunc func(ctx *Context)

func (tf TaskFunc)Do(ctx *Context) {
	tf(ctx)
}

// CommandTask does "Prepare"/"Response" before/after Do function.
// Clone method returns a new copy of commandTask.
type CommandTask interface {
	// Clone returns a copy of self
	Clone() CommandTask
	
	// Prepare does the preparation before calling Do. It works in application layer.
	Prepare(ctx *WebContext) ([]byte, error)
	
	// Do does real business logic. It works in domain layer.
	Do(ctx *WebContext)(resp CommandTaskResp, err error)
	
	// Finish does finishing works before writing response to client. It works in application layer.
	Finish(ctx *WebContext, resp CommandTaskResp) []byte
	
	// Response replies result to client. It works in application layer.
	Response(ctx *WebContext, resp []byte)
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

	ctx := NewWebContext(w, r)

	resp_pre, err := task.Prepare(ctx)
	if err != nil {
		task.Response(ctx, resp_pre)
		return
	}

	// err is ignored if task is typeof CommandTask
	resp, err := task.Do(ctx)
	if err != nil {
		task.Response(ctx, resp)
		return
	}
	
	rst := task.Finish(ctx, resp)
	task.Response(ctx, rst)
}
