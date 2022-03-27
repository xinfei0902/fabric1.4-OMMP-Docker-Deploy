package dservice

import "net/http"

// HTTP service
type makeHTTPHandleFunc func(debug, db bool) http.HandlerFunc

type httpHandleFuncPair struct {
	p string
	f makeHTTPHandleFunc
}

type makeHTTPHandle func(debug, db bool) http.Handler

type httpHandlePair struct {
	p string
	f makeHTTPHandle
}
