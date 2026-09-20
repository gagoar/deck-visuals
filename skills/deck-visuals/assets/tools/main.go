// Command dv-tools — Go port of the three deck-visuals standard-library Python
// scripts (check_media_url.py, coverage_report.py, validate_palette.py).
//
// Each script becomes a subcommand:
//
//	dv-tools check-media-url <url> [<url> ...]
//	dv-tools coverage --levity <levity.md> --catalog <show-catalog.md> [--results <results.json>] [--target N]
//	dv-tools validate-palette "#hex,#hex,..." [--mode light|dark] [--surface #hex] [--pairs adjacent|all] [--ordinal]
//
// This binary is built to reproduce the Python scripts' stdout byte-for-byte
// (see skills/deck-visuals/assets/tools/parity), with one intentional,
// approved behavior fix in the coverage subcommand's capacity-mode GAP
// suggestion text (see coverage.go).
package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: dv-tools <check-media-url|coverage|validate-palette> [args...]")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := os.Args[1]
	rest := os.Args[2:]

	var code int
	switch sub {
	case "check-media-url":
		code = runCheckMediaURL(rest)
	case "coverage":
		code = runCoverage(rest)
	case "validate-palette":
		code = runValidatePalette(rest)
	default:
		usage()
		os.Exit(2)
	}
	os.Exit(code)
}
