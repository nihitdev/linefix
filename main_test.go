package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime/debug"
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
	if got := stdout.String(); got != "linefix "+buildVersion()+"\n" {
		t.Fatalf("version output = %q", got)
	}
}

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name     string
		injected string
		info     *debug.BuildInfo
		ok       bool
		want     string
	}{
		{name: "release linker value", injected: "1.2.3", info: &debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}}, ok: true, want: "1.2.3"},
		{name: "go install module version", injected: "dev", info: &debug.BuildInfo{Main: debug.Module{Version: "v0.1.1"}}, ok: true, want: "0.1.1"},
		{name: "local development build", injected: "dev", info: &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, ok: true, want: "dev"},
		{name: "local VCS build", injected: "dev", info: &debug.BuildInfo{Main: debug.Module{Version: "v0.1.1-0.20260819000000-abcdef123456"}, Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abcdef123456"}}}, ok: true, want: "dev"},
		{name: "missing build info", injected: "dev", ok: false, want: "dev"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := resolveVersion(test.injected, test.info, test.ok); got != test.want {
				t.Fatalf("resolveVersion() = %q; want %q", got, test.want)
			}
		})
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
