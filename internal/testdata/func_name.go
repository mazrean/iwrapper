package testdata

import (
	"net/http"
)

//iwrapper:target func:"FuncNameWrapFunc"
type FuncName interface {
	//iwrapper:require
	http.ResponseWriter
	http.Hijacker
}
