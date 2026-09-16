#!/usr/bin/env python3
"""
validate_palette.py — categorical / ordinal palette validator for deck-visuals.

PROVENANCE: this is an independent MIT reimplementation of the checks
performed by Claude Code's built-in `dataviz` skill validator. That
built-in skill is the authority — when it is reachable, run it first and
treat its result as the real one. This file exists for the cases where
it isn't: a standalone shell, CI, or a PR review outside a Claude Code
session. It is kept deliberately in agreement with the built-in
validator (same thresholds, same checks, same CLI shape) rather than
copied from it — the color-science constants below (the Machado,
Oliveira & Fernandes 2009 CVD matrices, Bjoern Ottosson's OKLab
coefficients, the WCAG relative-luminance/contrast formula, and the
sRGB transfer function) are all public, standard values; the assembly
into these five checks is written from scratch here.

Checks, all computed in OKLCH / OKLab:
  1. Lightness band     - each color's OKLCH L sits inside the mode's
                          band: light [0.43, 0.77], dark [0.48, 0.67].
  2. Chroma floor         - each color's OKLCH C >= 0.10; below that a
                          hue reads as gray and stops doing identity work.
  3. CVD separation       - worst OKLab dE x100 (Euclidean distance in
                          OKLab, scaled by 100) across the active
                          pairlist (adjacent slots by default; every
                          combination with --pairs all), taking
                          min(protanopia, deuteranopia) per pair under
                          the Machado-Oliveira-Fernandes (2009)
                          severity-1.0 simulation. >=8.0 is a clean pass,
                          6.0-8.0 is a WARN legal only alongside a
                          secondary encoding (labels, gaps, texture),
                          below 6.0 is a hard FAIL. Tritanopia is
                          reported alongside for visibility but is not
                          gated.
  4. Normal-vision floor  - the same dE x100 formula, unsimulated, on the
                          same pairlist; the worst pair must clear 15.0.
                          This is a hard gate - no WARN band, and a
                          secondary encoding does not excuse it.
  5. Contrast vs surface  - WCAG contrast ratio of each color against the
                          given surface hex. >=3.0 is a pass; below that
                          is a WARN that obligates a relief channel
                          (visible direct labels or a table view), not a
                          hard FAIL.

--ordinal switches to the ramp checks for a single-hue sequential or
diverging-arm ramp instead of the five checks above:
  - Lightness monotone    - OKLCH L must strictly sort light->dark or
                          dark->light in input order.
  - Adjacent dL            - each step must move OKLCH L by >= 0.06 from
                          its neighbor.
  - Light-end contrast     - the lightest step must clear 2.0:1 against
                          the surface.
  - Single hue              - OKLab hue-angle spread across all steps
                          must stay <= 40deg, or it isn't a one-hue ramp.

Usage:
  python3 validate_palette.py "#hex,#hex,..." --mode dark --surface "#1a1a19"
  python3 validate_palette.py "#hex,#hex,..." --mode light
  python3 validate_palette.py "#hex,#hex,..." --mode dark --pairs all
  python3 validate_palette.py "#hex,#hex,..." --mode dark --ordinal

Flags:
  --mode {light,dark}    selects the lightness band and the default
                         surface. Defaults to "light".
  --surface HEX          surface color to contrast-check against.
                         Defaults to #fcfcfb (light) / #1a1a19 (dark).
  --pairs {adjacent,all} adjacent (default) checks only consecutive
                         slots - how stacks/bars/lines are read. all
                         checks every combination - required for
                         scatter/bubble/map/small-multiples, where any
                         two marks can end up side by side.
  --ordinal              validate a one-hue sequential/diverging-arm
                         ramp instead of a categorical set.

Exit code: 1 on any hard FAIL. WARN bands (CVD 6-8, sub-3:1 contrast)
still exit 0 - each requires a secondary/relief channel, but does not
fail the run by itself.
"""

import argparse
import math
import re
import sys

