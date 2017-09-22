# Typhoon
A general purpose web api framework. Its main purposes are:
1. Provide a simple, flexible and extensible framework aiming at Go Web API programs.
2. Raise the reusability of codes between different projects.
3. Identify go routines by Id, making it easy to track a unique request in log.

Typhoon是一个通用目的的Web API应用框架。实际项目中，发现大家用Go写的程序架构性不强，一部分人还停留在面向过程的世界里；另一部分人有面向对象的思想，但程序的扩展性和灵活性不够好；或者自己已经可以写出具有良好结构的项目应用，但是系统库不够灵活，项目之间重复代码比较多等等。

针对以上问题，该框架主要目的是：
1. 解决通用Web API类型的Go程序架构问题，实现简单，灵活，易于扩展。
2. 解决项目间重复编码的问题，提高代码复用性。
3. 为go routine标记id，方便log跟踪单个请求。

# Example

```go
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
// Implement Task interface
type ServiceTask struct {
}

func (st *ServiceTask)Do(ctx *task.TaskContext) {
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

func (ct *UserCommandTask)Do(ctx *task.TaskContext)([]byte, error) {
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
func (ct *UserCommandTask)Prepare(ctx *task.TaskContext) ([]byte, error) {
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

func (ct *UserCommandTask)Response(ctx *task.TaskContext, data []byte) {
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
	tp.AddTask(&ServiceTask{})
	tp.AddTask(&ServiceTask{})
	// Add web command task
	tp.AddRoute("/test", &UserCommandTask{})

	// start service tasks
	tp.StartTasks()
	// wait for web requests
	tp.Run(":8086")
}
```

You can send requests using command: **curl -d '{"name":"will","age":23, "tel":"112"}' http://127.0.0.1:8086/test**

**Output may be like this:**
```
[3] get req:{"name":"will","age":23, "tel":"112"}
[3] handling.
[3] write resp :{"name":"","tel":"","age":0}.
[3] done, consume 1683 us
[4] get req:{"name":"will","age":23, "tel":"112"}
[4] handling.
[4] write resp :{"name":"","tel":"","age":0}.
[4] done, consume 69 us
[1] service task count 0.
[2] service task count 0.
[2] service task count 1.
[1] service task count 1.
[1] service task count 2.
[2] service task count 2.
[2] service task count 3.
[1] service task count 3.
[1] service task count 4.
[5] get req:{"name":"will","age":23, "tel":"112"}
[5] handling.
[5] write resp :{"name":"","tel":"","age":0}.
[5] done, consume 380 us
[6] get req:{"name":"will","age":23, "tel":"112"}
[6] handling.
[6] write resp :{"name":"","tel":"","age":0}.
[6] done, consume 440 us
```
