package alpine

import (
	"encoding/json"
	"testing"
)

// TestPaletteStateDeclaresPicked locks the property name the Palette component
// binds with x-model. If this ever renames, Palette needs the same rename or it
// regresses to "picked is not defined" at Alpine init.
func TestPaletteStateDeclaresPicked(t *testing.T) {
	cases := []struct {
		name   string
		picked string
		want   string
	}{
		{name: "empty", picked: "", want: `{"picked":""}`},
		{name: "tailwind_shade", picked: "red-500", want: `{"picked":"red-500"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(PaletteState(tc.picked))
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if got := string(b); got != tc.want {
				t.Errorf("PaletteState(%q) = %s, want %s", tc.picked, got, tc.want)
			}
		})
	}
}
