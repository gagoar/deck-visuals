// palette.go — port of validate_palette.py: categorical / ordinal palette
// validator for deck-visuals. Algorithm and thresholds are identical to
// validate_palette.py / validate_palette.js (Machado-Oliveira-Fernandes CVD
// matrices, OKLab coefficients, sRGB transfer function, WCAG contrast).
//
// The committed parity goldens were captured from the Python script, whose
// f-strings print whole-number threshold constants WITH a trailing ".0"
// (e.g. "below 15.0 floor", "8.0"), unlike the JS twin which prints "15" /
// "8". pyFloat() below reproduces that Python float-to-string formatting so
// output matches the goldens byte-for-byte.
package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// -- thresholds ---------------------------------------------------------------

var lightnessBand = map[string][2]float64{
	"light": {0.43, 0.77},
	"dark":  {0.48, 0.67},
}

const (
	chromaFloor           = 0.10
	cvdTarget             = 8.0
	cvdFloor              = 6.0
	normalFloor           = 15.0
	contrastMin           = 3.0
	ordinalMinDeltaL      = 0.06
	ordinalLightEndFloor  = 2.0
	ordinalMaxHueSpread   = 40 // degrees; an int in the Python source, printed without ".0"
)

var defaultSurface = map[string]string{"light": "#fcfcfb", "dark": "#1a1a19"}

// Machado, Oliveira & Fernandes (2009), full (severity 1.0) dichromacy
// simulation matrices, applied in linear sRGB.
var cvdMatrix = map[string][3][3]float64{
	"protan": {
		{0.152286, 1.052583, -0.204868},
		{0.114503, 0.786281, 0.099216},
		{-0.003882, -0.048116, 1.051998},
	},
	"deutan": {
		{0.367322, 0.860646, -0.227968},
		{0.280085, 0.672501, 0.047413},
		{-0.011820, 0.042940, 0.968881},
	},
	"tritan": {
		{1.255528, -0.076749, -0.178779},
		{-0.078411, 0.930809, 0.147602},
		{0.004733, 0.691367, 0.303900},
	},
}

// -- input normalization -------------------------------------------------------

var hexRe = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

func isHex(v string) bool {
	return hexRe.MatchString(v)
}

