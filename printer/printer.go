// Package printer implements printing of avo files in various formats.
package printer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/mihailkirov/avo/internal/stack"
	"github.com/mihailkirov/avo/ir"
)

// Printer can produce output for an avo File.
type Printer interface {
	Print(*ir.File) ([]byte, error)
}

// Builder can construct a printer.
type Builder func(Config) Printer

// Config represents general printing configuration.
type Config struct {
	// Command-line arguments passed to the generator. If provided, this will be
	// included in a code generation warning.
	Argv []string

	// Name of the code generator.
	Name string

	// Name of Go package the generated code will belong to.
	Pkg string
}

// NewDefaultConfig produces a config with Name "avo".
// The package name is guessed from the current directory.
func NewDefaultConfig() Config {
	return Config{
		Name: "avo",
		Pkg:  pkg(),
	}
}

// NewArgvConfig constructs a Config from os.Args.
// The package name is guessed from the current directory.
func NewArgvConfig() Config {
	return Config{
		Argv: os.Args,
		Pkg:  pkg(),
	}
}

// NewGoRunConfig produces a Config for a generator that's expected to be
// executed via "go run ...".
func NewGoRunConfig() Config {
	path := mainfile()
	if path == "" {
		return NewDefaultConfig()
	}
	argv := []string{"go", "run", filepath.Base(path)}
	if len(os.Args) > 1 {
		argv = append(argv, os.Args[1:]...)
	}
	return Config{
		Argv: argv,
		Pkg:  pkg(),
	}
}

// GeneratedBy returns a description of the code generator.
func (c Config) GeneratedBy() string {
	if c.Argv == nil {
		return c.Name
	}
	return fmt.Sprintf("command: %s", strings.Join(c.Argv, " "))
}

// GeneratedWarning returns text for a code generation warning. Conforms to https://golang.org/s/generatedcode.
func (c Config) GeneratedWarning() string {
	return fmt.Sprintf("Code generated by %s. DO NOT EDIT.", c.GeneratedBy())
}

// mainfile attempts to determine the file path of the main function by
// inspecting the stack. Returns empty string on failure.
func mainfile() string {
	if m := stack.Main(); m != nil {
		return m.File
	}
	return ""
}

// pkg guesses the name of the package from the working directory.
func pkg() string {
	if cwd, err := os.Getwd(); err == nil {
		return filepath.Base(cwd)
	}
	return ""
}
