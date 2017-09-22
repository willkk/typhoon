package core

import (
	"sync/atomic"
	"context"
	"net/http"
)

var routineId uint64 = 0
var routineContext *RoutineContext = &RoutineContext{}

type RoutineContext struct {
}

func (rc *RoutineContext)Id() uint64 {
	return atomic.AddUint64(&routineId, 1)
}

// NewContext returns a new Context that carries value u.
func (rc *RoutineContext)NewContext(ctx context.Context, req *http.Request)*http.Request {
	new_ctx := context.WithValue(ctx, "routineid", routineContext.Id())
	return req.WithContext(new_ctx)
}

// FromContext returns the User value stored in ctx, if any.
func (rc *RoutineContext)FromContext(ctx context.Context, key interface{}) (uint64, bool) {
	u, ok := ctx.Value(key).(uint64)
	return u, ok
}

func NewContext(ctx context.Context, req *http.Request) *http.Request {
	return routineContext.NewContext(ctx, req)
}

func CtxId(ctx context.Context)(uint64) {
	if ctx == nil {
		return 0
	}
	id, _ := routineContext.FromContext(ctx, "routineid")
	return id
}