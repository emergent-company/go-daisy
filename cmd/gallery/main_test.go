package main

import (
	"os"
	"testing"
)

func TestParseBool(t *testing.T) {
	const envName = "GO_DAISY_TEST_BOOL"

	cases := []struct {
		name  string
		set   bool
		value string
		def   bool
		want  bool
	}{
		{"unset defaults true", false, "", true, true},
		{"unset defaults false", false, "", false, false},
		{"empty defaults true", true, "", true, true},
		{"true", true, "true", false, true},
		{"TRUE", true, "TRUE", false, true},
		{"1", true, "1", false, true},
		{"t", true, "t", false, true},
		{"false", true, "false", true, false},
		{"0", true, "0", true, false},
		{"F", true, "F", true, false},
		{"padded true", true, "  true  ", false, true},
		{"invalid defaults true", true, "notabool", true, true},
		{"invalid defaults false", true, "notabool", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.set {
				t.Setenv(envName, tc.value)
			} else {
				_ = os.Unsetenv(envName)
			}
			if got := parseBool(envName, tc.def); got != tc.want {
				t.Errorf("parseBool(%q, %v) with env %q = %v, want %v",
					envName, tc.def, tc.value, got, tc.want)
			}
		})
	}
}
