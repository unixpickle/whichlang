package tokens

import "testing"

func TestCountTokens(t *testing.T) {
	document := "Hello this is a Hello123\ttest\nhello hi is1 is123"
	actual := CountTokens(document)
	expected := map[string]int{
		"Hello123": 1,
		"is1":      1,
		"is123":    1,
		"Hello":    2,
		"this":     1,
		"is":       3,
		"a":        1,
		"123":      2,
		"test":     1,
		"hello":    1,
		"hi":       1,
		"1":        1,
	}

	for x, count := range expected {
		if actual[x] != count {
			t.Error("expected count", count, "for", x, "but got", actual[x])
		}
	}

	for x := range actual {
		if expected[x] == 0 {
			t.Error("got unexpected token:", x)
		}
	}
}
