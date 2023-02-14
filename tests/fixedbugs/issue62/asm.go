//go:build ignore

package main

import . "github.com/mihailkirov/avo/build"

func main() {
	Package("github.com/mihailkirov/avo/tests/fixedbugs/issue62")
	Implement("private")
	RET()
	Generate()
}
