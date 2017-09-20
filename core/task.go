package core

import (
	"net/http"
	"fmt"
)

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
	Response(w http.ResponseWriter, resp interface{}) error
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

	resp, err := task.Prepare(w, r)
	if err != nil {
		resp = []byte(err.Error())
		task.Response(w, resp)
		fmt.Printf("Prepare err:%s", err)
		return
	}

	resp, err = task.Do()
	task.Response(w, resp)
}
