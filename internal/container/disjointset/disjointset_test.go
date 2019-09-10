package disjointset

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	d := New()
	for _, k := range []string{"hello", "world", "testing"} {
		if d.Find(k) != k {
			t.FailNow()
		}
	}
}

func TestUnion(t *testing.T) {
	// Construct n elements in m sets.
	m, n := 7, 100

	// Make assignments.
	set := map[string]string{}
	members := map[string][]string{}

	for i := 0; i < m; i++ {
		k := key(i)
		set[k] = k
		members[k] = []string{k}
	}

	for i := m; i < n; i++ {
		k := key(i)
		s := key(rand.Intn(m))
		set[k] = s
		members[s] = append(members[s], k)
	}

	// Build set that should have the same assignments.
	d := New()
	for s := range members {
		size := len(members[s])
		rand.Shuffle(size, func(i, j int) {
			members[s][i], members[s][j] = members[s][j], members[s][i]
		})
		for i := 0; i < size; i++ {
			d.Union(members[s][i], members[s][(i+1)%size])
		}
	}

	// Confirm same set membership.
	for k, expect := range set {
		if !d.Same(k, expect) {
			t.FailNow()
		}
	}
}

func key(i int) string {
	return strconv.Itoa(i)
}
