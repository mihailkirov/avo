package printer_test

import (
	"testing"
	"github.com/mihailkirov/avo/build"
	"github.com/mihailkirov/avo/buildtags"
	"github.com/mihailkirov/avo/printer"
)

func TestStubsPragmas(t *testing.T) {
	ctx := build.NewContext()
	ctx.Function("f")
	ctx.Pragma("noescape")
	ctx.Pragma("linkname f remote.f")
	ctx.SignatureExpr("func(x *uint64)")
	ctx.RET()

	AssertPrintsLines(t, ctx, printer.NewStubs, []string{
		"// Code generated by avo. DO NOT EDIT.",
		"",
		"package printer",
		"",
		"//go:noescape",
		"//go:linkname f remote.f",
		"func f(x *uint64)",
		"",
	})
}

func TestStubsConstraints(t *testing.T) {
	ctx := build.NewContext()
	ctx.ConstraintExpr("linux darwin")
	ctx.ConstraintExpr("amd64 arm64 mips64x ppc64x")

	expect := []string{
		"// Code generated by avo. DO NOT EDIT.",
		"",
	}
	if buildtags.GoBuildSyntaxSupported() {
		expect = append(expect,
			"//go:build (linux || darwin) && (amd64 || arm64 || mips64x || ppc64x)",
		)
	}
	if buildtags.PlusBuildSyntaxSupported() {
		expect = append(expect,
			"// +build linux darwin",
			"// +build amd64 arm64 mips64x ppc64x",
		)
	}
	expect = append(expect,
		"",
		"package printer",
		"",
	)

	AssertPrintsLines(t, ctx, printer.NewStubs, expect)
}
