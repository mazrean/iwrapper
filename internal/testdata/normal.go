package testdata

import (
	"net/http"
)

//iwrapper:target
type Normal interface {
	//iwrapper:require
	http.ResponseWriter
	http.Hijacker
}
