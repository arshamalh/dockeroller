package tools

import "testing"

func TestIndexer(t *testing.T) {
	cases := []struct {
		num    int
		length int
		want   int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{-1, 10, 9},
		{-3, 10, 7},
	}
	for _, c := range cases {
		if got := Indexer(c.num, c.length); got != c.want {
			t.Error("it cannot be done.")
		}
	}
}
