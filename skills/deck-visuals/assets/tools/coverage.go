// coverage.go — port of coverage_report.py: does the presenter's recognized
// pool give the levity matcher real choice?
//
// APPROVED FIX (the only intentional behavior change from the Python
// original): coverage_report.py's capacity-mode "could add" suggestion is
// buggy. In capacity mode (no --results), every catalog title is treated as
// "selected", so `carriers[shape]` and `potential[shape]` are always
// identical sets — `missing` (potential minus carriers) is therefore always
// empty, and the Python script always falls through to "no catalog title
// carries this shape — extend show-catalog.md", even when the shape *does*
// have carriers and is simply short of --target.
//
// This port corrects that for both scopes:
//   - if there are carriers not yet counted (results mode, under-recognized
//     titles) -> "could add: <titles>" (unchanged)
//   - else if the shape has any carriers at all but is still under target
//     -> "only <N> title(s) carry this shape: <titles> — extend show-catalog.md"
//   - else (no catalog title carries the shape at all)
//     -> "no catalog title carries this shape — extend show-catalog.md"
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	altParenRe = regexp.MustCompile(`\(alt\)`)
	nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)
	altTrimRe  = regexp.MustCompile(`\s*\(alt\)\s*`)
	titleRe    = regexp.MustCompile(`^\*\*(.+?)\*\*$`)
	idRe       = regexp.MustCompile(`^-\s*id:\s*(.+)$`)
	shapesRe   = regexp.MustCompile(`^-\s*Shapes:\s*(.+)$`)
)

// normalizeShape is the lenient key: lowercase, drop a trailing "(alt)",
// punctuation -> spaces, matching Python's normalize().
func normalizeShape(shape string) string {
	s := strings.ToLower(shape)
	s = altParenRe.ReplaceAllString(s, " ")
	s = nonAlnumRe.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

// isSeparatorRow reports whether s is empty or made only of "-", ":" and
// spaces — the markdown table header-separator row (or a blank first cell).
func isSeparatorRow(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if r != '-' && r != ':' && r != ' ' {
			return false
		}
	}
	return true
}

func newScanner(f *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return scanner
}

// parseLevityShapes returns {normalized: display} for every row of the
// "## Seed library" table's first column, matching parse_levity_shapes() —
// scoped to that section only, so the reaction-beat table elsewhere in
// levity.md is not mistaken for a concept shape.
func parseLevityShapes(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	shapes := map[string]string{}
	inSection := false
	scanner := newScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			inSection = strings.Contains(strings.ToLower(line), "seed library")
			continue
		}
		if !inSection {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			continue
		}
		stripped := strings.Trim(strings.TrimSpace(line), "|")
		cells := strings.Split(stripped, "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		if len(cells) == 0 {
			continue
		}
		first := cells[0]
		if isSeparatorRow(first) {
			continue
		}
		if strings.ToLower(first) == "concept shape" {
			continue
		}
		key := normalizeShape(first)
		if key == "" {
			continue
		}
		display := strings.TrimSpace(altTrimRe.ReplaceAllString(first, ""))
		if _, ok := shapes[key]; !ok {
			shapes[key] = display
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return shapes, nil
}

type catalogEntry struct {
	title     string
	shapes    []string
	rawShapes []string
}

