package deqcc

import (
	"bufio"
	"encoding/binary"
	"io"
)

type ProgramsReader struct {
	Programs
	io.ReaderAt
}

func (p *ProgramsReader) DumpStrings(sr *StringRepo) error {
	var (
		r = bufio.NewReader(io.NewSectionReader(p, int64(p.Strings.Offset), int64(p.Strings.Num)))
		i int
	)
	for {
		data, err := r.ReadBytes(0)
		switch err {
		case nil:
			str := string(data)
			str = str[:len(str)-1]
			sr.Add(str, i)
			i += len(data)
		case io.EOF:
			str := string(data)
			str = str[:len(str)-1]
			sr.Add(str, i)
			return nil
		default:
			return err
		}
	}
	return nil
}

func (p *ProgramsReader) GetString(at int) (string, error) {
	r := bufio.NewReader(io.NewSectionReader(p, int64(p.Strings.Offset)+int64(at), int64(p.Strings.Num)))
	s, err := r.ReadString(0)
	if err != nil {
		return "", err
	}
	return s[:len(s)-1], nil
}

func (p *ProgramsReader) GetStatements() ([]Statement, error) {
	stmts := make([]Statement, p.Statements.Num)
	r := io.NewSectionReader(p, int64(p.Statements.Offset), int64(p.Statements.Num)*StatementSz)
	if err := read(r, stmts); err != nil {
		return nil, err
	}
	return stmts, nil
}

func (p *ProgramsReader) GetGlobalDefs() ([]Def, error) {
	defs := make([]Def, p.GlobalDefs.Num)
	r := io.NewSectionReader(p, int64(p.GlobalDefs.Offset), int64(p.GlobalDefs.Num)*DefSz)
	if err := read(r, defs); err != nil {
		return nil, err
	}
	return defs, nil
}

func (p *ProgramsReader) GetGlobals() (*GlobalVars, error) {
	var globals GlobalVars
	r := io.NewSectionReader(p, int64(p.Globals.Offset), GlobalVarsSz)
	if err := read(r, &globals); err != nil {
		return nil, err
	}
	return &globals, nil
}

func (p *ProgramsReader) GetEvals() ([]Eval, error) {
	evals := make([]Eval, p.Globals.Num)
	r := io.NewSectionReader(p, int64(p.Globals.Offset), int64(p.Globals.Num)*EvalSz)
	if err := read(r, evals); err != nil {
		return nil, err
	}
	return evals, nil
}

type Value struct {
	Immediate bool
	Name      string
	Type      Type
	Data      Data
}

type Data interface {
	Repr() string
}

type Program struct {
	strings *StringRepo
}

func (p *Program) NumStrings() int                     { return p.strings.Num() }
func (p *Program) String(i int) (string, bool)         { return p.strings.ById(i) }
func (p *Program) StringByOffset(i int) (string, bool) { return p.strings.ByOffset(i) }

type StringRepo struct {
	strs    []string
	strOffs map[int]int
}

func (r *StringRepo) Add(str string, offset int) {
	if r.strOffs == nil {
		r.strOffs = make(map[int]int)
	}
	i := len(r.strs)
	r.strs = append(r.strs, str)
	r.strOffs[offset] = i
}

func (r *StringRepo) Num() int {
	if r == nil {
		return 0
	}
	return len(r.strs)
}

func (r *StringRepo) ById(i int) (string, bool) {
	if r == nil || i < 0 || i >= len(r.strs) {
		return "", false
	}
	return r.strs[i], true
}

func (r *StringRepo) ByOffset(i int) (string, bool) {
	if r == nil {
		return "", false
	}
	i, ok := r.strOffs[i]
	if !ok {
		return "", false
	}
	return r.strs[i], true
}

func Open(r io.ReaderAt) (*Program, error) {
	var hdr Programs
	{
		r := io.NewSectionReader(r, 0, ProgramsSize)
		if err := read(r, &hdr); err != nil {
			return nil, err
		}
	}
	pr := ProgramsReader{hdr, r}
	// defs, err := pr.GetGlobalDefs()
	// if err != nil {
	// 	return nil, err
	// }
	// evals, err := pr.GetEvals()
	// if err != nil {
	// 	return nil, err
	// }
	// for _, d := range defs {
	// 	edef, err := NewEDef(d, &pr, evals)
	// 	if edef != nil && err == nil {
	// 		n, err := pr.GetString(int(edef.SName))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		fmt.Println(n, edef.Data)
	// 	}
	// }

	p := Program{
		strings: new(StringRepo),
	}
	if err := pr.DumpStrings(p.strings); err != nil {
		return nil, err
	}
	return &p, nil
}

func read(r io.Reader, data interface{}) error { return binary.Read(r, binary.LittleEndian, data) }