func splitColors(raw string) []string {
	var out []string
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// pyFloat formats a float the way a Python f-string embeds it with no
// explicit precision: the shortest decimal that round-trips, always with a
// decimal point (so 15.0 -> "15.0", not "15").
func pyFloat(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	if !strings.ContainsAny(s, ".eE") {
		s += ".0"
	}
	return s
}

// -- color math -----------------------------------------------------------------

func hexToSrgb(hexColor string) [3]float64 {
	h := strings.TrimPrefix(strings.TrimSpace(hexColor), "#")
	var out [3]float64
	for i, idx := range [3]int{0, 2, 4} {
		v, _ := strconv.ParseUint(h[idx:idx+2], 16, 8)
		out[i] = float64(v) / 255.0
	}
	return out
}

func srgbChannelToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func hexToLinear(hexColor string) [3]float64 {
	s := hexToSrgb(hexColor)
	return [3]float64{srgbChannelToLinear(s[0]), srgbChannelToLinear(s[1]), srgbChannelToLinear(s[2])}
}

func relativeLuminance(hexColor string) float64 {
	rgb := hexToLinear(hexColor)
	return 0.2126*rgb[0] + 0.7152*rgb[1] + 0.0722*rgb[2]
}

func contrastRatio(hexA, hexB string) float64 {
	la, lb := relativeLuminance(hexA), relativeLuminance(hexB)
	hi, lo := la, lb
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

func linearToOklab(rgbLin [3]float64) [3]float64 {
	r, g, b := rgbLin[0], rgbLin[1], rgbLin[2]
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b
	l3, m3, s3 := math.Cbrt(l), math.Cbrt(m), math.Cbrt(s)
	return [3]float64{
		0.2104542553*l3 + 0.7936177850*m3 - 0.0040720468*s3,
		1.9779984951*l3 - 2.4285922050*m3 + 0.4505937099*s3,
		0.0259040371*l3 + 0.7827717662*m3 - 0.8086757660*s3,
	}
}

func oklabOf(hexColor string) [3]float64 {
	return linearToOklab(hexToLinear(hexColor))
}

func oklchOf(hexColor string) (l, c float64) {
	lab := oklabOf(hexColor)
	return lab[0], math.Hypot(lab[1], lab[2])
}

func pymod(x, m float64) float64 {
	r := math.Mod(x, m)
	if r < 0 {
		r += m
	}
	return r
}

func hueAngleOf(hexColor string) float64 {
	lab := oklabOf(hexColor)
	deg := math.Atan2(lab[2], lab[1]) * 180 / math.Pi
	return pymod(deg, 360)
}

func simulateCvd(hexColor, kind string) [3]float64 {
	rgb := hexToLinear(hexColor)
	matrix := cvdMatrix[kind]
	clamp := func(c float64) float64 {
		if c < 0 {
			return 0
		}
		if c > 1 {
			return 1
		}
		return c
	}
	var out [3]float64
	for i := 0; i < 3; i++ {
		out[i] = clamp(matrix[i][0]*rgb[0] + matrix[i][1]*rgb[1] + matrix[i][2]*rgb[2])
	}
	return out
}

// deltaEOklab returns 100x the Euclidean OKLab distance between two colors,
// optionally under a CVD simulation ("protan"/"deutan"/"tritan"), or the
// unsimulated distance when kind == "".
func deltaEOklab(hexA, hexB, kind string) float64 {
	var a, b [3]float64
	if kind != "" {
		a = linearToOklab(simulateCvd(hexA, kind))
		b = linearToOklab(simulateCvd(hexB, kind))
	} else {
		a = oklabOf(hexA)
		b = oklabOf(hexB)
	}
	dx, dy, dz := a[0]-b[0], a[1]-b[1], a[2]-b[2]
	return 100 * math.Sqrt(dx*dx+dy*dy+dz*dz)
}

// -- pairlist -------------------------------------------------------------------

func buildPairlist(n int, mode string) [][2]int {
	var pairs [][2]int
	if mode == "all" {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				pairs = append(pairs, [2]int{i, j})
			}
		}
		return pairs
	}
	for i := 0; i < n-1; i++ {
		pairs = append(pairs, [2]int{i, i + 1})
	}
	return pairs
}

// -- report rows ------------------------------------------------------------------

type checkRow struct {
	Check, State, Detail string
}

func stateOf(fail bool) string {
	if fail {
		return "fail"
	}
	return "pass"
}

// -- categorical checks -----------------------------------------------------------