// parseCatalog returns entries[id] = {title, shapes, rawShapes} and the
// insertion order of ids, matching parse_catalog().
func parseCatalog(path string) (map[string]*catalogEntry, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	entries := map[string]*catalogEntry{}
	var order []string
	var curID string
	var title string

	scanner := newScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		stripped := strings.TrimSpace(line)

		if m := titleRe.FindStringSubmatch(stripped); m != nil {
			title = strings.TrimSpace(m[1])
			curID = ""
			continue
		}
		if m := idRe.FindStringSubmatch(stripped); m != nil {
			curID = strings.TrimSpace(m[1])
			t := title
			if t == "" {
				t = curID
			}
			entries[curID] = &catalogEntry{title: t}
			order = append(order, curID)
			continue
		}
		if m := shapesRe.FindStringSubmatch(stripped); m != nil && curID != "" {
			rawParts := strings.Split(m[1], ";")
			var raws []string
			for _, p := range rawParts {
				p = strings.TrimSpace(p)
				if p != "" {
					raws = append(raws, p)
				}
			}
			e := entries[curID]
			e.rawShapes = raws
			e.shapes = make([]string, len(raws))
			for i, r := range raws {
				e.shapes[i] = normalizeShape(r)
			}
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return entries, order, nil
}

type resultsFile struct {
	RecognizedIDs []string `json:"recognized_ids"`
}

// pyReprStr approximates Python's repr() for a plain ASCII string: wraps in
// single quotes (or double quotes if the string itself contains a single
// quote but no double quote), escaping backslashes and the chosen quote.
func pyReprStr(s string) string {
	quote := byte('\'')
	if strings.Contains(s, "'") && !strings.Contains(s, "\"") {
		quote = '"'
	}
	var b strings.Builder
	b.WriteByte(quote)
	for _, r := range s {
		switch r {
		case rune(quote):
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte(quote)
	return b.String()
}

func sortedUniqueStrings(items []string) []string {
	set := map[string]bool{}
	for _, it := range items {
		set[it] = true
	}
	out := make([]string, 0, len(set))
	for it := range set {
		out = append(out, it)
	}
	sort.Strings(out)
	return out
}

func runCoverage(args []string) int {
	fs := flag.NewFlagSet("coverage", flag.ExitOnError)
	levity := fs.String("levity", "", "path to levity.md")
	catalog := fs.String("catalog", "", "path to show-catalog.md")
	resultsPath := fs.String("results", "", "path to results.json")
	target := fs.Int("target", 3, "minimum carriers per shape")
	fs.Parse(args) //nolint:errcheck // ExitOnError already exits(2) on parse failure

	if *levity == "" || *catalog == "" {
		fmt.Fprintln(os.Stderr, "usage: dv-tools coverage --levity <levity.md> --catalog <show-catalog.md> [--results <results.json>] [--target N]")
		return 2
	}

	levityShapes, err := parseLevityShapes(*levity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading levity file: %v\n", err)
		return 1
	}
	entries, order, err := parseCatalog(*catalog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading catalog file: %v\n", err)
		return 1
	}

	selected := map[string]bool{}
	scope := "catalog capacity (all titles)"
	if *resultsPath != "" {
		data, rerr := os.ReadFile(*resultsPath)
		if rerr != nil {
			fmt.Fprintf(os.Stderr, "error reading results file: %v\n", rerr)
			return 1
		}
		var rf resultsFile
		if jerr := json.Unmarshal(data, &rf); jerr != nil {
			fmt.Fprintf(os.Stderr, "error parsing results file: %v\n", jerr)
			return 1
		}
		for _, id := range rf.RecognizedIDs {
			selected[id] = true
		}
		scope = "recognized"
	} else {
		for _, id := range order {
			selected[id] = true
		}
	}

	counts := map[string]int{}
	carriers := map[string][]string{}
	potential := map[string][]string{}
	drift := map[string][]string{}

	for _, cid := range order {
		e := entries[cid]
		for i, norm := range e.shapes {
			raw := e.rawShapes[i]
			if _, ok := levityShapes[norm]; !ok {
				drift[raw] = append(drift[raw], e.title)
				continue
			}
			potential[norm] = append(potential[norm], e.title)
			if selected[cid] {
				counts[norm]++
				carriers[norm] = append(carriers[norm], e.title)
			}
		}
	}

	fmt.Printf("Coverage report — scope: %s, target: >= %d per shape\n\n", scope, *target)

	type shapeDisp struct{ norm, display string }
	sortedShapes := make([]shapeDisp, 0, len(levityShapes))
	for norm, display := range levityShapes {
		sortedShapes = append(sortedShapes, shapeDisp{norm, display})
	}
	sort.Slice(sortedShapes, func(i, j int) bool {
		return strings.ToLower(sortedShapes[i].display) < strings.ToLower(sortedShapes[j].display)
	})

	gaps := 0
	for _, sd := range sortedShapes {
		n := counts[sd.norm]
		ok := n >= *target
		state := "OK "
		if !ok {
			state = "GAP"
			gaps++
		}
		fmt.Printf("[%s] %s: %d\n", state, sd.display, n)
		if ok {
			continue
		}

		carrierSet := map[string]bool{}
		for _, t := range carriers[sd.norm] {
			carrierSet[t] = true
		}
		var missing []string
		for _, t := range potential[sd.norm] {
			if !carrierSet[t] {
				missing = append(missing, t)
			}
		}

		switch {
		case len(missing) > 0:
			titles := sortedUniqueStrings(missing)
			fmt.Printf("        could add: %s\n", strings.Join(titles, ", "))
		case len(potential[sd.norm]) > 0:
			titles := sortedUniqueStrings(potential[sd.norm])
			fmt.Printf("        only %d title(s) carry this shape: %s — extend show-catalog.md\n", len(titles), strings.Join(titles, ", "))
		default:
			fmt.Println("        no catalog title carries this shape — extend show-catalog.md")
		}
	}

	if len(drift) > 0 {
		fmt.Println("\nDrift — catalog shape tags with no match in levity.md:")
		rawKeys := make([]string, 0, len(drift))
		for k := range drift {
			rawKeys = append(rawKeys, k)
		}
		sort.Strings(rawKeys)
		for _, raw := range rawKeys {
			titles := sortedUniqueStrings(drift[raw])
			fmt.Printf("  %s (in: %s)\n", pyReprStr(raw), strings.Join(titles, ", "))
		}
	}

	fmt.Printf("\n%d shapes, %d gap(s), %d drift tag(s).\n", len(levityShapes), gaps, len(drift))

	if gaps > 0 || len(drift) > 0 {
		return 1
	}
	return 0
}
