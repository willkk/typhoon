package core

import "net/http"

type taskType int

const (
	taskTypeCommand taskType = iota  // web api tasks
	taskTypeService         		 // local long-term live routines.
)

// Task is the interface that every "task" should implement.
type Task interface {
	// Do executes task.
	Do()(interface{}, error)
	// Clone clones a copy of self
	Clone() Task
}

// commandTask additionally does "Prepare"/"Response" before/after Do function.
type commandTask interface {
	Task

	// Prepare does the preparation before calling Do.
	Prepare(w http.ResponseWriter, r *http.Request) (interface{}, error)
	// Response replies result to client.
	Response(interface{}) error
}

func NewHandler(task Task, tasktype taskType) http.Handler {
	return &taskHandler{task, tasktype}
}

type taskHandler struct {
	task Task
	taskType taskType
}

func (th *taskHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	task := th.task.(commandTask)
	var resp interface{}
	var err error
	defer task.Response(resp)

	resp, err = task.Prepare(w, r)
	if err != nil {
		return
	}

	resp, err = task.Do()
}
