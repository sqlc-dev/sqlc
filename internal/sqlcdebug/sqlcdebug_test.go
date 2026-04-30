package sqlcdebug

import (
	"sync"
	"testing"
)

func resetForTest(t *testing.T) {
	t.Helper()
	registryMu.Lock()
	registry = map[string]*Setting{}
	envMap = nil
	envOnce = sync.Once{}
	registryMu.Unlock()
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  map[string]string
	}{
		{"", map[string]string{}},
		{"dumpast=1", map[string]string{"dumpast": "1"}},
		{"dumpast=1,trace=trace.out", map[string]string{"dumpast": "1", "trace": "trace.out"}},
		{"  dumpast=1 , processplugins=0 ", map[string]string{"dumpast": "1", "processplugins": "0"}},
		{"trace=", map[string]string{"trace": ""}},
		{"bare", map[string]string{}},
	}
	for _, tt := range tests {
		got := parse(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("parse(%q): got %v, want %v", tt.input, got, tt.want)
			continue
		}
		for k, v := range tt.want {
			if got[k] != v {
				t.Errorf("parse(%q)[%q] = %q, want %q", tt.input, k, got[k], v)
			}
		}
	}
}

func TestSettingValue(t *testing.T) {
	resetForTest(t)
	Update("dumpast=1,trace=foo.out")

	if v := New("dumpast").Value(); v != "1" {
		t.Errorf("dumpast = %q, want %q", v, "1")
	}
	if v := New("trace").Value(); v != "foo.out" {
		t.Errorf("trace = %q, want %q", v, "foo.out")
	}
	if !New("dumpast").IsSet() {
		t.Errorf("IsSet(dumpast) = false, want true")
	}

	// Unset key returns its registered default.
	if v := New("processplugins").Value(); v != "1" {
		t.Errorf("processplugins default = %q, want %q", v, "1")
	}
	if New("processplugins").IsSet() {
		t.Errorf("IsSet(processplugins) = true, want false")
	}
}

func TestUpdateRefreshesExistingSettings(t *testing.T) {
	resetForTest(t)
	dumpAST := New("dumpast")

	if v := dumpAST.Value(); v != "" {
		t.Errorf("initial dumpast = %q, want empty", v)
	}

	Update("dumpast=1")
	if v := dumpAST.Value(); v != "1" {
		t.Errorf("after Update dumpast = %q, want %q", v, "1")
	}

	Update("")
	if v := dumpAST.Value(); v != "" {
		t.Errorf("after reset dumpast = %q, want empty", v)
	}
	if dumpAST.IsSet() {
		t.Errorf("after reset IsSet = true, want false")
	}
}

func TestAny(t *testing.T) {
	resetForTest(t)
	if Any() {
		t.Errorf("Any() = true, want false on empty env")
	}
	Update("trace=1")
	if !Any() {
		t.Errorf("Any() = false, want true after Update")
	}
}

func TestNewPanicsOnUnknown(t *testing.T) {
	resetForTest(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("New(\"bogus\") did not panic")
		}
	}()
	New("bogus")
}
