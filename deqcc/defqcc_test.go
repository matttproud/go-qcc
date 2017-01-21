package deqcc

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var progs = flag.String("progs", "testdata/progs.dat", "path of progs.dat")

func TestOpenSmoke(t *testing.T) {
	if *progs == "" {
		t.Skipf("no progs.dat")
	}
	f, err := os.Open(*progs)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p, err := Open(f)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < p.NumStrings(); i++ {
		fmt.Println(p.String(i))
	}
}
