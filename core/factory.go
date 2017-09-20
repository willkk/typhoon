package core

var factory_ factory = factory{ tasks: make(map[string]Task) }

type factory struct {
	tasks map[string]Task
}

func RegisterTask(name string, task Task) error {
	factory_.tasks[name] = task
	return nil
}

func FactoryTask(name string) Task {
	return factory_.tasks[name].Clone()
}