# -- thresholds ---------------------------------------------------------------
LIGHTNESS_BAND = {"light": (0.43, 0.77), "dark": (0.48, 0.67)}  # OKLCH L
CHROMA_FLOOR = 0.10                                             # OKLCH C
CVD_TARGET = 8.0    # OKLab dE x100, min(protan, deutan) - clean pass
CVD_FLOOR = 6.0     # OKLab dE x100 - below this is a hard FAIL
NORMAL_FLOOR = 15.0  # OKLab dE x100, unsimulated, worst pair - hard gate
CONTRAST_MIN = 3.0  # WCAG ratio vs surface
DEFAULT_SURFACE = {"light": "#fcfcfb", "dark": "#1a1a19"}
ORDINAL_MIN_DELTA_L = 0.06       # OKLCH L step between ordinal ramp entries
ORDINAL_LIGHT_END_FLOOR = 2.0    # WCAG ratio for the ramp's lightest step
ORDINAL_MAX_HUE_SPREAD = 40      # degrees - beyond this it's not one hue

# Machado, Oliveira & Fernandes (2009), full (severity 1.0) dichromacy
# simulation matrices, applied in linear sRGB.
CVD_MATRIX = {
    "protan": (
        (0.152286, 1.052583, -0.204868),
        (0.114503, 0.786281, 0.099216),
        (-0.003882, -0.048116, 1.051998),
    ),
    "deutan": (
        (0.367322, 0.860646, -0.227968),
        (0.280085, 0.672501, 0.047413),
        (-0.011820, 0.042940, 0.968881),
    ),
    "tritan": (
        (1.255528, -0.076749, -0.178779),
        (-0.078411, 0.930809, 0.147602),
        (0.004733, 0.691367, 0.303900),
    ),
}

# -- input normalization -------------------------------------------------------
_HEX_RE = re.compile(r"^#?[0-9a-fA-F]{6}$")


def is_hex(v):
    return bool(_HEX_RE.match(v))


def split_colors(raw):
    return [c for c in (s.strip() for s in (raw or "").split(",")) if c]


# -- color math -----------------------------------------------------------------
def hex_to_srgb(hex_color):
    h = hex_color.strip().lstrip("#")
    return tuple(int(h[i:i + 2], 16) / 255.0 for i in (0, 2, 4))


def srgb_channel_to_linear(c):
    return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4


def hex_to_linear(hex_color):
    return tuple(srgb_channel_to_linear(c) for c in hex_to_srgb(hex_color))


def relative_luminance(hex_color):
    r, g, b = hex_to_linear(hex_color)
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def contrast_ratio(hex_a, hex_b):
    hi, lo = sorted((relative_luminance(hex_a), relative_luminance(hex_b)), reverse=True)
    return (hi + 0.05) / (lo + 0.05)


def linear_to_oklab(rgb_lin):
    r, g, b = rgb_lin
    l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b
    m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b
    s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b
    l3, m3, s3 = l ** (1 / 3), m ** (1 / 3), s ** (1 / 3)
    return (
        0.2104542553 * l3 + 0.7936177850 * m3 - 0.0040720468 * s3,  # L
        1.9779984951 * l3 - 2.4285922050 * m3 + 0.4505937099 * s3,  # a
        0.0259040371 * l3 + 0.7827717662 * m3 - 0.8086757660 * s3,  # b
    )


def oklab_of(hex_color):
    return linear_to_oklab(hex_to_linear(hex_color))


def oklch_of(hex_color):
    L, a, b = oklab_of(hex_color)
    return (L, math.hypot(a, b))


def hue_angle_of(hex_color):
    _, a, b = oklab_of(hex_color)
    return math.degrees(math.atan2(b, a)) % 360


def simulate_cvd(hex_color, kind):
    r, g, b = hex_to_linear(hex_color)
    matrix = CVD_MATRIX[kind]

    def clamp(c):
        return max(0.0, min(1.0, c))

    return tuple(clamp(m0 * r + m1 * g + m2 * b) for m0, m1, m2 in matrix)


def delta_e_oklab(hex_a, hex_b, kind=None):
    a = linear_to_oklab(simulate_cvd(hex_a, kind)) if kind else oklab_of(hex_a)
    b = linear_to_oklab(simulate_cvd(hex_b, kind)) if kind else oklab_of(hex_b)
    return 100 * math.dist(a, b)


# -- pairlist -------------------------------------------------------------------
def build_pairlist(n, mode):
    if mode == "all":
        return [(i, j) for i in range(n) for j in range(i + 1, n)]
    return [(i, i + 1) for i in range(n - 1)]


