package acc

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/ec3/prime"
)

func TestEvaluate(t *testing.T) {
	cases := []struct {
		Name   string
		Lines  []string
		Expect *big.Int
	}{
		{
			Name: "add",
			Lines: []string{
				"1 add 1",
			},
			Expect: big.NewInt(2),
		},
		{
			Name: "double",
			Lines: []string{
				"dbl 1",
			},
			Expect: big.NewInt(2),
		},
		{
			Name: "shl",
			Lines: []string{
				"1 shl 3",
			},
			Expect: big.NewInt(8),
		},
		{
			Name: "p25519-2",
			Lines: []string{
				"_1    = 1",
				"_10   = _1 << 1",
				"_1001 = _10 << 2 + _1",
				"_1011 = _1001 + _10",
				"x5    = _1011 << 1 + _1001",
				"x10   = x5 << 5 + x5",
				"x20   = x10 << 10 + x10",
				"x40   = x20 << 20 + x20",
				"x50   = x40 << 10 + x10",
				"x100  = x50 << 50 + x50",
				"x200  = x100 << 100 + x100",
				"x250  = x200 << 50 + x50",
				"return x250 << 5 + _1011",
			},
			Expect: new(big.Int).Sub(
				prime.P25519.Int(),
				big.NewInt(2),
			),
		},
	}
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			var src string
			for _, line := range c.Lines {
				src += line + "\n"
			}

			// Parse to AST.
			s, err := String(src)
			if err != nil {
				t.Fatal(err)
			}

			// Translate to IR.
			p, err := Translate(s)
			if err != nil {
				t.Fatal(err)
			}

			// Compile to Program.
			prog, err := Compile(p)
			if err != nil {
				t.Fatal(err)
			}

			// Evaluate chain.
			ch := prog.Evaluate()

			// Report.
			for i, x := range ch {
				t.Logf("[%3d] bin=%b", i, x)
			}

			if err := ch.Produces(c.Expect); err != nil {
				t.Fatalf("chain does not produce %d: %s", c.Expect, err)
			}
		})
	}
}