func validateCategorical(palette []string, mode, surface, pairs string) ([]checkRow, bool) {
	band := lightnessBand[mode]
	lo, hi := band[0], band[1]
	var report []checkRow
	hardFail := false

	var offBand []string
	for _, c := range palette {
		l, _ := oklchOf(c)
		if !(lo <= l && l <= hi) {
			offBand = append(offBand, c)
		}
	}
	if len(offBand) > 0 {
		hardFail = true
	}
	var lbDetail string
	if len(offBand) > 0 {
		parts := make([]string, len(offBand))
		for i, c := range offBand {
			l, _ := oklchOf(c)
			parts[i] = fmt.Sprintf("%s (L=%.3f)", c, l)
		}
		lbDetail = fmt.Sprintf("outside [%s, %s]: %s", pyFloat(lo), pyFloat(hi), strings.Join(parts, ", "))
	} else {
		lbDetail = fmt.Sprintf("all %d inside L %s–%s", len(palette), pyFloat(lo), pyFloat(hi))
	}
	report = append(report, checkRow{"Lightness band", stateOf(len(offBand) > 0), lbDetail})

	var belowChroma []string
	for _, c := range palette {
		_, cc := oklchOf(c)
		if cc < chromaFloor {
			belowChroma = append(belowChroma, c)
		}
	}
	if len(belowChroma) > 0 {
		hardFail = true
	}
	var cfDetail string
	if len(belowChroma) > 0 {
		parts := make([]string, len(belowChroma))
		for i, c := range belowChroma {
			_, cc := oklchOf(c)
			parts[i] = fmt.Sprintf("%s (C=%.3f)", c, cc)
		}
		cfDetail = fmt.Sprintf("below %s: %s", pyFloat(chromaFloor), strings.Join(parts, ", "))
	} else {
		cfDetail = fmt.Sprintf("all %d >= %s", len(palette), pyFloat(chromaFloor))
	}
	report = append(report, checkRow{"Chroma floor", stateOf(len(belowChroma) > 0), cfDetail})

	pairlist := buildPairlist(len(palette), pairs)
	pairLabel := "adjacent"
	if pairs == "all" {
		pairLabel = "all-pairs"
	}

	type worstCvdT struct {
		d    float64
		kind string
		a, b string
		set  bool
	}
	var worstCvd worstCvdT
	for _, kind := range []string{"protan", "deutan"} {
		for _, p := range pairlist {
			d := deltaEOklab(palette[p[0]], palette[p[1]], kind)
			if !worstCvd.set || d < worstCvd.d {
				worstCvd = worstCvdT{d, kind, palette[p[0]], palette[p[1]], true}
			}
		}
	}
	var worstTritan float64
	haveTritan := len(pairlist) > 0
	if haveTritan {
		worstTritan = math.Inf(1)
		for _, p := range pairlist {
			d := deltaEOklab(palette[p[0]], palette[p[1]], "tritan")
			if d < worstTritan {
				worstTritan = d
			}
		}
	}
	cvdD := math.Inf(1)
	if worstCvd.set {
		cvdD = worstCvd.d
	}
	cvdState := "fail"
	if cvdD >= cvdTarget {
		cvdState = "pass"
	} else if cvdD >= cvdFloor {
		cvdState = "warn"
	}
	if cvdState == "fail" {
		hardFail = true
	}
	var cvdDetail string
	if worstCvd.set {
		cvdDetail = fmt.Sprintf("worst %s %s↔%s ΔE %.1f (%s)", pairLabel, worstCvd.a, worstCvd.b, cvdD, worstCvd.kind)
		if haveTritan {
			cvdDetail += fmt.Sprintf(" · tritan %.1f", worstTritan)
		}
	} else {
		cvdDetail = "n/a (single color)"
	}
	report = append(report, checkRow{"CVD separation", cvdState, cvdDetail})

	type worstNormalT struct {
		d    float64
		a, b string
		set  bool
	}
	var worstNormal worstNormalT
	for _, p := range pairlist {
		d := deltaEOklab(palette[p[0]], palette[p[1]], "")
		if !worstNormal.set || d < worstNormal.d {
			worstNormal = worstNormalT{d, palette[p[0]], palette[p[1]], true}
		}
	}
	normD := math.Inf(1)
	if worstNormal.set {
		normD = worstNormal.d
	}
	normState := "pass"
	if normD < normalFloor {
		normState = "fail"
	}
	if normState == "fail" {
		hardFail = true
	}
	var normDetail string
	if worstNormal.set {
		normDetail = fmt.Sprintf("worst %s %s↔%s ΔE %.1f", pairLabel, worstNormal.a, worstNormal.b, normD)
		if normState == "fail" {
			normDetail += fmt.Sprintf(" — below %s floor", pyFloat(normalFloor))
		}
	} else {
		normDetail = "n/a (single color)"
	}
	report = append(report, checkRow{"Normal-vision floor", normState, normDetail})

	var belowContrast []string
	for _, c := range palette {
		if contrastRatio(c, surface) < contrastMin {
			belowContrast = append(belowContrast, c)
		}
	}
	var contrastDetail string
	if len(belowContrast) > 0 {
		parts := make([]string, len(belowContrast))
		for i, c := range belowContrast {
			parts[i] = fmt.Sprintf("%s (%.2f:1)", c, contrastRatio(c, surface))
		}
		contrastDetail = fmt.Sprintf("below %s:1, relief required (labels or table view): %s", pyFloat(contrastMin), strings.Join(parts, ", "))
	} else {
		contrastDetail = fmt.Sprintf("all %d >= %s:1", len(palette), pyFloat(contrastMin))
	}
	contrastState := "pass"
	if len(belowContrast) > 0 {
		contrastState = "warn"
	}
	report = append(report, checkRow{"Contrast vs surface", contrastState, contrastDetail})

	return report, !hardFail
}

