package testdata

import (
	"net/http"
)

type (
	//iwrapper:target
	TypeInBracketInsideComment interface {
		//iwrapper:require
		http.ResponseWriter
		http.Hijacker
	}
)
