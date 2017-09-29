package task

import (
	"sync/atomic"
	"net/http"
	"context"
)

var taskId uint64 = 0

type Context struct {
	// TaskId is used in both of serviceTask and commandTask.
	Id uint64

	// user context data
	UserContext context.Context
}

type WebContext struct {
	Context

	// W&R is only used in commandTask
	W http.ResponseWriter
	R *http.Request
}

// NewContext returns a new Context that will be used throughout processing-cycle.
func NewContext() *Context {
	return &Context{ Id: atomic.AddUint64(&taskId, 1)}
}

// NewContext returns a new Context that will be used throughout processing-cycle.
func NewWebContext(w http.ResponseWriter, r *http.Request) *WebContext {
	return &WebContext{ Context{Id: atomic.AddUint64(&taskId, 1)}, w, r}
}
