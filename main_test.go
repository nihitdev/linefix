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
		if !strings.Contains(stdout.String(), "Usage:") {
			t.Errorf("run(%q) did not print usage", arg)
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

func TestRunMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	lf := filepath.Join(dir, "lf.txt")
	crlf := filepath.Join(dir, "crlf.txt")
	if err := os.WriteFile(lf, []byte("one\ntwo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(crlf, []byte("one\r\ntwo\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"check", lf, crlf}, &stdout, &stderr); code != 0 {
		t.Fatalf("run(check multiple) exit code = %d; stderr = %q", code, stderr.String())
	}
	want := lf + ": LF\n" + crlf + ": CRLF\n"
	if got := stdout.String(); got != want {
		t.Fatalf("check output = %q; want %q", got, want)
	}
}

func TestRunDryRunAndQuiet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	original := []byte("one\r\ntwo\r\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"lf", "--dry-run", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("dry run exit code = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "would convert to LF") {
		t.Fatalf("dry-run output = %q", stdout.String())
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, original) {
		t.Fatalf("dry run modified file: data=%q error=%v", got, err)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"--quiet", "lf", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("quiet conversion exit code = %d; stderr = %q", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("quiet output = %q; want none", stdout.String())
	}
}

func TestRunRejectsInvalidArguments(t *testing.T) {
	tests := [][]string{{}, {"unknown", "file"}, {"lf"}, {"--version", "extra"}, {"--wat", "lf", "file"}, {"--help", "lf"}, {"--dry-run", "check", "file"}}
	for _, args := range tests {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code == 0 {
			t.Errorf("run(%q) unexpectedly succeeded", args)
		}
	}
}
