//go:build ignore

package main

import . "github.com/mkirov/avo/build"

func main() {
	Package("github.com/mkirov/avo/examples/ext")
	Implement("StructFieldB")
	b := Load(Param("e").Field("B"), GP8())
	Store(b, ReturnIndex(0))
	RET()
	Generate()
}
