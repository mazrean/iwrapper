package testdata

import (
	"net/http"
)

//iwrapper:target
type MultiRequire interface {
	//iwrapper:require
	http.ResponseWriter
	//iwrapper:require
	http.Hijacker
	http.Flusher
}
