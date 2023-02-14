//go:build ignore

package main

import (
	. "github.com/mkirov/avo/build"
	. "github.com/mkirov/avo/operand"
	. "github.com/mkirov/avo/reg"
)

func main() {
	TEXT("Issue65", NOSPLIT, "func()")
	VINSERTI128(Imm(1), Y0.AsX(), Y1, Y2)
	RET()
	Generate()
}
