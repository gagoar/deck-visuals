# Palette presets — the pool the palette picker chooses from

This file is **config, not documentation**. It is the fixed pool the palette
picker (`--form palette`) offers a human to choose among at deck-build time —
the visual analog of `show-catalog.md` for the levity picker. Every preset
here has already been run through `assets/tools/dv-tools validate-palette`
and must show `ALL CHECKS PASS`, exit 0, before it earns a place. Never add a
preset that doesn't pass; drop a themed idea rather than ship a WARN/FAIL.

## Format

Each `##` heading is one preset:

- `id:` — kebab-case, stable, unique. The picker and `results.json`
  (`chosen_preset_id`) key on it.
- `subject:` — the one-line theme tag shown under the swatches in the picker.
- `mode:` — `dark` (this plugin's decks are dark-only; see
  `references/brand.md` and `assets/brand-palette.md`).
- `surface:` — the hex the palette is validated against.
- `slots:` — the 8 categorical hexes, **in order** — reordering invalidates
  the adjacency guarantee the validator checked.

## house

- id: house
- subject: brand default
- mode: dark
- surface: `#1a1a19`
- slots: `#2894eb, #d95926, #199e70, #c98500, #d55181, #008300, #9085e9, #e66767`

The exact validated 8 slots + surface from `assets/brand-palette.md`, copied
verbatim — this is the plugin's own default categorical palette, not a new
derivation.

Validator:

```
assets/tools/dv-tools validate-palette \
  "#2894eb,#d95926,#199e70,#c98500,#d55181,#008300,#9085e9,#e66767" \
  --mode dark --surface "#1a1a19"
```

Result: **ALL CHECKS PASS**, exit 0 (identical run to `assets/brand-palette.md`'s
own record — worst CVD ΔE 8.4 protan `#199e70`↔`#c98500`, worst normal-vision
ΔE 19.3 `#c98500`↔`#d55181`, all 8 ≥ 3:1 contrast).

## space

- id: space
- subject: space / night sky
- mode: dark
- surface: `#1a1a19`
- slots: `#6a87f0, #d95926, #199e70, #c98500, #d55181, #008300, #9085e9, #e66767`

Derivation: `house` with slot 1 re-anchored from the deck-accent blue to a
cosmic indigo — same OKLCH lightness (L 0.650) and chroma (C 0.160) as
`house`'s slot 1, hue rotated from ≈248° to 270° (blue toward violet, short of
`house`'s own violet slot 7 at ≈287° so the two don't collide). Slots 2-8
unchanged. Evokes deep-space indigo/nebula without giving up any accessibility
margin.

Validator:

```
assets/tools/dv-tools validate-palette \
  "#6a87f0,#d95926,#199e70,#c98500,#d55181,#008300,#9085e9,#e66767" \
  --mode dark --surface "#1a1a19"
```

Result: **ALL CHECKS PASS**, exit 0 — worst pair is unchanged from `house`
(`#199e70`↔`#c98500` CVD ΔE 8.4; `#c98500`↔`#d55181` normal-vision ΔE 19.3),
since slot 1 isn't on either worst-case edge; all 8 ≥ 3:1 contrast.

## architecture

- id: architecture
- subject: architecture / glass & steel
- mode: dark
- surface: `#1a1a19`
- slots: `#00a4a4, #d95926, #199e70, #c98500, #d55181, #008300, #9085e9, #e66767`

Derivation: `house` with slot 1 re-anchored to a steel/glass teal-cyan — same
OKLCH lightness (L 0.650) as `house`'s slot 1, hue rotated to ≈195° (chroma
settles at C 0.111, still clear of the 0.10 floor, the tradeoff a pure hue
rotation makes to stay in the sRGB gamut at that hue). Slots 2-8 unchanged.
Warm terracotta/concrete hues (0°-150°) were tried first for slot 1 and every
one collided with the existing warm slots (2, 4, 5) badly enough to fail CVD
or normal-vision separation — dropped in favor of this cool glass/steel take,
which is the one that actually clears the bar.

Validator:

```
assets/tools/dv-tools validate-palette \
  "#00a4a4,#d95926,#199e70,#c98500,#d55181,#008300,#9085e9,#e66767" \
  --mode dark --surface "#1a1a19"
```

Result: **ALL CHECKS PASS**, exit 0 — worst pair is unchanged from `house`
(`#199e70`↔`#c98500` CVD ΔE 8.4; `#c98500`↔`#d55181` normal-vision ΔE 19.3),
since slot 1 isn't on either worst-case edge; all 8 ≥ 3:1 contrast.

## Provenance

Every result above is from the real `assets/tools/dv-tools validate-palette`
binary (build via `assets/tools/build.sh`, or `go run .` from
`assets/tools/`) — the same MIT reimplementation of Claude Code's built-in
`dataviz` validator that `assets/brand-palette.md` uses. `space` and
`architecture` were found by holding slot 1's OKLCH lightness and chroma
fixed and sweeping its hue angle, then re-checking the *real* validator at
each candidate — a pure hue rotation is provably invariant for the
Lightness-band and Normal-vision-floor checks (both are computed from
unsimulated OKLab, and rotation is a distance-preserving isometry of the
a*/b* plane), so only CVD separation and sRGB contrast needed re-verifying
per candidate, which the validator runs above did.
