//go:build ignore

package main

import . "github.com/mkirov/avo/build"

func main() {
	Package("github.com/mkirov/avo/tests/fixedbugs/issue62")
	Implement("private")
	RET()
	Generate()
}
