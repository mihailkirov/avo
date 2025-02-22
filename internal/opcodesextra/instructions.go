// Package opcodesextra provides curated extensions to the instruction database.
package opcodesextra

import (
	"sort"
	"github.com/mihailkirov/avo/internal/inst"
)

// sets of extra instructions.
var sets = [][]*inst.Instruction{
	movlqzx,
	gfni,
	vaes,
	vnni,
	vpclmulqdq,
	vpopcntdq,
	bitalg,
	vbmi2,
}

// Instructions returns a list of extras to add to the instructions database.
//
// Note that instructions returned are expected to be injected into the loading
// process, as if they had been read out of the Opcodes database.  As such, they
// are not expected to be in the final form required for the instruction
// database. For example, AVX-512 instruction form transformations do not need
// to be applied, and operand types such as xmm{k}{z} or m256/m64bcst may be
// used for brevity.
func Instructions() []*inst.Instruction {
	// Concatenate and clone the instruction lists.  It can be convenient for
	// forms lists and other data structures to be shared in the curated lists,
	// but we want to return distinct copies here to avoid subtle bugs in
	// consumers.
	var is []*inst.Instruction
	for _, set := range sets {
		for _, i := range set {
			c := i.Clone()
			is = append(is, &c)
		}
	}

	// Sort ISA lists. Similarly, this facilitates sharing helper functions for
	// building forms lists without worrying about whether the ISA list is in
	// the right order.
	for _, i := range is {
		for idx := range i.Forms {
			sort.Strings(i.Forms[idx].ISA)
		}
	}

	return is
}
