// linefix converts text-file line endings in place.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

// version is replaced for releases with: -ldflags "-X main.version=<version>".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("linefix", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { printUsage(stderr) }
	showVersion := fs.Bool("version", false, "print version")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if *showVersion {
		if fs.NArg() != 0 {
			fmt.Fprintln(stderr, "linefix: --version does not accept arguments")
			return 2
		}
		fmt.Fprintf(stdout, "linefix %s\n", buildVersion())
		return 0
	}
	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "linefix: expected a command and a file")
		fmt.Fprintln(stderr, "Try 'linefix --help' for usage.")
		return 2
	}

	command, path := fs.Arg(0), fs.Arg(1)
	switch command {
	case "check":
		ending, err := CheckFile(path)
		if err != nil {
			printError(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, ending)
		return 0
	case "lf", "crlf":
		changed, err := ConvertFile(path, LineEnding(command))
		if err != nil {
			printError(stderr, err)
			return 1
		}
		if changed {
			fmt.Fprintf(stdout, "%s: converted to %s\n", path, displayEnding(LineEnding(command)))
		} else {
			fmt.Fprintf(stdout, "%s: already %s\n", path, displayEnding(LineEnding(command)))
		}
		return 0
	default:
		fmt.Fprintf(stderr, "linefix: unknown command %q\n", command)
		fmt.Fprintln(stderr, "Try 'linefix --help' for usage.")
		return 2
	}
}

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	return resolveVersion(version, info, ok)
}

func resolveVersion(injected string, info *debug.BuildInfo, ok bool) string {
	if injected != "dev" {
		return injected
	}
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return injected
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return injected
		}
	}
	return strings.TrimPrefix(info.Main.Version, "v")
}

func printError(w io.Writer, err error) {
	if errors.Is(err, ErrBinaryFile) {
		fmt.Fprintf(w, "linefix: refusing to process likely binary file: %v\n", err)
		return
	}
	fmt.Fprintf(w, "linefix: %v\n", err)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `linefix — convert text-file line endings safely

Usage:
  linefix lf <file>       Convert CRLF line endings to LF
  linefix crlf <file>     Convert LF line endings to CRLF
  linefix check <file>    Detect the file's line endings

Options:
  -h, --help              Show this help and exit
      --version           Show the version and exit

Check output:
  LF                      All line endings are LF
  CRLF                    All line endings are CRLF
  Mixed                   Both LF and CRLF are present
  No line endings         No LF or CRLF endings are present

Examples:
  linefix check README.md
  linefix lf script.sh
  linefix crlf notes.txt

Files are modified in place. Existing permissions and trailing-newline
state are preserved. Likely binary files are refused.

Manual: man linefix
Source: https://github.com/nihitdev/linefix`)
}

func displayEnding(ending LineEnding) string {
	if ending == EndingCRLF {
		return "CRLF"
	}
	return "LF"
}
