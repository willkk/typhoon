package rest

import "strings"

// Parse URL into Request
// http url format: scheme://host/[path][?params]
// path: [version]/collection:
// 		 [version]/collection/resourceid:
type Request struct {
	Scheme string
	Host string
	Path string
	Paths []string // 0：version 1:collection 2:resourceid
	Params map[string]string
}

func NewRequest()*Request {
	return &Request{Params: make(map[string]string)}
}

func Parse(url string) *Request {
	req := NewRequest()
	if strings.HasPrefix(url, "http://") {
		req.Scheme = "http://"
	}
	if strings.HasPrefix(url, "https://") {
		req.Scheme = "https://"
	}
	url = strings.TrimPrefix(url, req.Scheme)
	if strings.Contains(url, "?") {
	 	urls := strings.Split(url, "?")
		url = urls[0]

		params := strings.Split(urls[1], "&")
		for _, param := range params {
			kvs := strings.Split(param, "=")
			req.Params[kvs[0]] = kvs[1]
		}
	}
	paths := strings.Split(url, "/")
	req.Host = paths[0]
	req.Path = strings.TrimPrefix(req.Path, req.Host)
	req.Paths = paths[1:]

	return req
}