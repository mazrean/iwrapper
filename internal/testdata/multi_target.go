package testdata

import (
	"net/http"
)

//iwrapper:target
type MultiTarget1 interface {
	//iwrapper:require
	http.ResponseWriter
	http.Hijacker
}

//iwrapper:target
type MultiTarget2 interface {
	//iwrapper:require
	http.ResponseWriter
	http.Flusher
}
