package testdata

import (
	"net/http"
)

type (
	//iwrapper:target
	MultiTargetInBracket1 interface {
		//iwrapper:require
		http.ResponseWriter
		http.Hijacker
	}

	//iwrapper:target
	MultiTargetInBracket2 interface {
		//iwrapper:require
		http.ResponseWriter
		http.Flusher
	}
)
