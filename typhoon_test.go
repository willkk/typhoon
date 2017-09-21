package typhoon

import (
	"testing"
	"typhoon/core"
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
)

type ServiceTask struct {

}

func (st *ServiceTask)Do()(interface{}) {
	return nil
}

func (st *ServiceTask)Clone() core.Task {
	return nil
}

type UserCommandTask struct {
	Name string `json:"name"`
	Tel string 	`json:"tel"`
	Age int 	`json:"age"`
}

func (ct *UserCommandTask)Do()(interface{}) {
	resp, _ := json.Marshal(ct)
	return resp
}

func (ct *UserCommandTask)Clone() core.Task {
	task := new(UserCommandTask)
	*task = *ct
	return task
}

func (ct *UserCommandTask)Prepare(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(400)
		return errors.New("Invalid Method")
	}
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, ct)
	if err != nil {
		return err
	}

	return nil
}

func (ct *UserCommandTask)Response(w http.ResponseWriter, data []byte) error {
	if data != nil {
		w.Write(data)
	}

	return nil
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()
	tp.AddServiceRoute("timer/print", &ServiceTask{})
	tp.AddCommandRoute("/test", &UserCommandTask{})

	tp.Run(":8086")
}

