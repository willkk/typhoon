package typhoon

import (
	"testing"
	"net/http"
	"core"
)

func TestNew(t *testing.T) {
	tp := New()
	var r *http.Request = nil
	task := core.NewTask(nil, nil)
	tp.AddRoute("/status", )
}
