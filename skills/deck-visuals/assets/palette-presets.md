# Palette presets — accent themes for the palette picker

This file is **config, not documentation**. It is the pool the palette picker
(`--form palette`) offers a human to choose among at deck-build time. Two groups:

- **Favorites** — always offered, regardless of the deck's subject. Pinned by the user.
- **Matched to your library** — a theme drawn from a famous show/film's signature
  colors, offered when that title is in the presenter's recognized set (the onboarding
  `results.json` `recognized_ids` / `assets/familiar-sources.md`). Each `match:` is an
  `assets/show-catalog.md` id.

**How it works (accent + coordinated surface, not full palettes).** Every theme
shares the one **validated house data palette** — the 8 categorical colors charts
and diagrams use never change (an 8-slot categorical palette must stay 8 separable
hues to pass the colorblind/contrast validator, so it can't be re-themed). The
**accent** varies per theme, and — as of this revision — so does the deck
**surface** (bg) and **panel** (card) color: each is the *same* navy tokens'
OKLCH lightness and chroma, just rotated to the accent's hue (a monochromatic
pairing — surface and accent share a hue family, so a warm accent never lands on
a cold navy floor; see the derivation note below). Only the hue rotates; lightness
and chroma are pinned, which is why every theme clears the same contrast bars
without per-theme re-tuning. Every accent below clears **WCAG AA (≥4.5:1)**
against its own theme surface (worst case 5.2:1, Stranger Things); a few carry an
`accent_set` (2-4 colors) for a genuinely multi-color palette in the preview.

- Shared data palette (house): `#2894eb, #d95926, #199e70, #c98500, #d55181, #008300, #9085e9, #e66767`
- Reference navy tokens (house theme, `references/brand.md`): surface `#0a1428`,
  panel `#101d38`, validate surface `#1a1a19`

**Surface derivation (per theme).** Convert the navy `surface`/`panel` tokens to
OKLCH, read off their lightness `L` and chroma `C` (surface: L≈0.19, C≈0.043;
panel: L≈0.235, C≈0.055 — these are slightly lighter/more saturated than
`brand.md`'s nominal `#0a1428`/`#101d38` because they're re-derived from the OKLCH
round-trip, not hand-picked), then re-render at the theme accent's hue. Every
theme's `surface`/`panel` below was generated this way and independently
contrast-checked (`assets/tools/dv-tools validate-palette ... --surface <theme
surface>`, ≥3:1) — see `assets/brand-palette.md`'s per-theme summary.

