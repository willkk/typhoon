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

