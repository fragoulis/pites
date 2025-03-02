package autofix_test

import (
	"testing"
)

func TestAddressCity(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		ok   bool
	}{
		{
			name: "ok",
			in:   "αρχιμηδους 18",
			out:  "asdf",
			ok:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Skip()
		})
	}
}
