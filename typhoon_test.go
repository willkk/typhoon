package typhoon

import (
	"testing"
	"typhoon/core/task"
	"io/ioutil"
	"errors"
	"encoding/json"
	"time"
	"fmt"
	"context"
)

//--------------------------------service task----------------------------------
// Implement func(ctx *task.Context)
func TrivialTask(ctx *task.Context) {
	var count int
	for {
		select {
		case <- time.After(time.Second*10):
			fmt.Printf("[%d] service task count %d.\n", ctx.Id, count)
			count++
		}
		if count == 100 {
			break
		}
	}
}

//--------------------------------command task----------------------------------
// Implement CommandTask interface
type UserCommandTask struct {
	Name string `json:"name"`
	Tel string 	`json:"tel"`
	Age int 	`json:"age"`
}

func (ct *UserCommandTask)Do(ctx *task.WebContext)([]byte, error) {
	resp, err := json.Marshal(ct)
	fmt.Printf("[%d] handling.\n", ctx.Id)
	return resp, err
}

func (ct *UserCommandTask)Clone() task.CommandTask {
	task := new(UserCommandTask)
	return task
}

type userContext struct {
	start int // us
}
func (ct *UserCommandTask)Prepare(ctx *task.WebContext) ([]byte, error) {
	if ctx.R.Method != "POST" {
		ctx.W.WriteHeader(400)
		return []byte("Invalid Method"), errors.New("Invalid Method")
	}

	now := time.Now()

	data, _ := ioutil.ReadAll(ctx.R.Body)
	fmt.Printf("[%d] get req:%v\n", ctx.Id, string(data))
	err := json.Unmarshal(data, ct)
	if err != nil {
		return []byte(err.Error()), err
	}

	uctx := &userContext{now.Second()*1000000 + now.Nanosecond()/1000}
	ctx.UserContext = context.WithValue(nil, "user_ctx", uctx)

	return nil, nil
}

func (ct *UserCommandTask)Response(ctx *task.WebContext, data []byte) {
	now := time.Now()
	if data != nil {
		ctx.W.Write(data)
		fmt.Printf("[%d] write resp :%s.\n", ctx.Id, string(data))
	}

	uctx := ctx.UserContext.Value("user_ctx").(*userContext)
	now_us := now.Second()*1000000 + now.Nanosecond()/1000
	fmt.Printf("[%d] done, consume %d us\n", ctx.Id, now_us - uctx.start)
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()

	// Add normal service task
	tp.AddTask(TrivialTask)
	// Add web command task
	tp.AddRoute("/test", &UserCommandTask{})

	// start service tasks
	tp.StartTasks()
	// start task immediately
	ExecTask(TrivialTask)
	// wait for web requests
	tp.Run(":8086")
}

