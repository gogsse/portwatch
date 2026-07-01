package portmap

import "testing"

func TestNew_ContainsBuiltIns(t *testing.T) {
	m := New()
	want := map[int]string{
		22:    "ssh",
		80:    "http",
		443:   "https",
		3306:  "mysql",
		6379:  "redis",
		27017: "mongodb",
	}
	for port, svc := range want {
		e, ok := m.Lookup(port)
		if !ok {
			t.Errorf("port %d not found in built-in map", port)
			continue
		}
		if e.Service != svc {
			t.Errorf("port %d: got service %q, want %q", port, e.Service, svc)
		}
	}
}

func TestLookup_UnknownPort_ReturnsFalse(t *testing.T) {
	m := New()
	_, ok := m.Lookup(99999)
	if ok {
		t.Error("expected false for unknown port, got true")
	}
}

func TestAdd_OverridesEntry(t *testing.T) {
	m := New()
	m.Add(Entry{Port: 22, Service: "custom-ssh", Proto: "tcp"})
	e, ok := m.Lookup(22)
	if !ok {
		t.Fatal("port 22 not found after Add")
	}
	if e.Service != "custom-ssh" {
		t.Errorf("got %q, want %q", e.Service, "custom-ssh")
	}
}

func TestAdd_NewCustomPort(t *testing.T) {
	m := New()
	m.Add(Entry{Port: 9200, Service: "elasticsearch", Proto: "tcp"})
	e, ok := m.Lookup(9200)
	if !ok {
		t.Fatal("custom port 9200 not found")
	}
	if e.Service != "elasticsearch" {
		t.Errorf("got %q, want %q", e.Service, "elasticsearch")
	}
}

func TestLabel_KnownPort(t *testing.T) {
	m := New()
	got := m.Label(22)
	want := "ssh (22)"
	if got != want {
		t.Errorf("Label(22) = %q, want %q", got, want)
	}
}

func TestLabel_UnknownPort(t *testing.T) {
	m := New()
	got := m.Label(65000)
	want := "unknown (65000)"
	if got != want {
		t.Errorf("Label(65000) = %q, want %q", got, want)
	}
}

func TestKnown_TrueForBuiltIn(t *testing.T) {
	m := New()
	if !m.Known(443) {
		t.Error("Known(443) returned false, want true")
	}
}

func TestKnown_FalseForUnknown(t *testing.T) {
	m := New()
	if m.Known(55555) {
		t.Error("Known(55555) returned true, want false")
	}
}
