package task

// commandTask additionally does "Prepare"/"Response" before/after Do function.
type commandTask interface {
	Task

	// Clone clones a copy of self
	Clone() Task
	// Prepare does the preparation before calling Do.
	Prepare(ctx *TaskContext) ([]byte, error)
	// Response replies result to client.
	Response(ctx *TaskContext, resp []byte)
}
