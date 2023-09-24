package testdata

import (
	"net/http"
)

//iwrapper:target
type OtherFileDeclare interface {
	//iwrapper:require
	http.ResponseWriter
	Hijacker
}
