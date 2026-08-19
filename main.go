// linefix converts text-file line endings in place.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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
		fmt.Fprintf(stdout, "linefix %s\n", version)
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

func printError(w io.Writer, err error) {
	if errors.Is(err, ErrBinaryFile) {
		fmt.Fprintf(w, "linefix: refusing to process likely binary file: %v\n", err)
		return
	}
	fmt.Fprintf(w, "linefix: %v\n", err)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `Usage:
  linefix lf <file>       Convert CRLF line endings to LF
  linefix crlf <file>     Convert LF line endings to CRLF
  linefix check <file>    Report LF, CRLF, Mixed, or No line endings
  linefix --version       Print the version
  linefix -h | --help     Show this help`)
}

func displayEnding(ending LineEnding) string {
	if ending == EndingCRLF {
		return "CRLF"
	}
	return "LF"
}
