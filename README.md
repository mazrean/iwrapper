# iwrapper
[![](https://github.com/mazrean/iwrapper/workflows/Release/badge.svg)](https://github.com/mazrean/iwrapper/actions)
[![go report](https://goreportcard.com/badge/mazrean/iwrapper)](https://goreportcard.com/report/mazrean/iwrapper)

[日本語](./README.ja.md)

A tool to assist in creating wrappers for Go interfaces.

## Motivation
Consider creating a wrapper for `http.ResponseWriter`. The simplest approach is to use embedding as shown below:
```go
type MyResponseWriter struct {
  http.ResponseWriter
}

func WrapResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
  return &MyResponseWriter{rw}
}
```
However, with this method, the result of type assertion changes after wrapping.
```go
// ResponseWriter implemented with hijacker
var rw http.ResponseWriter

// Before wrap: ok => true
_, ok = rw.(http.Hijacker)

// After wrap: ok => false
_, ok = WrapResponseWriter(rw).(http.Hijacker)
```
Such type assertions, often performed in standard libraries to maintain backward compatibility, can lead to unexpected behavioral changes.

By using the function (`ResponseWriterWrapper`) generated with `iwrapper`, this problem can be addressed as shown below:
```go
// 1. Implement http.Hijacker in MyResponseWriter
func (w *MyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// 2. Use ResponseWriterWrapper for wrapping
func WrapResponseWriterWithIWrapper(rw http.ResponseWriter) http.ResponseWriter {
	return ResponseWriterWrapper(rw, func(rw http.ResponseWriter) ResponseWriter {
		return &MyResponseWriter{rw}
	})
}

// After wrapping using iwrapper: ok => true
_, ok = WrapResponseWriterWithIWrapper(rw).(http.Hijacker)
```

## Usage
Consider wrapping `http.ResponseWriter` so that the results of type assertions for `http.Hijacker`, `http.CloseNotifier`, and `http.Flusher` remain unchanged:
1. Create the following Go file as the configuration for iwrapper:
   ```go
   package testdata

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
   ```
   - You can customize the generated function name with `iwrapper:target func:"ResponseWriterWrapFunc"`.
2. Execute `go generate`.
   - This produces the wrapping function (`ResponseWriterWrapper`) in `iwrapper_<configuration filename>.go`.

Detailed generated code and its usage examples are available in [`/example/`](https://github.com/mazrean/iwrapper/tree/main/example/).

## License

MIT
