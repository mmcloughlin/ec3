package cost

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
)

func TestCompare(t *testing.T) {
	// Load expected costs.
	expected := map[string]string{}
	b, err := ioutil.ReadFile("testdata/expect.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatal(err)
	}

	for _, f := range efd.All {
		if f.Program == nil {
			continue
		}

		counts, err := Operations(f)
		if err != nil {
			t.Fatal(err)
		}

		expect, ok := expected[f.ID]
		if !ok {
			t.Logf("missing expected cost for %q", f.ID)
			continue
		}

		got := counts.Summary()

		if !CountsStringsEqual(expect, got) {
			t.Logf("    id: %s", f.ID)
			t.Logf("   got: %s", got)
			t.Logf("expect: %s", expect)
			t.Fail()
		}
	}
}

func CountsStringsEqual(a, b string) bool {
	pa := strings.Split(a, " + ")
	pb := strings.Split(b, " + ")

	sort.Strings(pa)
	sort.Strings(pb)

	return reflect.DeepEqual(pa, pb)
}
