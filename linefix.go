package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LineEnding is a supported output line-ending format.
type LineEnding string

const (
	EndingLF   LineEnding = "lf"
	EndingCRLF LineEnding = "crlf"
)

// ErrBinaryFile is returned when input appears to contain binary data.
var ErrBinaryFile = errors.New("binary data detected")

// Detect reports the line-ending style used by data. Bare carriage returns are
// not line endings; a file must contain both LF and CRLF endings to be mixed.
func Detect(data []byte) string {
	var lf, crlf bool
	for i, b := range data {
		if b != '\n' {
			continue
		}
		if i > 0 && data[i-1] == '\r' {
			crlf = true
		} else {
			lf = true
		}
	}
	switch {
	case lf && crlf:
		return "Mixed"
	case crlf:
		return "CRLF"
	case lf:
		return "LF"
	default:
		return "No line endings"
	}
}

// Convert returns data with the requested line endings and whether it changed.
func Convert(data []byte, ending LineEnding) ([]byte, bool, error) {
	if isLikelyBinary(data) {
		return nil, false, ErrBinaryFile
	}
	var converted []byte
	switch ending {
	case EndingLF:
		converted = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	case EndingCRLF:
		// Normalize first; directly replacing every LF would turn CRLF into CRCRLF.
		normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
		converted = bytes.ReplaceAll(normalized, []byte("\n"), []byte("\r\n"))
	default:
		return nil, false, fmt.Errorf("unsupported line ending %q", ending)
	}
	if bytes.Equal(data, converted) {
		return data, false, nil
	}
	return converted, true, nil
}

// CheckFile detects the endings in a text file without modifying it.
func CheckFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", path, err)
	}
	if isLikelyBinary(data) {
		return "", fmt.Errorf("%q: %w", path, ErrBinaryFile)
	}
	return Detect(data), nil
}

// ConvertFile converts a text file in place. An unchanged file is not rewritten.
func ConvertFile(path string, ending LineEnding) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("stat %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return false, fmt.Errorf("%q is not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("read %q: %w", path, err)
	}
	converted, changed, err := Convert(data, ending)
	if err != nil {
		return false, fmt.Errorf("%q: %w", path, err)
	}
	if !changed {
		return false, nil
	}
	if err := replaceFile(path, converted, info.Mode().Perm()); err != nil {
		return false, err
	}
	return true, nil
}

func replaceFile(path string, data []byte, mode os.FileMode) (err error) {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".linefix-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %q: %w", path, err)
	}
	tmpName := tmp.Name()
	defer func() {
		if tmp != nil {
			_ = tmp.Close()
		}
		if err != nil {
			_ = os.Remove(tmpName)
		}
	}()

	if err = tmp.Chmod(mode); err != nil {
		return fmt.Errorf("set permissions on temporary file: %w", err)
	}
	if _, err = tmp.Write(data); err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err = tmp.Sync(); err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		tmp = nil
		return fmt.Errorf("close temporary file: %w", err)
	}
	tmp = nil
	if err = os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %q: %w", path, err)
	}
	return nil
}

func isLikelyBinary(data []byte) bool {
	const sampleSize = 8192
	if len(data) > sampleSize {
		data = data[:sampleSize]
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return true
	}
	if len(data) == 0 {
		return false
	}
	controls := 0
	for _, b := range data {
		if b < 0x20 && b != '\n' && b != '\r' && b != '\t' && b != '\f' && b != '\b' {
			controls++
		}
	}
	return controls*10 > len(data)
}
