package core

import (
	"net/http"
	"fmt"
	"reflect"
)

type taskType int

const (
	taskTypeCommand taskType = iota  // web api tasks
	taskTypeService         		 // local long-term live routines.
)

// Task is the interface that every "task" should implement.
type Task interface {
	// Do executes task.
	Do()(resp interface{})
	// Clone clones a copy of self
	Clone() Task
}

// commandTask additionally does "Prepare"/"Response" before/after Do function.
type commandTask interface {
	Task

	// Prepare does the preparation before calling Do.
	Prepare(w http.ResponseWriter, r *http.Request) error
	// Response replies result to client.
	Response(w http.ResponseWriter, resp []byte) error
}

func NewHandler(task Task, tasktype taskType) http.Handler {
	return &taskHandler{task, tasktype}
}

type taskHandler struct {
	task Task
	taskType taskType
}

func (th *taskHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	task := th.task.Clone()
	if task == nil {
		panic(fmt.Sprintf("(%v).Clone failed.", reflect.TypeOf(task)))
	}

	cmdtask, ok := task.(commandTask)
	if !ok {
		panic(fmt.Sprintf("(%v) type assertion failed.", reflect.TypeOf(task)))
	}

	err := cmdtask.Prepare(w, r)
	if err != nil {
		cmdtask.Response(w, []byte(err.Error()))
		return
	}

	resp := task.Do()
	if resp != nil {
		cmdtask.Response(w, resp.([]byte))
	}

}
