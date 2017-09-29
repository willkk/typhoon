package task

import (
	"sync/atomic"
	"net/http"
	"context"
)

var taskId uint64 = 0

type Context struct {
	// W&R is only used in commandTask
	W http.ResponseWriter
	R *http.Request

	// TaskId is used in both of serviceTask and commandTask.
	Id uint64

	// user context data
	UserContext context.Context
}

// NewContext returns a new Context that will be used throughout processing-cycle.
func NewContext(w http.ResponseWriter, r *http.Request)*Context {
	return &Context{
		W: w,
		R: r,
		Id: atomic.AddUint64(&taskId, 1),
	}
}