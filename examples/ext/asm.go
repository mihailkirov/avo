//go:build ignore

package main

import . "github.com/mihailkirov/avo/build"

func main() {
	Package("github.com/mihailkirov/avo/examples/ext")
	Implement("StructFieldB")
	b := Load(Param("e").Field("B"), GP8())
	Store(b, ReturnIndex(0))
	RET()
	Generate()
}
