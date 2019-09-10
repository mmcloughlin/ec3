// Package disjointset implements the disjoint-set data structure.
package disjointset

// DisjointSets represents the set of string keys partitioned into disjoint sets.
type DisjointSets struct {
	parent map[string]string
}

// New builds a disjoint set structure. Initially every key is in its own
// singleton set.
func New() *DisjointSets {
	return &DisjointSets{
		parent: make(map[string]string),
	}
}

// Find returns the key of the set k belongs to.
func (d *DisjointSets) Find(k string) string {
	p, ok := d.parent[k]
	if !ok || p == k {
		return k
	}
	return d.Find(p)
}

// Same reports whether k and l are in the same set.
func (d *DisjointSets) Same(k, l string) bool {
	return d.Find(k) == d.Find(l)
}

// Union the sets k and l are in.
func (d *DisjointSets) Union(k, l string) {
	sk := d.Find(k)
	sl := d.Find(l)

	// If they are already in the same set, nothing to be done.
	if sk == sl {
		return
	}

	// Point one at the other.
	d.parent[sk] = sl
}
