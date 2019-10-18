package m64

import (
	"math/bits"
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/mmcloughlin/ec3/arith/eval"
	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestInterface(t *testing.T) {
	proc := New()
	eval.NewEvaluator(proc)
}

func TestUint64(t *testing.T) {
	proc := New()
	x := rand.Uint64()
	v := proc.Uint64(x)
	if uint64(v.(Word)) != x {
		t.Fail()
	}
}

func TestADD(t *testing.T) {
	proc := New()
	got := func(x, y, c uint64) (uint64, uint64) {
		s, cout := proc.ADD(Word(x), Word(y), Word(c&1))
		return uint64(s.(Word)), uint64(cout.(Word))
	}
	expect := func(x, y, c uint64) (uint64, uint64) {
		return bits.Add64(x, y, c&1)
	}
	if err := quick.CheckEqual(got, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSUB(t *testing.T) {
	proc := New()
	got := func(x, y, b uint64) (uint64, uint64) {
		s, bout := proc.SUB(Word(x), Word(y), Word(b&1))
		return uint64(s.(Word)), uint64(bout.(Word))
	}
	expect := func(x, y, b uint64) (uint64, uint64) {
		return bits.Sub64(x, y, b&1)
	}
	if err := quick.CheckEqual(got, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestMUL(t *testing.T) {
	proc := New()
	got := func(x, y uint64) (uint64, uint64) {
		hi, lo := proc.MUL(Word(x), Word(y))
		return uint64(hi.(Word)), uint64(lo.(Word))
	}
	expect := func(x, y uint64) (uint64, uint64) {
		return bits.Mul64(x, y)
	}
	if err := quick.CheckEqual(got, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSHL(t *testing.T) {
	proc := New()
	got := func(x uint64, s uint) uint64 {
		sh := proc.SHL(Word(x), s&0x3f)
		return uint64(sh.(Word))
	}
	expect := func(x uint64, s uint) uint64 {
		return x << (s & 0x3f)
	}
	if err := quick.CheckEqual(got, expect, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSHR(t *testing.T) {
	proc := New()
	got := func(x uint64, s uint) uint64 {
		sh := proc.SHR(Word(x), s&0x3f)
		return uint64(sh.(Word))
	}
	expect := func(x uint64, s uint) uint64 {
		return x >> (s & 0x3f)
	}
	if err := quick.CheckEqual(got, expect, nil); err != nil {
		t.Fatal(err)
	}
}

type BadValue int

func (BadValue) Bits() uint { return 0 }

func TestBadValueType(t *testing.T) {
	v := BadValue(42)
	proc := New()
	proc.ADD(v, v, v)
	proc.SUB(v, v, v)
	proc.MUL(v, v)
	proc.SHL(v, 7)
	proc.SHR(v, 7)
	errs := proc.Errors()
	if len(errs) != 10 {
		t.FailNow()
	}
	for _, err := range errs {
		assert.ErrorContains(t, err, "unexpected")
	}
}
