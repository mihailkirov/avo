//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"github.com/mihailkirov/avo/attr"
)

var (
	output = flag.String("output", "", "path to output file (default stdout)")
	seed   = flag.Int64("seed", 0, "random seed")
	num    = flag.Int("num", 32, "number of attributes to generate")
)

func GenerateAttributes(n int) []attr.Attribute {
	as := make([]attr.Attribute, 0, n)

	// Include each bitlevel.
	for i := 0; i < 16 && i < n; i++ {
		a := attr.Attribute(1 << uint(i))
		as = append(as, a)
	}

	// Add randomly generated attributes.
	for len(as) < n {
		a := attr.Attribute(rand.Uint32())
		as = append(as, a)
	}

	return as
}

func PrintAttributesTest(w io.Writer, as []attr.Attribute) {
	_, self, _, _ := runtime.Caller(0)
	fmt.Fprintf(w, "// Code generated by %s. DO NOT EDIT.\n\n", filepath.Base(self))
	fmt.Fprintf(w, "#include \"textflag.h\"\n\n")
	fmt.Fprintf(w, "TEXT ·attrtest(SB), $0-1\n")
	fmt.Fprintf(w, "\tMOVB $0, ret+0(FP)\n")

	for i, a := range as {
		fmt.Fprintf(w, "\tMOVW $(%d), R8\n", a)
		fmt.Fprintf(w, "\tMOVW $(%s), R9\n", a.Asm())
		fmt.Fprintf(w, "\tCMPW R8, R9\n")

		cont := fmt.Sprintf("cont%d", i)
		fmt.Fprintf(w, "\tJE   %s\n", cont)
		fmt.Fprintf(w, "\tRET\n")

		fmt.Fprintf(w, "\n%s:\n", cont)
	}

	fmt.Fprintf(w, "\tMOVB $1, ret+0(FP)\n")
	fmt.Fprintf(w, "\tRET\n")
}

func main() {
	flag.Parse()

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = f
	}

	rand.Seed(*seed)
	as := GenerateAttributes(*num)

	buf := bytes.NewBuffer(nil)
	PrintAttributesTest(buf, as)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
	}
}
