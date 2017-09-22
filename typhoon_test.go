package typhoon

import (
	"testing"
	"typhoon/core"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
)

type ServiceTask struct {

}

func (st *ServiceTask)Do()([]byte, error) {
	return nil, nil
}

type UserCommandTask struct {
	Name string `json:"name"`
	Tel string 	`json:"tel"`
	Age int 	`json:"age"`
}

func (ct *UserCommandTask)Do()([]byte, error) {
	resp, err := json.Marshal(ct)
	return resp, err
}

func (ct *UserCommandTask)Clone() core.Task {
	task := new(UserCommandTask)
	return task
}

func (ct *UserCommandTask)Prepare(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		return []byte("Invalid Method"), errors.New("Invalid Method")
	}
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, ct)
	if err != nil {
		return []byte(err.Error()), err
	}

	return nil, nil
}

func (ct *UserCommandTask)Response(w http.ResponseWriter, data []byte) {
	if data != nil {
		w.Write(data)
	}
}

func TestTyphoon_Run(t *testing.T) {
	tp := New()
	tp.AddServiceRoute("timer/print", &ServiceTask{})
	tp.AddCommandRoute("/test", &UserCommandTask{})

	tp.Run(":8086")
}