**Building the picker's `--data`.** Claude assembles the presets shown = **all
Favorites** + one theme per recognized title (`match:` ∈ `recognized_ids`). If nothing
matches, offer Favorites alone (house is the default). Each preset object needs
`id`, `subject`, `accent`, `data_slots` (the shared 8-color array, same on every
preset), and — map this file's `surface`/`panel` fields to the JSON keys
`bg`/`panel` (`palette-picker.html`'s `mockSlide` reads `preset.bg`) — this is
what makes the mock preview (and the emitted `results.json`) show the real
coordinated surface instead of falling back to navy.

## Favorites (always offered)

**house**
- id: house
- subject: brand default
- accent: `#3ea6ff` · surface: `#041627` · panel: `#051f37`
- The plugin's signature accent (`--code-line` in `references/brand.md`). Surface
  is the house hue re-derived at the pinned L/C (see above) — near-identical to
  the nominal navy, by construction.

**wes-anderson**
- id: wes-anderson
- subject: Wes Anderson — Grand Budapest
- accent: `#f4a8b9` · surface: `#240b12` · panel: `#33111b`
- accent_set: `#f4a8b9, #e8b04b, #4bb6a3`
- The user's favorite, always at hand. Grand Budapest pastel pink lead, with mustard +
  teal for the multi-color preview. Pink 9.85:1 on its own surface; all three clear AA.

## Matched to your library

Offered when the `match:` id is in the recognized set. Accent is drawn from the
title's signature palette; surface/panel are that hue at the pinned navy
lightness/chroma. All clear AA on their own theme surface.

**the-matrix** — id: the-matrix · subject: The Matrix — digital rain · accent: `#22e06a` · surface: `#051a09` · panel: `#06250f` · match: the-matrix
**breaking-bad** — id: breaking-bad · subject: Breaking Bad — desert cook · accent: `#7cb518` · surface: `#0e1803` · panel: `#162304` · match: breaking-bad
**the-simpsons** — id: the-simpsons · subject: The Simpsons — Springfield yellow · accent: `#ffd90f` · surface: `#1b1400` · panel: `#261d00` · match: the-simpsons
**star-wars** — id: star-wars · subject: Star Wars — lightsaber blue · accent: `#3bd6ff` · surface: `#001922` · panel: `#002330` · match: star-wars
**jurassic-park** — id: jurassic-park · subject: Jurassic Park — amber DNA · accent: `#ff9d2f` · surface: `#221000` · panel: `#301700` · match: jurassic-park
**severance** — id: severance · subject: Severance — Lumon blue · accent: `#6ab7e6` · surface: `#001726` · panel: `#002135` · match: severance
**stranger-things** — id: stranger-things · subject: Stranger Things — Upside Down red · accent: `#ff3b3b` · surface: `#250c0a` · panel: `#341210` · match: stranger-things
**rick-and-morty** — id: rick-and-morty · subject: Rick and Morty — portal green · accent: `#8ce651` · surface: `#0c1905` · panel: `#122407` · match: rick-and-morty
**spongebob** — id: spongebob · subject: SpongeBob — Bikini Bottom · accent: `#29c7d6` · surface: `#001a1e` · panel: `#00252a` · match: spongebob
**futurama** — id: futurama · subject: Futurama — Planet Express · accent: `#46c8b0` · surface: `#001b15` · panel: `#00261e` · match: futurama
**avatar-last-airbender** — id: avatar-last-airbender · subject: Avatar — waterbender blue · accent: `#4bb1e0` · surface: `#001825` · panel: `#002233` · match: avatar-last-airbender
**game-of-thrones** — id: game-of-thrones · subject: Game of Thrones — winter · accent: `#7fb6d6` · surface: `#001725` · panel: `#002234` · match: game-of-thrones
**the-last-of-us** — id: the-last-of-us · subject: The Last of Us — cordyceps · accent: `#e0762f` · surface: `#240e03` · panel: `#321503` · match: the-last-of-us
**mario-kart** — id: mario-kart · subject: Mario Kart — Mario red · accent: `#ff5a4d` · surface: `#250c09` · panel: `#34120f` · match: mario-kart

## User-provided (formal / brand)

For a formal deck the accent can come from the presenter's own brand instead of a
favorite or a show. The palette picker has a **custom accent** field: paste a hex and
it previews live with a WCAG contrast read on the dark deck.

**Accent-first (safe default).** The brand color rides as the **accent**; the 8
chart colors stay the validated house data palette. A brand palette rarely passes the
colorblind/contrast bar as categorical chart colors, so we don't force it there — the
brand shows up where it's safe (headings, underlines, callouts) and charts stay
colorblind-safe. The surface/panel are derived the same way as every preset — the
pasted hue at the pinned navy lightness/chroma — so a custom brand color also gets a
coordinated surface, not the fixed navy. If the pasted accent is under AA on its own
derived surface, the picker flags it and notes it needs a relief channel (direct
labels / a table view).

The chosen custom accent comes back as `{"chosen_preset_id":"custom","accent":"#…"}`.
Pasting a hex is built now; pulling colors from a brand **website** (fetch + scrape
hexes, propose, confirm via swatches) is a planned follow-up.

## Provenance

Accents were taken from each title's widely-recognized signature colors. Each
theme's `surface`/`panel` were then derived by converting the navy tokens
(`references/brand.md`) to OKLCH, reading off `L`/`C`, and re-rendering at the
accent's hue — never hand-picked. Accent-vs-own-surface WCAG contrast (≥4.5:1)
was checked for all 16 (worst case 5.23:1, Stranger Things; best 13.24:1, The
Simpsons) with a standalone check (recorded in the PR). They are **accents**, not
categorical slots — the shared 8-slot data palette is the one validated by
`assets/tools/dv-tools validate-palette` (see `assets/brand-palette.md`). Add a
title's theme only from a signature palette the show is actually known for, derive
its surface/panel the same way, and confirm the accent clears AA against its own
surface before it ships.
