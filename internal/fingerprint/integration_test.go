package fingerprint_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stevezaluk/portwatch/internal/fingerprint"
)

// buildExtFakeProc mirrors buildFakeProc but is exported for integration tests.
func buildExtFakeProc(t *testing.T, port, inode, pid int, comm string) *fingerprint.Resolver {
	t.Helper()
	root := t.TempDir()

	netDir := filepath.Join(root, "net")
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hexPort := strconv.FormatInt(int64(port), 16)
	for len(hexPort) < 4 {
		hexPort = "0" + hexPort
	}
	tcpContent := "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
		"   0: 00000000:" + hexPort + " 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 " +
		strconv.Itoa(inode) + " 1 0000000000000000 100 0 0 10 0\n"
	if err := os.WriteFile(filepath.Join(netDir, "tcp"), []byte(tcpContent), 0o644); err != nil {
		t.Fatal(err)
	}

	fdDir := filepath.Join(root, strconv.Itoa(pid), "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	socketTarget := "socket:[" + strconv.Itoa(inode) + "]"
	if err := os.Symlink(socketTarget, filepath.Join(fdDir, "5")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, strconv.Itoa(pid), "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	return fingerprint.NewWithRoot(root)
}

func TestFingerprint_FullCycle(t *testing.T) {
	ports := []struct {
		port  int
		inode int
		pid   int
		comm  string
	}{
		{3000, 11111, 501, "node"},
		{5432, 22222, 502, "postgres"},
	}
	for _, tc := range ports {
		r := buildExtFakeProc(t, tc.port, tc.inode, tc.pid, tc.comm)
		info, err := r.Lookup(tc.port)
		if err != nil {
			t.Errorf("port %d: unexpected error: %v", tc.port, err)
			continue
		}
		if info.Name != tc.comm {
			t.Errorf("port %d: name got %q want %q", tc.port, info.Name, tc.comm)
		}
		if info.PID != tc.pid {
			t.Errorf("port %d: pid got %d want %d", tc.port, info.PID, tc.pid)
		}
	}
}
