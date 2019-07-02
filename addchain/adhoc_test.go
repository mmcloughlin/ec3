package addchain

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/mmcloughlin/ec3/prime"
)

func TestAdhoc(t *testing.T) {
	a := NewDictAlgorithm(
		RunLength{T: 0},
		NewContinuedFractions(DichotomicStrategy{}),
	)
	p := prime.NISTP256.Int()
	q := new(big.Int).Sub(p, big.NewInt(3))
	c, err := a.FindChain(q)
	if err != nil {
		t.Fatal(err)
	}

	for i, x := range c {
		t.Logf("%3d: %x", i+1, x)
	}

	if err := c.Produces(q); err != nil {
		t.Fatal(err)
	}
}

func TestStuck(t *testing.T) {
	// opt(dictionary(run_length(0),continued_fractions(co_binary)))
	a := NewDictAlgorithm(
		RunLength{T: 0},
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
	)

	p := prime.P25519.Int()
	q := new(big.Int).Sub(p, big.NewInt(2))

	c, err := a.FindChain(q)
	if err != nil {
		t.Fatal(err)
	}

	for i, x := range c {
		t.Logf("%3d: %x", i+1, x)
	}

	if err := c.Produces(q); err != nil {
		t.Fatal(err)
	}
}

func TestEnsembleSerial(t *testing.T) {
	p := prime.P25519.Int()
	q := new(big.Int).Sub(p, big.NewInt(2))

	as := Ensemble()

	for _, a := range as {
		fmt.Println(a)
		_, err := a.FindChain(q)
		if err != nil {
			t.Fatal(err)
		}
	}
}
