package main

import (
	"testing"
)

func TestValueToLiteral(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"ab c", "`ab c`"},
		{"ab\nc", "`ab\nc`"},
		{"ab`\nc", "\"ab`\\nc\""},
		{"\x00", "\"\\x00\""},
	} {
		value := ValueToLiteral(test.in)
		if test.out != value {
			t.Errorf("ValueToLiteral(%+v): %+v != %+v", test.in, value, test.out)
		}
	}
}
