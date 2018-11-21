package load

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mmcloughlin/avo/internal/inst"
	"github.com/mmcloughlin/avo/internal/opcodescsv"
	"github.com/mmcloughlin/avo/internal/opcodesxml"
)

const (
	defaultCSVName        = "x86.v0.2.csv"
	defaultOpcodesXMLName = "x86_64.xml"
)

type Loader struct {
	X86CSVPath     string
	OpcodesXMLPath string

	alias map[opcodescsv.Alias]string
}

func NewLoaderFromDataDir(dir string) *Loader {
	return &Loader{
		X86CSVPath:     filepath.Join(dir, defaultCSVName),
		OpcodesXMLPath: filepath.Join(dir, defaultOpcodesXMLName),
	}
}

func (l *Loader) Load() ([]*inst.Instruction, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	// Load Opcodes XML file.
	iset, err := opcodesxml.ReadFile(l.OpcodesXMLPath)
	if err != nil {
		return nil, err
	}

	// Load opcodes XML data, grouped by Go opcode.
	im := map[string]*inst.Instruction{}
	for _, i := range iset.Instructions {
		for _, f := range i.Forms {
			if !l.include(f) {
				continue
			}

			opcode := l.goname(f)

			if im[opcode] == nil {
				im[opcode] = &inst.Instruction{
					Opcode:  opcode,
					Summary: i.Summary,
				}
			}

			im[opcode].Forms = append(im[opcode].Forms, l.form(f))
		}
	}

	// Convert to a slice to return.
	is := make([]*inst.Instruction, 0, len(im))
	for _, i := range im {
		is = append(is, i)
	}

	return is, nil
}

func (l *Loader) init() error {
	icsv, err := opcodescsv.ReadFile(l.X86CSVPath)
	if err != nil {
		return err
	}

	l.alias, err = opcodescsv.BuildAliasMap(icsv)
	if err != nil {
		return err
	}

	return nil
}

// include decides whether to include the instruction form in the avo listing.
// This discards some opcodes that are not supported in Go.
func (l Loader) include(f opcodesxml.Form) bool {
	// Exclude certain ISAs simply not present in Go.
	for _, isa := range f.ISA {
		switch isa.ID {
		case "TBM", "CLZERO", "MONITORX", "FEMMS":
			return false
		}
	}

	// x86 csv contains a number of CMOV* instructions which are actually not valid
	// Go instructions. The valid Go forms should have different opcodes from GNU.
	// Therefore a decent "heuristic" is CMOV* instructions that do not have
	// aliases.
	if strings.HasPrefix(f.GASName, "cmov") && l.lookupAlias(f) == "" {
		return false
	}

	// Some specific exclusions.
	switch f.GASName {
	case "callq":
		return false
	}

	return true
}

func (l Loader) lookupAlias(f opcodesxml.Form) string {
	a := opcodescsv.Alias{Opcode: f.GASName, DataSize: datasize(f)}
	return l.alias[a]
}

func (l Loader) goname(f opcodesxml.Form) string {
	// Use go opcode from Opcodes XML where available.
	if f.GoName != "" {
		return f.GoName
	}

	// Return alias if available.
	if a := l.lookupAlias(f); a != "" {
		return a
	}

	n := strings.ToUpper(f.GASName)

	// Some need data sizes added to them.
	// TODO(mbm): is there a better way of determining which ones these are?
	s := datasize(f)
	suffix := map[int]string{16: "W", 32: "L", 64: "Q"}
	switch n {
	case "RDRAND", "RDSEED":
		n += suffix[s]
	}

	return n
}

func (l Loader) form(f opcodesxml.Form) inst.Form {
	return inst.Form{
		Operands: operands(f.Operands),
	}
}

// operands maps Opcodes XML operands to avo format.
func operands(ops []opcodesxml.Operand) []inst.Operand {
	r := make([]inst.Operand, 0, len(ops))
	for _, op := range ops {
		r = append(r, operand(op))
	}
	return r
}

// operand maps an Opcodes XML operand to avo format.
func operand(op opcodesxml.Operand) inst.Operand {
	return inst.Operand{
		Type:   op.Type,
		Action: inst.ActionFromReadWrite(op.Input, op.Output),
	}
}

// datasize (intelligently) guesses the datasize of an instruction form.
func datasize(f opcodesxml.Form) int {
	for _, op := range f.Operands {
		if !op.Output {
			continue
		}
		for s := 8; s <= 64; s *= 2 {
			if strings.HasSuffix(op.Type, strconv.Itoa(s)) {
				return s
			}
		}
	}
	return 0
}
