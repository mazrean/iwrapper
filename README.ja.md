# iwrapper
[![](https://github.com/mazrean/iwrapper/workflows/Release/badge.svg)](https://github.com/mazrean/iwrapper/actions)
[![go report](https://goreportcard.com/badge/mazrean/iwrapper)](https://goreportcard.com/report/mazrean/iwrapper)

[English](./README.md)

Goのinterfaceのラッパーを作成するのを補助するツールです。

## Motivation
`http.ResponseWriter`のラッパーを作成することを考えてみましょう。
もっとも簡単なのは以下のようにembedを利用する方法です。
```go
type MyResponseWriter struct {
  http.ResponseWriter
}

func WrapResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
  return MyResponseWriter{rw}
}
```
しかし、この方法では型アサーションの結果がwrap前と変わってしまいます。
```go
// hijackerが実装されたResponseWriter
var rw http.ResponseWriter

// wrap前: ok => true
_, ok = rw.(http.Hijacker)

// wrap後: ok => false
_, ok = WrapResponseWriter(rw).(http.Hijacker)
```
このような型アサーションは後方互換性維持のため標準ライブラリなどで良く行われており、不意の挙動の変更を招きます。

`iwrapper`で生成した関数(`ResponseWriterWrapper`)を利用すると以下のようにしてこの問題を解決できます。
```go
// 1. http.HijackerをMyResponseWriterに実装
func (w MyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// 2. ResponseWriterWrapperを利用してwrap
func WrapResponseWriterWithIWrapper(rw http.ResponseWriter) http.ResponseWriter {
	return ResponseWriterWrapper(rw, func(rw http.ResponseWriter) ResponseWriter {
		return MyResponseWriter{rw}
	})
}

// wrap(iwrapper利用)後: ok => true
_, ok = WrapResponseWriterWithIWrapper(rw).(http.Hijacker)
```

## Usage
`http.ResponseWriter`を`http.Hijacker`、`http.CloseNotifier`、`http.Flusher`での型アサーション結果が変わらないようWrapする場合を考えます。
1. iwrapperの設定として以下のGoファイルを作成します。
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
   - `iwrapper:target func:"ResponseWriterWrapFunc"のようにして、生成関数名をカスタマイズできます。
2. `go generate`を実行します
   - `iwrapper_<設定ファイル名>.go`にwrap用関数(`ResponseWriterWrapper`)が生成されます

詳細な生成コード・生成コードの使用例は[`/example/`](./example/)にあります。

## License

MIT
