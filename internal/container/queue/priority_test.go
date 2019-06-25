package queue

import "testing"

func TestPriority(t *testing.T) {
	q := NewPriority()

	q.Insert(1, 1)
	q.Insert(13, 10000)
	q.Insert(42, 100)

	for _, expect := range []int{1, 42, 13} {
		if q.Empty() {
			t.Fail()
		}
		if q.Peek() != expect {
			t.Fail()
		}
		if q.Pop() != expect {
			t.Fail()
		}
	}

	if !q.Empty() {
		t.Fail()
	}
}
