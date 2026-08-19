package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		ending  LineEnding
		want    string
		changed bool
	}{
		{"CRLF to LF", "one\r\ntwo\r\n", EndingLF, "one\ntwo\n", true},
		{"LF to CRLF", "one\ntwo\n", EndingCRLF, "one\r\ntwo\r\n", true},
		{"already LF", "one\ntwo\n", EndingLF, "one\ntwo\n", false},
		{"already CRLF", "one\r\ntwo\r\n", EndingCRLF, "one\r\ntwo\r\n", false},
		{"mixed to LF", "one\r\ntwo\nthree", EndingLF, "one\ntwo\nthree", true},
		{"empty", "", EndingLF, "", false},
		{"no newline", "one", EndingCRLF, "one", false},
		{"final newline", "one\n", EndingCRLF, "one\r\n", true},
		{"prevent doubled CR", "one\r\ntwo\n", EndingCRLF, "one\r\ntwo\r\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed, err := Convert([]byte(tt.input), tt.ending)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want || changed != tt.changed {
				t.Fatalf("Convert() = %q, %v; want %q, %v", got, changed, tt.want, tt.changed)
			}
			if strings.Contains(string(got), "\r\r\n") {
				t.Fatal("conversion created CRCRLF")
			}
		})
	}
}

func TestDetect(t *testing.T) {
	tests := map[string]string{
		"one\ntwo\n":     "LF",
		"one\r\ntwo\r\n": "CRLF",
		"one\r\ntwo\n":   "Mixed",
		"":               "No line endings",
		"one":            "No line endings",
	}
	for input, want := range tests {
		if got := Detect([]byte(input)); got != want {
			t.Errorf("Detect(%q) = %q; want %q", input, got, want)
		}
	}
}

func TestBinaryFileRejected(t *testing.T) {
	data := []byte{'t', 'e', 'x', 't', 0, '\n'}
	if _, _, err := Convert(data, EndingLF); !errors.Is(err, ErrBinaryFile) {
		t.Fatalf("Convert() error = %v; want ErrBinaryFile", err)
	}
}

func TestConvertFilePreservesPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("one\r\ntwo\r\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ConvertFile(path, EndingLF); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := after.Mode().Perm(), before.Mode().Perm(); got != want {
		t.Fatalf("permissions = %o; want original permissions %o", got, want)
	}
}