// -- ordinal ramp check -----------------------------------------------------------

func validateOrdinal(palette []string, mode, surface string) ([]checkRow, bool) {
	var report []checkRow
	hardFail := false

	ls := make([]float64, len(palette))
	for i, c := range palette {
		l, _ := oklchOf(c)
		ls[i] = l
	}

	ascending, descending := true, true
	for i := 1; i < len(ls); i++ {
		if ls[i] < ls[i-1] {
			ascending = false
		}
		if ls[i] > ls[i-1] {
			descending = false
		}
	}
	monotone := ascending || descending
	if !monotone {
		hardFail = true
	}
	var lmDetail string
	if monotone {
		lmDetail = "reads light→dark (or reverse)"
	} else {
		parts := make([]string, len(ls))
		for i, l := range ls {
			parts[i] = fmt.Sprintf("%.3f", l)
		}
		lmDetail = "not monotone: L = " + strings.Join(parts, ", ")
	}
	report = append(report, checkRow{"Lightness monotone", stateOf(!monotone), lmDetail})

	gaps := make([]float64, 0, len(ls)-1)
	for i := 0; i < len(ls)-1; i++ {
		gaps = append(gaps, math.Abs(ls[i+1]-ls[i]))
	}
	thin := 0
	for _, g := range gaps {
		if g < ordinalMinDeltaL {
			thin++
		}
	}
	if thin > 0 {
		hardFail = true
	}
	gapParts := make([]string, len(gaps))
	for i, g := range gaps {
		gapParts[i] = fmt.Sprintf("%.3f", g)
	}
	gapStr := strings.Join(gapParts, ", ")
	var adjDetail string
	if thin > 0 {
		adjDetail = fmt.Sprintf("%d step(s) below %s: gaps = %s", thin, pyFloat(ordinalMinDeltaL), gapStr)
	} else {
		adjDetail = fmt.Sprintf("all steps >= %s (gaps = %s)", pyFloat(ordinalMinDeltaL), gapStr)
	}
	report = append(report, checkRow{"Adjacent ΔL", stateOf(thin > 0), adjDetail})

	sortedIdx := make([]int, len(palette))
	for i := range sortedIdx {
		sortedIdx[i] = i
	}
	sort.SliceStable(sortedIdx, func(i, j int) bool { return ls[sortedIdx[i]] < ls[sortedIdx[j]] })
	var lightestStep string
	if len(sortedIdx) > 0 {
		if mode == "light" {
			lightestStep = palette[sortedIdx[len(sortedIdx)-1]]
		} else {
			lightestStep = palette[sortedIdx[0]]
		}
	}
	lightContrast := contrastRatio(lightestStep, surface)
	lightOk := lightContrast >= ordinalLightEndFloor
	if !lightOk {
		hardFail = true
	}
	lcDetail := fmt.Sprintf("%s at %.2f:1 vs surface", lightestStep, lightContrast)
	if !lightOk {
		lcDetail += fmt.Sprintf(" — below %s:1 floor", pyFloat(ordinalLightEndFloor))
	}
	report = append(report, checkRow{"Light-end contrast", stateOf(!lightOk), lcDetail})

	hues := make([]float64, len(palette))
	for i, c := range palette {
		hues[i] = hueAngleOf(c)
	}
	var spread float64
	if len(hues) > 0 {
		mx, mn := hues[0], hues[0]
		for _, h := range hues {
			if h > mx {
				mx = h
			}
			if h < mn {
				mn = h
			}
		}
		spread = mx - mn
	}
	if spread > 180 {
		spread = 360 - spread
	}
	oneHue := spread <= float64(ordinalMaxHueSpread)
	if !oneHue {
		hardFail = true
	}
	shDetail := fmt.Sprintf("hue spread %.0f°", spread)
	if !oneHue {
		shDetail += fmt.Sprintf(" — exceeds %d°, not a one-hue ramp", ordinalMaxHueSpread)
	}
	report = append(report, checkRow{"Single hue", stateOf(!oneHue), shDetail})

	return report, !hardFail
}

// -- CLI --------------------------------------------------------------------------

var stateGlyph = map[string]string{"pass": "PASS", "warn": "WARN", "fail": "FAIL"}

