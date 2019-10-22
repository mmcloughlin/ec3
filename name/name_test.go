package name

import "testing"

func TestSchemes(t *testing.T) {
	cases := []struct {
		Name   string
		Scheme Scheme
		Expect string
	}{
		{"name", Plain, "name"},
		{"Name", Plain, "Name"},

		{"nAmE", LowerCase, "name"},
		{"name", LowerCase, "name"},

		{"name", Prefixed("prefix"), "prefixname"},
		{"Name", Prefixed("prefix"), "PrefixName"},

		{"Name", CompositeScheme(Prefixed("prefix"), LowerCase), "prefixname"},
	}
	for _, c := range cases {
		if got := c.Scheme.Name(c.Name); got != c.Expect {
			t.FailNow()
		}
	}
}

func TestSetExported(t *testing.T) {
	cases := []struct {
		Name     string
		Exported bool
		Expect   string
	}{
		{"hello", false, "hello"},
		{"hello", true, "Hello"},
		{"Hello", false, "hello"},
		{"Hello", true, "Hello"},
		{"", false, ""},
		{"", true, ""},
	}
	for _, c := range cases {
		if got := SetExported(c.Name, c.Exported); got != c.Expect {
			t.Errorf("SetExported(%v,%v) = %v; expect %v", c.Name, c.Exported, got, c.Expect)
		}
	}
}
