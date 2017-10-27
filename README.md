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

# Example(typhoon_test.go)

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
	"math/rand"
)

const (
	Err_Success = iota
	Err_Method
	Err_Json
)

//--------------------------------service task----------------------------------
// Implement func(ctx *task.Context)
func TrivialTask(ctx *task.Context) {
	var count int
	if ctx.UserContext != nil {
		fmt.Printf("[%d] userctx:%v.\n", ctx.Id, ctx.UserContext.Value("now"))
	}

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
type PaymentTask struct {
	OrderId int `json:"order_id"`
	SrcBankNo string `json:"src_bank_no"`
	DstBankNo string `json:"dst_bank_no"`
	Amount int `json:"amount"`
}

type PaymentTaskResp struct {
	Code int `json:"code"`
	Err string `json:"err"`
	Data interface{} `json:"data"`
}

func (ptr *PaymentTaskResp)Response() []byte {
	resp, err := json.Marshal(ptr)
	if err != nil {
		resp := fmt.Sprintf(`{"code":%d, "err":"%s"}`, Err_Json, err)
		return []byte(resp)
	}
	return resp
}

func (pt *PaymentTask)Clone() task.CommandTask {
	task := new(PaymentTask)
	return task
}

type userContext struct {
	start int // us
}
func (pt *PaymentTask)Prepare(ctx *task.WebContext) (task.TaskResponse, error) {
	if ctx.R.Method != "POST" {
		ctx.W.WriteHeader(400)
		return &PaymentTaskResp{
			Err_Method,
			"Invalid Method",
			nil},
			errors.New("Invalid Method")
	}

	now := time.Now()

	data, _ := ioutil.ReadAll(ctx.R.Body)
	fmt.Printf("[%d] get req:%v\n", ctx.Id, string(data))
	err := json.Unmarshal(data, pt)
	if err != nil {
		return &PaymentTaskResp{Err_Json, "json.Unmarshal failed", nil}, err
	}

	uctx := &userContext{now.Second()*1000000 + now.Nanosecond()/1000}
	ctx.UserContext = context.WithValue(nil, "user_ctx", uctx)

	return &PaymentTaskResp{Err_Success, "success", ""}, nil
}

func (pt *PaymentTask)Do(ctx *task.WebContext)(task.TaskResponse, error) {
	fmt.Printf("[%d] handling payment.\n", ctx.Id)

	// Do payment business logic.
	sleep := time.Duration(rand.Int()%100)
	time.Sleep(time.Millisecond*sleep)

	return &PaymentTaskResp{Err_Success, "success", "user data"}, nil
}

// Finishing works if there is any. There is no problem if you keep it empty.
func (pt *PaymentTask)Finish(ctx *task.WebContext, reps task.TaskResponse) {
	// add bonus points
	// ...

	// send payment success mails
	// ...

	fmt.Printf("[%d] Finish done.\n", ctx.Id)
}

func (pt *PaymentTask)Response(ctx *task.WebContext, resp task.TaskResponse) {
	now := time.Now()
	if resp != nil {
		data := resp.Response()
		ctx.W.Write(data)
		fmt.Printf("[%d] write resp :%s.\n", ctx.Id, string(data))
	}

	uctx := ctx.UserContext.Value("user_ctx").(*userContext)
	now_us := now.Second()*1000000 + now.Nanosecond()/1000
	fmt.Printf("[%d] task is done, consume %d us\n", ctx.Id, now_us - uctx.start)
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()

	// Add normal service task
	tp.AddTask(TrivialTask)
	// Add web command task
	tp.AddRoute("/test", &PaymentTask{})

	// start service tasks
	tp.StartTasks()
	// start task immediately
	ExecTask(TrivialTask, context.WithValue(nil, "now", time.Now().Unix()))
	// wait for web requests
	tp.Run(":8086")
}

```

You can send requests using command: **curl -d '{"name":"will","age":23, "tel":"112"}' http://127.0.0.1:8086/test**

**Output may be like this:**
```
[2] userctx:1509088227.
[1] service task count 0.
[2] service task count 0.
[3] get req:{"order_id":123,"src_bank_no":"6147258369123","dst_bank_no":"7412589631203","amount":1000}
[3] handling payment.
[3] write resp :{"code":0,"err":"success","data":"user data"}.
[3] task is done, consume 11223 us
[3] Finish done.
[4] get req:{"order_id":123,"src_bank_no":"6147258369123","dst_bank_no":"7412589631203","amount":1000}
[4] handling payment.
[4] write resp :{"code":0,"err":"success","data":"user data"}.
[4] task is done, consume 51472 us
[4] Finish done.
[2] service task count 1.
[1] service task count 1.
[1] service task count 2.
[2] service task count 2.

```
