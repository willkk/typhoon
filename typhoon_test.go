package typhoon

import (
	"testing"
	"typhoon/core"
	"fmt"
)

type ServiceTask struct {

}

func (st *ServiceTask)Do()(interface{}, error) {
	fmt.Println("ServiceTask.Do()")
	return nil, nil
}

func (st *ServiceTask)Clone() core.Task {
	return nil
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()
	tp.AddServiceRoute("timer/print", &ServiceTask{})

	tp.Run(":8086")
}

