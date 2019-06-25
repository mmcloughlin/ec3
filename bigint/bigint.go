package bigint

import (
	"math/big"
	"sort"
)

// Zero returns 0.
func Zero() *big.Int {
	return big.NewInt(0)
}

// TrailingZeros returns the number of trailing zero bits in x. Returns 0 if x is 0.
func TrailingZeros(x *big.Int) int {
	if x.BitLen() == 0 {
		return 0
	}
	n := 0
	for ; x.Bit(n) == 0; n++ {
	}
	return n
}

// Sort in ascending order.
func Sort(xs []*big.Int) {
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].Cmp(xs[j]) < 0
	})
}
