package testdata

import (
	"net/http"
)

//iwrapper:target
type (
	TypeInBracketOutsideComment interface {
		//iwrapper:require
		http.ResponseWriter
		http.Hijacker
	}
)
