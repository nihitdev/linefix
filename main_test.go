package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelpAndVersion(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		var stdout, stderr bytes.Buffer
		if code := run([]string{arg}, &stdout, &stderr); code != 0 {
			t.Errorf("run(%q) exit code = %d; want 0", arg, code)
		}
		if !strings.Contains(stderr.String(), "Usage:") {
			t.Errorf("run(%q) did not print usage", arg)
		}
		for _, section := range []string{"Options:", "Check output:", "Examples:", "man linefix"} {
			if !strings.Contains(stderr.String(), section) {
				t.Errorf("run(%q) help missing %q", arg, section)
			}
		}
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("run(--version) exit code = %d; want 0", code)
	}
	if got := stdout.String(); got != "linefix "+version+"\n" {
		t.Fatalf("version output = %q", got)
	}
}

func TestRunCheck(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mixed.txt")
	if err := os.WriteFile(path, []byte("one\r\ntwo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{"check", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("run(check) exit code = %d; stderr = %q", code, stderr.String())
	}
	if got := stdout.String(); got != "Mixed\n" {
		t.Fatalf("check output = %q; want Mixed", got)
	}
}

func TestRunRejectsInvalidArguments(t *testing.T) {
	tests := [][]string{{}, {"unknown", "file"}, {"lf"}, {"--version", "extra"}}
	for _, args := range tests {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code == 0 {
			t.Errorf("run(%q) unexpectedly succeeded", args)
		}
	}
}
