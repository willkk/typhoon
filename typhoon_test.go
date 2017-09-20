package typhoon

import (
	"testing"
	"typhoon/core"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
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

type CommandTask struct {
	w http.ResponseWriter
	r *http.Request
}

func (ct *CommandTask)Do()(interface{}, error) {
	fmt.Println("ServiceTask.Do()")
	time.Sleep(time.Second*3)
	return nil, nil
}

func (ct *CommandTask)Clone() core.Task {
	return nil
}

func (ct *CommandTask)Prepare(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		return nil, nil
	}
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Println("get req:%s", string(data))
	fmt.Println("Prepare")
	ct.w = w
	ct.r = r

	return nil, nil
}

func (ct *CommandTask)Response(data interface{}) error {
	ct.w.Write([]byte("123456"))
	return nil
}


func TestTyphoon_Run(t *testing.T) {
	tp := New()
	tp.AddServiceRoute("timer/print", &ServiceTask{})
	tp.AddCommandRoute("/test", &CommandTask{})

	tp.Run(":8086")
}

