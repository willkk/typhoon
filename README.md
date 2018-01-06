# Typhoon
A general purpose web api framework. Its main purposes are:
1. Provide a simple, flexible and extensible framework aiming at Go Web API programs.
2. Raise the reusability of codes between different projects.
3. Identify go routines by Id, making it easy to track a unique request in log.


It's based on Layered Architecture thought:
```
+-----------------------------------------------+
|		 |				|
|  Application	 |   Prepare/Response/Finish	|
|		 |				|
|-----------------------------------------------|
|		 |				|
|    Domain	 |             Do		|
|		 |				|
|-----------------------------------------------|
|		 |				|
| Infrastructure |     Db, MQ, Logging, ...	|
|		 |				|
+-----------------------------------------------+

Finish is used to call downstream services if Do completes successfully. You can also leave it empty freely.
		     +---------------+
		     |    Service    |
		     |---------------|
		     |     Finish    |
		     +---------------+
		    /		      \
		   /		       \
  +---------------+			+---------------+   
  |		  |			|		|
  |   Service2	  |			|    Service3	|
  |		  |			|		|
  +---------------+			+---------------+
```  
该框架主要目的是：
1. 解决通用Web API类型的Go程序架构问题，实现简单，灵活，易于扩展。
2. 解决项目间重复编码的问题，提高代码复用性。
3. 为go routine标记id，方便log跟踪单个请求。

该框架基于分层架构模型:  

| 成员函数 | 层次   |  说明  |
| --------   | :-----:  | ----  |
|  Prepare   |   应用层<br>(Application)	|  用于接收客户端请求，并进行身份验证，授权验证，参数检查等，然后对数据进行预处理，封装成领域层对象所需的格式。如果失败，返回失败条件下的响应消息和失败原因。
|  Do        |  领域层<br>(Domain)         |  接受Prepare封装好的格式化数据，然后执行领域层逻辑进行业务处理。此函数只关心领域逻辑处理。返回处理结果和失败原因。
| Response   |    应用层<br>(Application)  |  接受Prepare或Do的返回结果数据，将响应消息进行序列化等操作，返回给客户端。
| Finish     |   应用层<br>(Application)   |  在Do操作执行成功时，调用下游服务接口。比如用户支付成功，亚马逊发送下单通知邮件；用户修改送货地址，亚马逊发送地址变更通知邮件等。


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

// Finishing works if there is any. Typically, it calls downstream services.
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
