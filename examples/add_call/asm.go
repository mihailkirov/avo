//go:build ignore
// +build ignore

package main

import (
	"reflect"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	// func1
	TEXT("Add", NOSPLIT|NOPTR, "func(x, y uint64) uint64")
	Doc("Adds x and y")
	x := Load(Param("x"), GP64())
	y := Load(Param("y"), GP64())
	ADDQ(x, y)
	Store(y, ReturnIndex(0))
	RET()

	//func 2
	var tmp uint64
	TEXT("Addcall", NOSPLIT, "func(x, y uint64) uint64")
	Doc("Performs a call to Add")
	x = Load(Param("x"), GP64())
	y = Load(Param("y"), GP64())
	z := GP64()
	XORQ(z, z)
	Doc("Allocate space on the stack for the return value + the two arguments")
	AllocLocal(3 * int(reflect.TypeOf(tmp).Size())) // return + 2 args
	MOVQ(z, Mem{Base: SP, Disp: 2 * int(reflect.TypeOf(tmp).Size())})
	MOVQ(x, Mem{Base: SP, Disp: 1 * int(reflect.TypeOf(tmp).Size())})
	MOVQ(y, Mem{Base: SP, Disp: 0 * int(reflect.TypeOf(tmp).Size())})
	w := GP64()
	LEAQ(Mem{Base: StaticBase, Symbol: Symbol{Name: "·Add", Static: false}}, w)
	CALL(w)
	MOVQ(Mem{Base: SP, Disp: 2 * int(reflect.TypeOf(tmp).Size())}, z)
	Store(z, ReturnIndex(0))
	RET()
	Generate()
}
