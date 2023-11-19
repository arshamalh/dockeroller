package tools

import "testing"

func TestSizeToHumanReadable(t *testing.T) {
	cases := []struct {
		size int64
		want string
	}{
		{10, "10 Bytes"},
		{1000, "1.00 KiloBytes"},
		{10001, "10.00 KiloBytes"},
		{1000000, "1.00 MegaBytes"},
		{10000000000, "10.00 GigaBytes"},
	}
	for _, c := range cases {
		if got := SizeToHumanReadable(c.size); got != c.want {
			t.Error("it cannot be done.")
		}
	}
}
