package main

import "testing"

func TestIsValidCurrency(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"ASDF", false},
	}

	for _, c := range cases {
		result := isValidCurrency(c.input)
		if result != c.want {
			t.Errorf("error on currency %q ­ got %v, wanted %v\n", c.input, result, c.want)
		}
	}
}
