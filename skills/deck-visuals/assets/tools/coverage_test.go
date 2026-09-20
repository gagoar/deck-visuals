package main

import "testing"

func TestNormalizeShape(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Excess", "excess"},
		{"Excess (alt)", "excess"},
		{"Excess(Alt)", "excess"},
		{"Denial of consequence", "denial of consequence"},
		{"  Weird!!Punct***  ", "weird punct"},
		{"Coordination breakdown under pressure", "coordination breakdown under pressure"},
		{"Multiple   spaces", "multiple spaces"},
		{"", ""},
		{"---", ""},
		{"Trailing-Hyphen-", "trailing hyphen"},
		{"Mixed CASE Shape", "mixed case shape"},
	}
	for _, tc := range cases {
		if got := normalizeShape(tc.in); got != tc.want {
			t.Errorf("normalizeShape(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsSeparatorRow(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"---", true},
		{":---:", true},
		{"  -- ", true},
		{"Concept shape", false},
		{"Excess", false},
		{"- ", true},
	}
	for _, tc := range cases {
		if got := isSeparatorRow(tc.in); got != tc.want {
			t.Errorf("isSeparatorRow(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestPyReprStr(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Not A Real Shape", "'Not A Real Shape'"},
		{"simple", "'simple'"},
		{"has 'single' quotes", `"has 'single' quotes"`},
	}
	for _, tc := range cases {
		if got := pyReprStr(tc.in); got != tc.want {
			t.Errorf("pyReprStr(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSortedUniqueStrings(t *testing.T) {
	got := sortedUniqueStrings([]string{"Show B", "Show A", "Show B", "Show C"})
	want := []string{"Show A", "Show B", "Show C"}
	if len(got) != len(want) {
		t.Fatalf("sortedUniqueStrings length = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("sortedUniqueStrings()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseLevityShapesScopesToSeedLibrarySection(t *testing.T) {
	shapes, err := parseLevityShapes("parity/fixtures/mini-levity.md")
	if err != nil {
		t.Fatalf("parseLevityShapes: %v", err)
	}
	want := map[string]string{
		"excess":         "Excess",
		"fragility":      "Fragility",
		"recurring bugs": "Recurring bugs",
	}
	if len(shapes) != len(want) {
		t.Fatalf("got %d shapes %v, want %d %v", len(shapes), shapes, len(want), want)
	}
	for k, v := range want {
		if shapes[k] != v {
			t.Errorf("shapes[%q] = %q, want %q", k, shapes[k], v)
		}
	}
	// The reaction-beat table (Facepalm) lives outside "## Seed library" and
	// must not leak in.
	if _, ok := shapes["facepalm"]; ok {
		t.Errorf("facepalm from the reaction-beat table leaked into seed-library shapes: %v", shapes)
	}
}

func TestParseCatalogMini(t *testing.T) {
	entries, order, err := parseCatalog("parity/fixtures/mini-catalog.md")
	if err != nil {
		t.Fatalf("parseCatalog: %v", err)
	}
	wantOrder := []string{"show-a", "show-b", "show-c"}
	if len(order) != len(wantOrder) {
		t.Fatalf("order = %v, want %v", order, wantOrder)
	}
	for i := range wantOrder {
		if order[i] != wantOrder[i] {
			t.Errorf("order[%d] = %q, want %q", i, order[i], wantOrder[i])
		}
	}
	showC := entries["show-c"]
	if showC == nil {
		t.Fatalf("show-c not found in entries: %v", entries)
	}
	if showC.title != "Show C" {
		t.Errorf("show-c title = %q, want %q", showC.title, "Show C")
	}
	wantRaw := []string{"Fragility", "Recurring bugs", "Not A Real Shape"}
	if len(showC.rawShapes) != len(wantRaw) {
		t.Fatalf("show-c rawShapes = %v, want %v", showC.rawShapes, wantRaw)
	}
	for i := range wantRaw {
		if showC.rawShapes[i] != wantRaw[i] {
			t.Errorf("show-c rawShapes[%d] = %q, want %q", i, showC.rawShapes[i], wantRaw[i])
		}
	}
}
