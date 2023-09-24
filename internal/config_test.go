package main

import "github.com/google/go-cmp/cmp"

func diff[T any](x, y T) string {
	return cmp.Diff(x, y, cmp.AllowUnexported(Package{}, AnonymousInterface{}, NamedInterface{}, Interface{}))
}
