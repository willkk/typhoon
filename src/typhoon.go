package typhoon

import (
	"net/http"
	"core"
)

type Typhoon struct {
	server http.Server
}

// Build new Typhoon instance
func New() *Typhoon {
	return &Typhoon{
		server: http.Server{
			Handler: core.Router(),
		},
	}
}

