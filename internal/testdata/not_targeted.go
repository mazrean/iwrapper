package testdata

import (
	"net/http"
)

type NotTargeted interface {
	http.ResponseWriter
	http.Hijacker
}
