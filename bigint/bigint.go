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

// InsertSortedUnique inserts an integer into a slice of sorted distinct
// integers.
func InsertSortedUnique(xs []*big.Int, x *big.Int) []*big.Int {
	return MergeUnique([]*big.Int{x}, xs)
}

// MergeUnique merges two slices of sorted distinct integers. Elements in both
// slices are deduplicated.
func MergeUnique(xs, ys []*big.Int) []*big.Int {
	r := make([]*big.Int, 0, len(xs)+len(ys))

	for len(xs) > 0 && len(ys) > 0 {
		switch xs[0].Cmp(ys[0]) {
		case -1:
			r = append(r, xs[0])
			xs = xs[1:]
		case 0:
			r = append(r, xs[0])
			xs = xs[1:]
			ys = ys[1:]
		case 1:
			r = append(r, ys[0])
			ys = ys[1:]
		}
	}

	r = append(r, xs...)
	r = append(r, ys...)

	return r
}
