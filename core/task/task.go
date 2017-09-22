package task

import (
	"net/http"
)

type taskType int

// Task is the interface that every "task" should implement.
type Task interface {
	// Do executes task. err is used for ServiceTask and resp is used for CommandTask.
	Do(ctx *TaskContext)(resp []byte, err error)
}


func NewHandler(task Task) http.Handler {
	return &taskHandler{task}
}

type taskHandler struct {
	task Task
}

func (th *taskHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmdtask, ok := th.task.(commandTask)
	if !ok {
		panic("type assertion failed.")
	}

	task := cmdtask.Clone()
	if task == nil {
		panic("Clone return nil.")
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
