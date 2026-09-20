package main

import (
	"math"
	"testing"
)

func approxEqual(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %v, want %v (+/- %v)", name, got, want, tol)
	}
}

func TestPyFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{8.0, "8.0"},
		{6.0, "6.0"},
		{15.0, "15.0"},
		{3.0, "3.0"},
		{2.0, "2.0"},
		{0.10, "0.1"},
		{0.06, "0.06"},
		{0.43, "0.43"},
		{0.77, "0.77"},
		{0.48, "0.48"},
		{0.67, "0.67"},
	}
	for _, tc := range cases {
		if got := pyFloat(tc.in); got != tc.want {
			t.Errorf("pyFloat(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestOklchOfBlackAndWhite(t *testing.T) {
	lWhite, cWhite := oklchOf("#ffffff")
	approxEqual(t, "L(white)", lWhite, 1.0, 0.001)
	approxEqual(t, "C(white)", cWhite, 0.0, 0.001)

	lBlack, cBlack := oklchOf("#000000")
	approxEqual(t, "L(black)", lBlack, 0.0, 0.001)
	approxEqual(t, "C(black)", cBlack, 0.0, 0.001)
}

func TestContrastRatioBlackWhiteIs21To1(t *testing.T) {
	// WCAG's own worked example: contrast of pure black against pure white
	// is exactly (1+0.05)/(0+0.05) = 21.
	approxEqual(t, "contrastRatio(white, black)", contrastRatio("#ffffff", "#000000"), 21.0, 1e-9)
	// Symmetric regardless of argument order.
	approxEqual(t, "contrastRatio(black, white)", contrastRatio("#000000", "#ffffff"), 21.0, 1e-9)
}

func TestContrastRatioIdenticalColorsIsOne(t *testing.T) {
	approxEqual(t, "contrastRatio(same, same)", contrastRatio("#2894eb", "#2894eb"), 1.0, 1e-9)
}

// TestDeltaEOklab_BrandDiverging pins the brand diverging pair's OKLab
// separation (unsimulated and under protan simulation) against the numbers
// the parity golden captured from the Python script
// (validate-palette-diverging.out: "ΔE 19.6 (protan)", "ΔE 29.4" normal).
func TestDeltaEOklab_BrandDiverging(t *testing.T) {
	normal := deltaEOklab("#2894eb", "#e66767", "")
	approxEqual(t, "deltaEOklab normal", normal, 29.4, 0.05)

	protan := deltaEOklab("#2894eb", "#e66767", "protan")
	approxEqual(t, "deltaEOklab protan", protan, 19.6, 0.05)
}

// TestDeltaEOklab_NearIdenticalBlues pins the deliberately-failing pair from
// validate-palette-fail.case: two blues one bit apart in the last hex byte.
func TestDeltaEOklab_NearIdenticalBlues(t *testing.T) {
	d := deltaEOklab("#2894eb", "#2894ec", "")
	if d >= normalFloor {
		t.Errorf("expected near-identical blues to fail the normal-vision floor (%v), got ΔE=%v", normalFloor, d)
	}
	approxEqual(t, "deltaEOklab near-identical", d, 0.2, 0.1)
}

func TestHexToSrgbAndBack(t *testing.T) {
	rgb := hexToSrgb("#ff0000")
	approxEqual(t, "r", rgb[0], 1.0, 1e-9)
	approxEqual(t, "g", rgb[1], 0.0, 1e-9)
	approxEqual(t, "b", rgb[2], 0.0, 1e-9)
}

func TestBuildPairlist(t *testing.T) {
	adj := buildPairlist(4, "adjacent")
	wantAdj := [][2]int{{0, 1}, {1, 2}, {2, 3}}
	if len(adj) != len(wantAdj) {
		t.Fatalf("adjacent pairlist = %v, want %v", adj, wantAdj)
	}
	for i := range wantAdj {
		if adj[i] != wantAdj[i] {
			t.Errorf("adjacent[%d] = %v, want %v", i, adj[i], wantAdj[i])
		}
	}

	all := buildPairlist(3, "all")
	wantAll := [][2]int{{0, 1}, {0, 2}, {1, 2}}
	if len(all) != len(wantAll) {
		t.Fatalf("all pairlist = %v, want %v", all, wantAll)
	}
	for i := range wantAll {
		if all[i] != wantAll[i] {
			t.Errorf("all[%d] = %v, want %v", i, all[i], wantAll[i])
		}
	}
}

func TestValidateCategorical_BrandPalettePasses(t *testing.T) {
	palette := []string{"#2894eb", "#d95926", "#199e70", "#c98500", "#d55181", "#008300", "#9085e9", "#e66767"}
	report, ok := validateCategorical(palette, "dark", "#1a1a19", "adjacent")
	if !ok {
		t.Fatalf("expected brand categorical palette to pass all checks, report=%+v", report)
	}
	if len(report) != 5 {
		t.Fatalf("expected 5 checks, got %d: %+v", len(report), report)
	}
}

func TestValidateCategorical_NearIdenticalBluesFails(t *testing.T) {
	_, ok := validateCategorical([]string{"#2894eb", "#2894ec"}, "dark", "#1a1a19", "adjacent")
	if ok {
		t.Fatalf("expected near-identical blues to hard-fail")
	}
}

func TestValidateOrdinal_BrandRampPasses(t *testing.T) {
	palette := []string{"#035590", "#0a6cb3", "#0b84da", "#359df5", "#6db7fd"}
	report, ok := validateOrdinal(palette, "dark", "#1a1a19")
	if !ok {
		t.Fatalf("expected brand ordinal ramp to pass all checks, report=%+v", report)
	}
	if len(report) != 4 {
		t.Fatalf("expected 4 checks, got %d: %+v", len(report), report)
	}
}

func TestIsHex(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"#2894eb", true},
		{"2894eb", true},
		{"#2894EB", true},
		{"not-a-hex", false},
		{"#2894e", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isHex(tc.in); got != tc.want {
			t.Errorf("isHex(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestParsePaletteFlags_PositionalBeforeFlags(t *testing.T) {
	f, err := parsePaletteFlags([]string{"#2894eb,#2894ec", "--mode", "dark", "--surface", "#1a1a19"})
	if err != nil {
		t.Fatalf("parsePaletteFlags: %v", err)
	}
	if f.positional != "#2894eb,#2894ec" {
		t.Errorf("positional = %q, want %q", f.positional, "#2894eb,#2894ec")
	}
	if f.mode != "dark" {
		t.Errorf("mode = %q, want %q", f.mode, "dark")
	}
	if f.surface != "#1a1a19" {
		t.Errorf("surface = %q, want %q", f.surface, "#1a1a19")
	}
}

func TestParsePaletteFlags_Ordinal(t *testing.T) {
	f, err := parsePaletteFlags([]string{"#035590,#0a6cb3", "--ordinal", "--mode", "dark"})
	if err != nil {
		t.Fatalf("parsePaletteFlags: %v", err)
	}
	if !f.ordinal {
		t.Errorf("expected ordinal=true")
	}
}
