// linefix converts text-file line endings in place.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// version is replaced for releases with: -ldflags "-X main.version=<version>".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	options, positional, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "linefix: %v\n", err)
		fmt.Fprintln(stderr, "Try 'linefix --help' for usage.")
		return 2
	}
	if options.help {
		if len(positional) != 0 || options.version || options.dryRun || options.quiet {
			fmt.Fprintln(stderr, "linefix: --help does not accept commands or other options")
			return 2
		}
		printUsage(stdout)
		return 0
	}
	if options.version {
		if len(positional) != 0 || options.dryRun || options.quiet {
			fmt.Fprintln(stderr, "linefix: --version does not accept arguments")
			return 2
		}
		fmt.Fprintf(stdout, "linefix %s\n", version)
		return 0
	}
	if len(positional) < 2 {
		fmt.Fprintln(stderr, "linefix: expected a command and at least one file")
		fmt.Fprintln(stderr, "Try 'linefix --help' for usage.")
		return 2
	}

	command, paths := positional[0], positional[1:]
	switch command {
	case "check":
		if options.dryRun || options.quiet {
			fmt.Fprintln(stderr, "linefix: --dry-run and --quiet are only valid with lf or crlf")
			return 2
		}
		failed := false
		for _, path := range paths {
			ending, err := CheckFile(path)
			if err != nil {
				printError(stderr, err)
				failed = true
				continue
			}
			if len(paths) == 1 {
				fmt.Fprintln(stdout, ending)
			} else {
				fmt.Fprintf(stdout, "%s: %s\n", path, ending)
			}
		}
		return exitCode(failed)
	case "lf", "crlf":
		failed := false
		ending := LineEnding(command)
		for _, path := range paths {
			changed, err := convertFile(path, ending, !options.dryRun)
			if err != nil {
				printError(stderr, err)
				failed = true
				continue
			}
			if options.quiet {
				continue
			}
			if changed {
				verb := "converted"
				if options.dryRun {
					verb = "would convert"
				}
				fmt.Fprintf(stdout, "%s: %s to %s\n", path, verb, displayEnding(ending))
			} else {
				fmt.Fprintf(stdout, "%s: already %s\n", path, displayEnding(ending))
			}
		}
		return exitCode(failed)
	default:
		fmt.Fprintf(stderr, "linefix: unknown command %q\n", command)
		fmt.Fprintln(stderr, "Try 'linefix --help' for usage.")
		return 2
	}
}

type cliOptions struct {
	help, version, dryRun, quiet bool
}

func parseArgs(args []string) (cliOptions, []string, error) {
	var options cliOptions
	var positional []string
	parseOptions := true
	for _, arg := range args {
		if parseOptions && arg == "--" {
			parseOptions = false
			continue
		}
		if !parseOptions || !strings.HasPrefix(arg, "-") || arg == "-" {
			positional = append(positional, arg)
			continue
		}
		switch arg {
		case "-h", "--help":
			options.help = true
		case "--version":
			options.version = true
		case "-n", "--dry-run":
			options.dryRun = true
		case "-q", "--quiet":
			options.quiet = true
		default:
			return cliOptions{}, nil, fmt.Errorf("unknown option %q", arg)
		}
	}
	return options, positional, nil
}

func exitCode(failed bool) int {
	if failed {
		return 1
	}
	return 0
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
  linefix [options] lf <file>...       Convert CRLF line endings to LF
  linefix [options] crlf <file>...     Convert LF line endings to CRLF
  linefix check <file>...              Report each file's line endings

Options:
  -n, --dry-run    Show changes without modifying files
  -q, --quiet      Suppress successful conversion output
  -h, --help       Show this help
      --version    Print the version

Use -- before a file name that begins with a dash.`)
}

func displayEnding(ending LineEnding) string {
	if ending == EndingCRLF {
		return "CRLF"
	}
	return "LF"
}