func printPaletteReport(report []checkRow, ok bool, mode, surface string, ordinal bool, count int) {
	kind := "categorical"
	if ordinal {
		kind = "ordinal ramp"
	}
	fmt.Printf("\nPalette (%s, surface %s, %s): %d slot(s)\n", mode, surface, kind, count)
	for _, r := range report {
		fmt.Printf("  [%-4s] %-22s %s\n", stateGlyph[r.State], r.Check, r.Detail)
	}
	verdict := "FAILED — fix the marked checks"
	if ok {
		verdict = "ALL CHECKS PASS"
	}
	fmt.Printf("\n  → %s\n", verdict)
	if !ordinal {
		fmt.Println("  WARN on CVD (6–8 floor band) is legal only with a secondary encoding.")
		fmt.Println("  WARN on contrast (sub-3:1) obligates a relief channel (labels or table view).")
	}
}

// paletteFlags is a hand-rolled parser (not flag.FlagSet) because the CLI
// shape here — a positional palette argument followed by flags — is
// something Go's stdlib flag package cannot parse: flag.Parse stops at the
// first non-flag argument, and the positional comes first in every real
// invocation (matching validate_palette.py's argparse and .js's parseArgs).
type paletteFlags struct {
	mode, surface, pairs, positional string
	ordinal                          bool
	havePositional                   bool
}

func parsePaletteFlags(argv []string) (*paletteFlags, error) {
	f := &paletteFlags{}
	valueFlags := map[string]*string{"--mode": &f.mode, "--surface": &f.surface, "--pairs": &f.pairs}

	for i := 0; i < len(argv); i++ {
		a := argv[i]

		if eq := strings.Index(a, "="); eq > 0 && strings.HasPrefix(a, "--") {
			key, val := a[:eq], a[eq+1:]
			if ptr, ok := valueFlags[key]; ok {
				*ptr = val
				continue
			}
		}

		if ptr, ok := valueFlags[a]; ok {
			if i+1 >= len(argv) {
				return nil, fmt.Errorf("missing value for %s", a)
			}
			i++
			*ptr = argv[i]
			continue
		}
		if a == "--ordinal" {
			f.ordinal = true
			continue
		}
		if strings.HasPrefix(a, "--") {
			return nil, fmt.Errorf("unknown flag: %s", a)
		}
		if !f.havePositional {
			f.positional = a
			f.havePositional = true
			continue
		}
		return nil, fmt.Errorf("unexpected extra positional: %s", a)
	}
	return f, nil
}

func runValidatePalette(args []string) int {
	usage := `usage: dv-tools validate-palette "#hex,#hex,..." [--mode light|dark] [--surface #hex] [--pairs adjacent|all] [--ordinal]`

	f, err := parsePaletteFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 2
	}

	palette := splitColors(f.positional)
	if len(palette) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		return 2
	}

	mode := f.mode
	if mode == "" {
		mode = "light"
	}
	if mode != "light" && mode != "dark" {
		fmt.Fprintf(os.Stderr, "--mode must be \"light\" or \"dark\" (got %q)\n", mode)
		return 2
	}

	pairs := f.pairs
	if pairs == "" {
		pairs = "adjacent"
	}
	if pairs != "adjacent" && pairs != "all" {
		fmt.Fprintf(os.Stderr, "--pairs must be \"adjacent\" or \"all\" (got %q)\n", pairs)
		return 2
	}

	surface := strings.TrimSpace(f.surface)
	if surface == "" {
		surface = defaultSurface[mode]
	}

	allColors := append(append([]string{}, palette...), surface)
	var bad []string
	for _, c := range allColors {
		if !isHex(c) {
			bad = append(bad, c)
		}
	}
	if len(bad) > 0 {
		fmt.Fprintf(os.Stderr, "invalid hex value(s): %s — expected #rrggbb\n", strings.Join(bad, ", "))
		return 2
	}

	var report []checkRow
	var ok bool
	if f.ordinal {
		report, ok = validateOrdinal(palette, mode, surface)
	} else {
		report, ok = validateCategorical(palette, mode, surface, pairs)
	}
	printPaletteReport(report, ok, mode, surface, f.ordinal, len(palette))
	if ok {
		return 0
	}
	return 1
}