# -- categorical checks -----------------------------------------------------------
def validate_categorical(palette, mode, surface, pairs):
    lo, hi = LIGHTNESS_BAND[mode]
    report = []
    hard_fail = False

    off_band = [c for c in palette if not (lo <= oklch_of(c)[0] <= hi)]
    if off_band:
        hard_fail = True
    report.append((
        "Lightness band", "fail" if off_band else "pass",
        (f"outside [{lo}, {hi}]: " + ", ".join(f"{c} (L={oklch_of(c)[0]:.3f})" for c in off_band))
        if off_band else f"all {len(palette)} inside L {lo}–{hi}",
    ))

    below_chroma = [c for c in palette if oklch_of(c)[1] < CHROMA_FLOOR]
    if below_chroma:
        hard_fail = True
    report.append((
        "Chroma floor", "fail" if below_chroma else "pass",
        (f"below {CHROMA_FLOOR}: " + ", ".join(f"{c} (C={oklch_of(c)[1]:.3f})" for c in below_chroma))
        if below_chroma else f"all {len(palette)} >= {CHROMA_FLOOR}",
    ))

    pairlist = build_pairlist(len(palette), pairs)
    pair_label = "all-pairs" if pairs == "all" else "adjacent"

    worst_cvd = None
    for kind in ("protan", "deutan"):
        for i, j in pairlist:
            d = delta_e_oklab(palette[i], palette[j], kind)
            if worst_cvd is None or d < worst_cvd[0]:
                worst_cvd = (d, kind, palette[i], palette[j])
    worst_tritan = min((delta_e_oklab(palette[i], palette[j], "tritan") for i, j in pairlist), default=None)
    cvd_d = worst_cvd[0] if worst_cvd else math.inf
    cvd_state = "pass" if cvd_d >= CVD_TARGET else ("warn" if cvd_d >= CVD_FLOOR else "fail")
    if cvd_state == "fail":
        hard_fail = True
    if worst_cvd:
        detail = f"worst {pair_label} {worst_cvd[2]}↔{worst_cvd[3]} ΔE {cvd_d:.1f} ({worst_cvd[1]})"
        if worst_tritan is not None:
            detail += f" · tritan {worst_tritan:.1f}"
    else:
        detail = "n/a (single color)"
    report.append(("CVD separation", cvd_state, detail))

    worst_normal = None
    for i, j in pairlist:
        d = delta_e_oklab(palette[i], palette[j])
        if worst_normal is None or d < worst_normal[0]:
            worst_normal = (d, palette[i], palette[j])
    norm_d = worst_normal[0] if worst_normal else math.inf
    norm_state = "pass" if norm_d >= NORMAL_FLOOR else "fail"
    if norm_state == "fail":
        hard_fail = True
    if worst_normal:
        detail = f"worst {pair_label} {worst_normal[1]}↔{worst_normal[2]} ΔE {norm_d:.1f}"
        if norm_state == "fail":
            detail += f" — below {NORMAL_FLOOR} floor"
    else:
        detail = "n/a (single color)"
    report.append(("Normal-vision floor", norm_state, detail))

    below_contrast = [c for c in palette if contrast_ratio(c, surface) < CONTRAST_MIN]
    report.append((
        "Contrast vs surface", "warn" if below_contrast else "pass",
        (f"below {CONTRAST_MIN}:1, relief required (labels or table view): "
         + ", ".join(f"{c} ({contrast_ratio(c, surface):.2f}:1)" for c in below_contrast))
        if below_contrast else f"all {len(palette)} >= {CONTRAST_MIN}:1",
    ))

    return report, not hard_fail


