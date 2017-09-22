package task

// Manage all service tasks
type serviceTask interface {
	Task
}

var servTasks []serviceTask

func AddServiceTask(task Task) {
	st := task.(serviceTask)
	servTasks = append(servTasks, st)
}

func StartAllServices() {
	for _, st := range servTasks {
		ctx := NewContext(nil, nil)
		go st.Do(ctx)
	}
}

