package testdata

import (
	"net/http"
)

//iwrapper:target
type MultiOptional interface {
	//iwrapper:require
	http.ResponseWriter
	http.Hijacker
	http.Flusher
}
