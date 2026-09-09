package alpine

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestComboboxStateMarshalsEmptySelected(t *testing.T) {
	b, err := json.Marshal(ComboboxState(nil, nil))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if strings.Contains(got, `"selected":null`) {
		t.Errorf("ComboboxState(nil) = %s, want [] not null (null breaks selected.length/indexOf)", got)
	}
	if !strings.Contains(got, `"selected":[]`) {
		t.Errorf("ComboboxState(nil) = %s, want selected:[]", got)
	}
}

func TestComboboxStateIncludesLabels(t *testing.T) {
	b, err := json.Marshal(ComboboxState(
		[]string{"tech"},
		map[string]string{"tech": "Technology", "health": "Healthcare"},
	))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, `"selected":["tech"]`) {
		t.Errorf("ComboboxState() = %s, want selected:[\"tech\"]", got)
	}
	if !strings.Contains(got, `"labels":{"health":"Healthcare","tech":"Technology"}`) {
		t.Errorf("ComboboxState() = %s, want labels map", got)
	}
}

func TestComboboxInitUsesDataNotThis(t *testing.T) {
	for _, multi := range []bool{false, true} {
		expr := ComboboxInit(multi)
		if strings.HasPrefix(expr, "var ") {
			t.Errorf("ComboboxInit(%v) starts with 'var' — Alpine rejects it (Unexpected token 'var'); use 'let'", multi)
		}
		if !strings.Contains(expr, "root=$data") {
			t.Errorf("ComboboxInit(%v) must grab $data so Alpine.$data() finds the methods: %s", multi, expr)
		}
		if strings.Contains(expr, "root=this") {
			t.Errorf("ComboboxInit(%v) binds root=this (window); use root=$data: %s", multi, expr)
		}
		if multi && !strings.Contains(expr, "root.selected.splice") {
			t.Errorf("ComboboxInit(true) missing toggle-select splice: %s", expr)
		}
	}
}

func TestStructuredInputInitUsesDataNotThis(t *testing.T) {
	expr := StructuredInputInit()
	if !strings.Contains(expr, "root=$data") || strings.Contains(expr, "root=this") {
		t.Errorf("StructuredInputInit must use root=$data (not root=this): %s", expr)
	}
}

func TestTagListInitUsesDataNotThis(t *testing.T) {
	expr := TagListInit()
	if !strings.Contains(expr, "root=$data") || strings.Contains(expr, "root=this") {
		t.Errorf("TagListInit must use root=$data (not root=this): %s", expr)
	}
}
