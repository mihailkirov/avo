//go:build ignore

package main

import (
	. "github.com/mkirov/avo/build"
	. "github.com/mkirov/avo/operand"
	. "github.com/mkirov/avo/reg"
)

func main() {
	TEXT("Formatting", NOSPLIT, "func()")
	Doc("Formatting contains various cases to test the formatter.")

	ADDQ(R8, R8)
	Comment("One comment line between instructions.")
	ADDQ(R8, R8)

	Comment("Comment before label.")
	Label("label")
	Comment("Comment after label.")
	ADDQ(R8, R8)
	JMP(LabelRef("label"))

	RET()

	Generate()
}
