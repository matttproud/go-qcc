package deqcc

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type ProgramsReader struct {
	Programs
	io.ReaderAt
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

func Open(r io.ReaderAt) error {
	var hdr Programs
	{
		r := io.NewSectionReader(r, 0, ProgramsSize)
		if err := read(r, &hdr); err != nil {
			return err
		}
	}
	pr := ProgramsReader{hdr, r}
	defs, err := pr.GetGlobalDefs()
	if err != nil {
		return err
	}
	for i, d := range defs {
		n, err := pr.GetString(int(d.SName))
		if err != nil {
			return err
		}
		fmt.Printf("%v. %v (%v)\n", i, d, n)
	}
	return nil
}

func read(r io.Reader, data interface{}) error { return binary.Read(r, binary.LittleEndian, data) }
