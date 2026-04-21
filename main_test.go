//go:build integration

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEnvSuppressedBumpIsNotSilent verifies that when BUMP_NO_ALPHA suppresses
// a -alpha bump, the CLI exits non-zero or prints a clear warning rather than
// silently writing nothing and exiting 0. This catches the class of bug where
// a user's environment silently neuters bump operations.
func TestEnvSuppressedBumpIsNotSilent(t *testing.T) {
	dir := t.TempDir()
	versionFile := filepath.Join(dir, "VERSION")
	if err := os.WriteFile(versionFile, []byte("v1.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".", "-in", versionFile, "-alpha", "-write")
	cmd.Env = append(os.Environ(),
		"BUMP_NO_ALPHA=true",
	)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	// When BUMP_NO_ALPHA suppresses the only requested bump operation,
	// the binary must NOT silently exit 0 having written nothing.
	// It should either exit non-zero, or emit a warning that the operation
	// was suppressed by environment configuration.
	if err == nil {
		// Exited 0 — only acceptable if output contains a suppression warning
		if !strings.Contains(output, "suppress") &&
			!strings.Contains(output, "no bump") &&
			!strings.Contains(output, "BUMP_NO_ALPHA") &&
			!strings.Contains(output, "No bump operation") {
			t.Errorf("bump exited 0 with BUMP_NO_ALPHA=true but gave no suppression warning.\nOutput: %s", output)
		}
	}

	// Either way, the file must be unchanged
	raw, _ := os.ReadFile(versionFile)
	if strings.TrimSpace(string(raw)) != "v1.0.0" {
		t.Errorf("file should be unchanged when bump is suppressed, got %q", strings.TrimSpace(string(raw)))
	}
}
