package rest

import "testing"

func TestURLParser_Parse(t *testing.T) {
	url := "http://api.github.com/v1/users/willkk:?type=3"
	t.Log(Parse(url))
}
