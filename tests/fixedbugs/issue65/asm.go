//go:build ignore

package main

import (
	. "github.com/mihailkirov/avo/build"
	. "github.com/mihailkirov/avo/operand"
	. "github.com/mihailkirov/avo/reg"
)

func main() {
	TEXT("Issue65", NOSPLIT, "func()")
	VINSERTI128(Imm(1), Y0.AsX(), Y1, Y2)
	RET()
	Generate()
}
