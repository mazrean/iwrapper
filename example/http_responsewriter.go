package example

//go:generate go run github.com/mazrean/iwrapper -src=$GOFILE -dst=iwrapper_$GOFILE

import (
	"net/http"
)

//iwrapper:target
type ResponseWriter interface {
	//iwrapper:require
	http.ResponseWriter
	http.Hijacker
	http.CloseNotifier
	http.Flusher
}
