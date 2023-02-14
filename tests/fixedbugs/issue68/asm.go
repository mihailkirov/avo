//go:build ignore

package main

import (
	. "github.com/mihailkirov/avo/build"
	. "github.com/mihailkirov/avo/operand"
)

func main() {
	TEXT("Issue68", NOSPLIT, "func() uint64")
	Doc("Issue68 tests custom package names.")
	x := GP64()
	MOVQ(U32(68), x)
	Store(x, ReturnIndex(0))
	RET()
	Generate()
}
