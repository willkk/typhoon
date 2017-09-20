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

func (st *ServiceTask)Do()([]byte, error) {
	fmt.Println("ServiceTask.Do()")
	return nil, nil
}

func (st *ServiceTask)Clone() core.Task {
	return nil
}

type UserCommandTask struct {
	Name string `json:"name"`
	Tel string 	`json:"tel"`
	Age int 	`json:"age"`
}

func (ct *UserCommandTask)Do()([]byte, error) {
	resp, _ := json.Marshal(ct)
	return resp, nil
}

func (ct *UserCommandTask)Clone() core.Task {
	return nil
}

func (ct *UserCommandTask)Prepare(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		fmt.Println("Prepare err.")
		return nil, errors.New("Invalid Method")
	}
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Println("get req:", string(data))
	err := json.Unmarshal(data, ct)
	if err != nil {
		fmt.Println("Unmarshal failed. err=", err)
		return nil, err
	}
	fmt.Println("get object:", ct)

	return nil, nil
}

func (ct *UserCommandTask)Response(w http.ResponseWriter, data []byte) error {
	fmt.Println("Response data:", string(data))
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

