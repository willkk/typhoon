package core

import "net/http"

// Task is the interface that every "task" should implement.
type Task interface {
	// Prepare does the preparation before calling Do.
	Prepare() error
	// Do executes task.
	Do() error
	// Response replies result to client.
	Response() error
	// Clone clones a copy of self
	Clone() *Task

	http.Handler
}

func NewTask(w http.ResponseWriter, r *http.Request) *Task {
	return &task{
		w: w,
		r: r,
	}
}

type task struct {
	w http.ResponseWriter
	r *http.Request
}

func(t *task)Prepare() error {
	return nil
}

func(t *task)Do() error {
	return nil
}

func(t *task)Response() error {
	return nil
}

func (t *task)Clone() *Task {
	copy := new(task)
	*copy = *t
	return copy
}

func(t *task)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer t.Response()

	if err := t.Prepare(); err != nil {
		return
	}
	if err := t.Do(); err != nil {
		return
	}
}