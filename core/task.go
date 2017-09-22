package core

import (
	"net/http"
)

type taskType int

const (
	taskTypeCommand taskType = iota  // web api tasks
	taskTypeService         		 // local long-term live routines.
)

// Task is the interface that every "task" should implement.
type Task interface {
	// Do executes task. err is used for ServiceTask and resp is used for CommandTask.
	Do()(resp []byte, err error)
}

// commandTask additionally does "Prepare"/"Response" before/after Do function.
type commandTask interface {
	Task

	// Clone clones a copy of self
	Clone() Task
	// Prepare does the preparation before calling Do.
	Prepare(w http.ResponseWriter, r *http.Request) ([]byte, error)
	// Response replies result to client.
	Response(w http.ResponseWriter, resp []byte)
}

func NewHandler(task Task, tasktype taskType) http.Handler {
	return &taskHandler{task, tasktype}
}

type taskHandler struct {
	task Task
	taskType taskType
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

	resp_err, err := cmdtask.Prepare(w, r)
	if err != nil {
		cmdtask.Response(w, resp_err)
		return
	}

	// err is ignored if task is typeof CommandTask
	resp, _ := task.Do()
	cmdtask.Response(w, resp)
}
