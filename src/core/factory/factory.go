package factory

import (
	. "core"
)

var factory_ factory = factory{ tasks: make(map[string]*Task) }

// Factory registers and
type Factory interface {
	RegisterFactoryTask(name string, task *Task) error
	FactoryTask(name string) *Task
}

type factory struct {
	tasks map[string]*Task
}

func RegisterTask(name string, task * Task) error {
	factory_.tasks[name] = task
	return nil
}

func FactoryTask(name string) *Task {
	return factory_.tasks[name].Clone()
}


