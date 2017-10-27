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
	
	// Prepare does the preparation before calling Do.
	// [Application Layer]
	Prepare(ctx *WebContext) (TaskResponse, error)
	
	// Do does real business logic.
	// [Domain Layer]
	Do(ctx *WebContext)(resp TaskResponse, err error)
	
	// Response replies result to client.
	// [Application Layer]
	Response(ctx *WebContext, resp TaskResponse)

	// Finish does finishing works after writing response to client.
	// Finishing works can be like, sending mails, recalculating bonus points
	// after user's successful payment. Theoretically, it should have nothing to do with
	// business logic in Do. Typically, it calls downstream services.
	// [Application Layer]
	Finish(ctx *WebContext, resp TaskResponse)
}

type TaskResponse interface {
	Response() []byte
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

	resp, err := task.Do(ctx)
	task.Response(ctx, resp)

	if err == nil {
		task.Finish(ctx, resp)
	}
}
