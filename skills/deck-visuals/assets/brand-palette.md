# Brand palette — validated (dark)

This is the **dark-mode** categorical palette for every chart and diagram this
plugin injects. The decks this plugin serves are dark-only (see
`references/brand.md`), so only dark has been derived and validated. A
light-mode set is a **separate, future selection** — dataviz treats light and
dark as independently-validated palettes, not a single hex list that auto-flips
between surfaces. Do not use these hexes on a light surface; none of the numbers
below were checked against one.

The authoritative check is Claude Code's built-in `dataviz` skill validator —
run that first when it's reachable. `assets/tools/dv-tools validate-palette`
is an MIT reimplementation of the same checks, kept in agreement with the
built-in validator, for standalone/CI/PR use where the built-in skill isn't
reachable (see `assets/tools/palette.go`'s header). Every number below is from
the real, built-in validator, not the vendored copy.

Derivation: started from dataviz's own validated dark default (`palette.md`)
and re-anchored slot 1 to a shade of the deck accent (`#3ea6ff`, see
`references/brand.md`) that fits the dark OKLCH lightness band — `#3ea6ff`
itself is too light (L 0.71, above the 0.67 ceiling), so slot 1 steps that same
hue (≈248°) down to L 0.65. The other seven hues and steps are dataviz's dark
defaults, unchanged.

## Categorical — 8 slots

| Slot | Hex | Name | Anchor |
|---|---|---|---|
| 1 | `#2894eb` | Signal Blue | shade of `#3ea6ff`, the deck accent |
| 2 | `#d95926` | Orange | dataviz dark default |
| 3 | `#199e70` | Aqua | dataviz dark default |
| 4 | `#c98500` | Yellow | dataviz dark default |
| 5 | `#d55181` | Magenta | dataviz dark default |
| 6 | `#008300` | Green | dataviz dark default |
| 7 | `#9085e9` | Violet | dataviz dark default |
| 8 | `#e66767` | Red | dataviz dark default |

Use the slots **in this order** for a legend or a series list — the validator
checks adjacent pairs in this exact order, and reordering invalidates the
adjacency guarantee below. This 8-set validates on the default *adjacent*
pairlist only (stacks, bars, lines); it has not been checked with `--pairs all`
(scatter/bubble/choropleth/small-multiples carry a lower series cap under that
harder test — see `color-formula.md` in the built-in `dataviz` skill).

### Validator command and real result

```
assets/tools/dv-tools validate-palette \
  "#2894eb,#d95926,#199e70,#c98500,#d55181,#008300,#9085e9,#e66767" \
  --mode dark --surface "#1a1a19"
```

Run against the **real, built-in `dataviz` validator** (not the vendored copy):

| Check | Result |
|---|---|
| Lightness band | PASS — all 8 inside OKLCH L 0.48–0.67 |
| Chroma floor | PASS — all 8 ≥ 0.10 |
| CVD separation | PASS — worst adjacent ΔE **8.4** (protanopia), `#c98500`↔`#199e70` |
| Normal-vision floor | PASS — worst adjacent ΔE **19.3**, `#d55181`↔`#c98500` |
| Contrast vs surface | PASS — all 8 ≥ 3:1 |

Exit code: **0**. All checks clean PASS — no WARN needed on this set.

## Sequential ramp

Single hue, the same blue as slot 1 (`#3ea6ff` family, OKLCH hue ≈ 248°), five
steps from near-surface (dark) to bright:

`#035590` → `#0a6cb3` → `#0b84da` → `#359df5` → `#6db7fd`

```
assets/tools/dv-tools validate-palette \
  "#035590,#0a6cb3,#0b84da,#359df5,#6db7fd" --mode dark --surface "#1a1a19" --ordinal
```

Real validator result — **ALL CHECKS PASS**, exit 0:

| Check | Result |
|---|---|
| Lightness monotone | PASS — steps read light→dark |
| Adjacent ΔL | PASS — all gaps ≥ 0.06 (actual ≈ 0.08 per step) |
| Light-end contrast | PASS — `#035590` (the darkest, near-surface step) at 2.24:1, clears the 2.0:1 floor |
| Single hue | PASS — hue spread 0° |

## Diverging pair

Cool pole **slot 1** (`#2894eb`, blue) against warm pole **slot 8** (`#e66767`,
red), with dataviz's own dark neutral gray as the midpoint (`#383835`):

`#2894eb` (cool / negative) — `#383835` (neutral midpoint) — `#e66767` (warm / positive)

The two poles, checked pairwise on the real validator:

```
assets/tools/dv-tools validate-palette "#2894eb,#e66767" --mode dark --surface "#1a1a19"
```

PASS, exit 0 — worst CVD ΔE **19.6** (protanopia), worst normal-vision ΔE
**29.4**, both poles ≥ 3:1 contrast. The neutral midpoint is intentionally
outside the categorical checks (chroma ≈ 0.005, reads as gray by design — that
is its job); dataviz's own dark midpoint shows the identical profile against
the categorical checker, so this is the expected shape for a diverging
midpoint, not a defect.

## Status colors

Reuse the deck's existing semantic tokens directly (see `references/brand.md`)
rather than inventing new ones:

| Status | Hex | Source token |
|---|---|---|
| Success | `#34d399` | `--success` |
| Warning / marker | `#ffb020` | `--pharos` |
| Error / regression | `#ff5470` | `--drift` |
| Info | `#3ea6ff` | `--code-line` (the deck accent itself) |

## Provenance

The authority for every check on this page is Claude Code's built-in `dataviz`
skill validator. `assets/tools/dv-tools validate-palette` is an
independent MIT reimplementation of the same checks (same thresholds, same
CLI), kept in agreement with the built-in validator for standalone/CI/PR use
where the built-in skill isn't reachable — see `assets/tools/palette.go`'s
header for detail.
