package web

import "net/http"

func init() {
	initHandle()
}

// StartService with http service
func StartService(address string) error {
	return http.ListenAndServe(address, baseSignUPHandle())
}

// StartTLSService with http TLS service
func StartTLSService(address string, certFile, keyFile string) error {
	return http.ListenAndServeTLS(address, certFile, keyFile, baseSignUPHandle())
}
