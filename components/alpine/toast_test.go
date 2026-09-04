package alpine

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToastQueueStateMarshalsEmptyArray(t *testing.T) {
	b, err := json.Marshal(ToastQueueState())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(b); got != `{"toasts":[]}` {
		t.Errorf("ToastQueueState() = %s, want {\"toasts\":[]} (null breaks Alpine push)", got)
	}
}

func TestToastQueueStateSeedsInitialToasts(t *testing.T) {
	b, err := json.Marshal(ToastQueueState(
		ToastItem{Type: "success", Message: "Saved", Action: "Undo"},
		ToastItem{Type: "error", Message: "Failed", Duration: 6000},
	))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"toasts":[{"id":"seed-0","type":"success","message":"Saved","action":"Undo","duration":0},{"id":"seed-1","type":"error","message":"Failed","duration":6000}]}`
	if got := string(b); got != want {
		t.Errorf("ToastQueueState(seed...) = %s\nwant %s", got, want)
	}
}

func TestToastQueueInitIsValidAlpineExpression(t *testing.T) {
	expr := ToastQueueInit()
	if strings.HasPrefix(expr, "var ") {
		t.Errorf("ToastQueueInit starts with 'var' — Alpine rejects it (Unexpected token 'var'); use 'let'")
	}
	if !strings.Contains(expr, "queue.add") || !strings.Contains(expr, "queue.dismiss") {
		t.Errorf("ToastQueueInit missing add/dismiss: %s", expr)
	}
	if !strings.Contains(expr, "queue=$data") {
		t.Errorf("ToastQueueInit must grab the data object via $data so Alpine.$data() finds the methods: %s", expr)
	}
}
