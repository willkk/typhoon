package core

import "net/http"

// Task is the interface that every "task" should implement.
type Task interface {
	// Prepare does the preparation before calling Do.
	Prepare(w http.ResponseWriter, r *http.Request) (interface{}, error)
	// Do executes task.
	Do() (interface{}, error)
	// Response replies result to client.
	Response(interface{}) error
	// Clone clones a copy of self
	Clone() *Task
}

func NewHandler(task Task) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		var resp interface{}
		var err error
		defer task.Response(resp)

		resp, err = task.Prepare(w, r)
		if err != nil {
			return
		}

		resp, err = task.Do()
	})
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