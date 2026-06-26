package fingerprint

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// buildFakeProc creates a minimal /proc-like tree under dir and returns
// a Resolver pointed at it.
func buildFakeProc(t *testing.T, port int, inode, pid int, comm string) *Resolver {
	t.Helper()
	root := t.TempDir()

	// /proc/net/tcp
	netDir := filepath.Join(root, "net")
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hexPort := strconv.FormatInt(int64(port), 16)
	for len(hexPort) < 4 {
		hexPort = "0" + hexPort
	}
	tcpLine := "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"
	tcpLine += "   0: 00000000:" + hexPort + " 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 " + strconv.Itoa(inode) + " 1 0000000000000000 100 0 0 10 0\n"
	if err := os.WriteFile(filepath.Join(netDir, "tcp"), []byte(tcpLine), 0o644); err != nil {
		t.Fatal(err)
	}

	// /proc/<pid>/fd/<n> -> socket:[inode]
	fdDir := filepath.Join(root, strconv.Itoa(pid), "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Use a regular file as a stand-in; we override Readlink by writing a
	// symlink to the socket target string.
	socketTarget := "socket:[" + strconv.Itoa(inode) + "]"
	if err := os.Symlink(socketTarget, filepath.Join(fdDir, "3")); err != nil {
		t.Fatal(err)
	}

	// /proc/<pid>/comm
	if err := os.WriteFile(filepath.Join(root, strconv.Itoa(pid), "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	return newWithRoot(root)
}

func TestLookup_ReturnsProcessInfo(t *testing.T) {
	r := buildFakeProc(t, 8080, 99999, 1234, "myserver")
	info, err := r.Lookup(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Port != 8080 {
		t.Errorf("port: got %d, want 8080", info.Port)
	}
	if info.PID != 1234 {
		t.Errorf("pid: got %d, want 1234", info.PID)
	}
	if info.Name != "myserver" {
		t.Errorf("name: got %q, want %q", info.Name, "myserver")
	}
}

func TestLookup_UnknownPort_ReturnsError(t *testing.T) {
	r := buildFakeProc(t, 8080, 99999, 1234, "myserver")
	_, err := r.Lookup(9999)
	if err == nil {
		t.Fatal("expected error for unknown port, got nil")
	}
}

func TestInfo_String_WithName(t *testing.T) {
	i := Info{Port: 80, PID: 42, Name: "nginx"}
	got := i.String()
	want := "port 80 -> nginx (pid 42)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInfo_String_WithoutName(t *testing.T) {
	i := Info{Port: 80, PID: 42}
	got := i.String()
	want := "port 80 (pid 42)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
