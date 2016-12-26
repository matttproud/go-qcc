package deqcc

import (
	"flag"
	"os"
	"testing"
)

var progs = flag.String("progs", "./progs.dat", "path of progs.dat")

func TestOpenSmoke(t *testing.T) {
	if *progs == "" {
		t.Skipf("no progs.dat")
	}
	f, err := os.Open(*progs)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := Open(f); err != nil {
		t.Fatal(err)
	}
}
