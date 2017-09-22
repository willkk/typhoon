package typhoon

import (
	"testing"
	"typhoon/core/task"
	"io/ioutil"
	"errors"
	"encoding/json"
	"time"
	"fmt"
)

type ServiceTask struct {

}

func (st *ServiceTask)Do(ctx *task.TaskContext)([]byte, error) {
	var count int
	for {
		select {
		case <- time.After(time.Second*30):
			fmt.Printf("[%d]ServiceTask Do.\n", ctx.Id)
			count++
		}
		if count == 10 {
			break
		}
	}

	return nil, nil
}

type UserCommandTask struct {
	Name string `json:"name"`
	Tel string 	`json:"tel"`
	Age int 	`json:"age"`
}

func (ct *UserCommandTask)Do(ctx *task.TaskContext)([]byte, error) {
	resp, err := json.Marshal(ct)
	fmt.Printf("[%d] handling.\n", ctx.Id)
	return resp, err
}

func (ct *UserCommandTask)Clone() task.Task {
	task := new(UserCommandTask)
	return task
}

func (ct *UserCommandTask)Prepare(ctx *task.TaskContext) ([]byte, error) {
	if ctx.R.Method != "POST" {
		ctx.W.WriteHeader(400)
		return []byte("Invalid Method"), errors.New("Invalid Method")
	}

	data, _ := ioutil.ReadAll(ctx.R.Body)
	fmt.Printf("[%d] get req:%v\n", ctx.Id, string(data))
	err := json.Unmarshal(data, ct)
	if err != nil {
		return []byte(err.Error()), err
	}

	return nil, nil
}

func (ct *UserCommandTask)Response(ctx *task.TaskContext, data []byte) {
	if data != nil {
		ctx.W.Write(data)
	}
	fmt.Printf("[%d] write resp :%s\n", ctx.Id, string(data))
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()

	tp.AddServiceRoute(&ServiceTask{})
	tp.AddCommandRoute("/test", &UserCommandTask{})

	tp.Run(":8086")
}