# -- ordinal ramp check -----------------------------------------------------------
def validate_ordinal(palette, mode, surface):
    report = []
    hard_fail = False
    Ls = [oklch_of(c)[0] for c in palette]

    ascending = all(Ls[i] >= Ls[i - 1] for i in range(1, len(Ls)))
    descending = all(Ls[i] <= Ls[i - 1] for i in range(1, len(Ls)))
    monotone = ascending or descending
    if not monotone:
        hard_fail = True
    report.append((
        "Lightness monotone", "pass" if monotone else "fail",
        "reads light→dark (or reverse)" if monotone
        else "not monotone: L = " + ", ".join(f"{l:.3f}" for l in Ls),
    ))

    gaps = [abs(Ls[i + 1] - Ls[i]) for i in range(len(Ls) - 1)]
    thin = [g for g in gaps if g < ORDINAL_MIN_DELTA_L]
    if thin:
        hard_fail = True
    gap_str = ", ".join(f"{g:.3f}" for g in gaps)
    report.append((
        "Adjacent ΔL", "fail" if thin else "pass",
        f"{len(thin)} step(s) below {ORDINAL_MIN_DELTA_L}: gaps = {gap_str}" if thin
        else f"all steps >= {ORDINAL_MIN_DELTA_L} (gaps = {gap_str})",
    ))

    sorted_by_l = sorted(palette, key=lambda c: oklch_of(c)[0])
    lightest_step = sorted_by_l[-1] if mode == "light" else sorted_by_l[0]
    light_contrast = contrast_ratio(lightest_step, surface)
    light_ok = light_contrast >= ORDINAL_LIGHT_END_FLOOR
    if not light_ok:
        hard_fail = True
    report.append((
        "Light-end contrast", "pass" if light_ok else "fail",
        f"{lightest_step} at {light_contrast:.2f}:1 vs surface"
        + ("" if light_ok else f" — below {ORDINAL_LIGHT_END_FLOOR}:1 floor"),
    ))

    hues = [hue_angle_of(c) for c in palette]
    spread = (max(hues) - min(hues)) if hues else 0
    if spread > 180:
        spread = 360 - spread
    one_hue = spread <= ORDINAL_MAX_HUE_SPREAD
    if not one_hue:
        hard_fail = True
    report.append((
        "Single hue", "pass" if one_hue else "fail",
        f"hue spread {spread:.0f}°"
        + ("" if one_hue else f" — exceeds {ORDINAL_MAX_HUE_SPREAD}°, not a one-hue ramp"),
    ))

    return report, not hard_fail


# -- CLI --------------------------------------------------------------------------
_STATE_GLYPH = {"pass": "PASS", "warn": "WARN", "fail": "FAIL"}


def print_report(report, ok, mode, surface, ordinal, count):
    kind = "ordinal ramp" if ordinal else "categorical"
    print(f"\nPalette ({mode}, surface {surface}, {kind}): {count} slot(s)")
    for check, state, detail in report:
        print(f"  [{_STATE_GLYPH[state]:<4}] {check:<22} {detail}")
    print(f"\n  → {'ALL CHECKS PASS' if ok else 'FAILED — fix the marked checks'}")
    if not ordinal:
        print("  WARN on CVD (6–8 floor band) is legal only with a secondary encoding.")
        print("  WARN on contrast (sub-3:1) obligates a relief channel (labels or table view).")


def main():
    ap = argparse.ArgumentParser(description="Validate a categorical or ordinal chart palette.")
    ap.add_argument("palette", help="comma-separated hex colors, e.g. \"#3ea6ff,#c860f0\"")
    ap.add_argument("--mode", choices=["light", "dark"], default="light")
    ap.add_argument("--surface", default=None, help="surface hex; defaults per --mode")
    ap.add_argument("--pairs", choices=["adjacent", "all"], default="adjacent")
    ap.add_argument("--ordinal", action="store_true", help="validate as a one-hue ramp")
    args = ap.parse_args()

    palette = split_colors(args.palette)
    if not palette:
        print('usage: python3 validate_palette.py "#hex,#hex,..." [--mode light|dark] '
              '[--surface #hex] [--pairs adjacent|all] [--ordinal]', file=sys.stderr)
        sys.exit(2)

    surface = (args.surface.strip() if args.surface else "") or DEFAULT_SURFACE[args.mode]
    bad_hex = [c for c in [*palette, surface] if not is_hex(c)]
    if bad_hex:
        print(f"invalid hex value(s): {', '.join(bad_hex)} — expected #rrggbb", file=sys.stderr)
        sys.exit(2)

    if args.ordinal:
        report, ok = validate_ordinal(palette, args.mode, surface)
    else:
        report, ok = validate_categorical(palette, args.mode, surface, args.pairs)

    print_report(report, ok, args.mode, surface, args.ordinal, len(palette))
    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
