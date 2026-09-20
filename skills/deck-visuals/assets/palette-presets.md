# Palette presets — accent themes for the palette picker

This file is **config, not documentation**. It is the pool the palette picker
(`--form palette`) offers a human to choose among at deck-build time. Two groups:

- **Favorites** — always offered, regardless of the deck's subject. Pinned by the user.
- **Matched to your library** — a theme drawn from a famous show/film's signature
  colors, offered when that title is in the presenter's recognized set (the onboarding
  `results.json` `recognized_ids` / `assets/familiar-sources.md`). Each `match:` is an
  `assets/show-catalog.md` id.

**How it works (accent themes, not full palettes).** Every theme shares the one
**validated house data palette** — the 8 categorical colors charts and diagrams use
never change (an 8-slot categorical palette must stay 8 separable hues to pass the
colorblind/contrast validator, so it can't be re-themed). Only the **accent** varies —
it colors headings, underlines, and chart callouts. Every accent below clears **WCAG
AA (≥4.5:1)** against the dark deck surfaces `#0a1428` and `#1a1a19`; a few carry an
`accent_set` (2-4 colors) for a genuinely multi-color palette in the preview.

- Shared data palette (house): `#2894eb, #d95926, #199e70, #c98500, #d55181, #008300, #9085e9, #e66767`
- Surfaces: deck bg `#0a1428`, validate surface `#1a1a19`

**Building the picker's `--data`.** Claude assembles the presets shown = **all
Favorites** + one theme per recognized title (`match:` ∈ `recognized_ids`). If nothing
matches, offer Favorites alone (house is the default).

## Favorites (always offered)

**house**
- id: house
- subject: brand default
- accent: `#3ea6ff`
- The plugin's signature accent (`--code-line` in `references/brand.md`).

**wes-anderson**
- id: wes-anderson
- subject: Wes Anderson — Grand Budapest
- accent: `#f4a8b9`
- accent_set: `#f4a8b9, #e8b04b, #4bb6a3`
- The user's favorite, always at hand. Grand Budapest pastel pink lead, with mustard +
  teal for the multi-color preview. Pink 9.75:1 on `#0a1428`; all three clear AA.

## Matched to your library

Offered when the `match:` id is in the recognized set. Accent is drawn from the
title's signature palette; all clear AA on the dark deck.

**the-matrix** — id: the-matrix · subject: The Matrix — digital rain · accent: `#22e06a` · match: the-matrix
**breaking-bad** — id: breaking-bad · subject: Breaking Bad — desert cook · accent: `#7cb518` · match: breaking-bad
**the-simpsons** — id: the-simpsons · subject: The Simpsons — Springfield yellow · accent: `#ffd90f` · match: the-simpsons
**star-wars** — id: star-wars · subject: Star Wars — lightsaber blue · accent: `#3bd6ff` · match: star-wars
**jurassic-park** — id: jurassic-park · subject: Jurassic Park — amber DNA · accent: `#ff9d2f` · match: jurassic-park
**severance** — id: severance · subject: Severance — Lumon blue · accent: `#6ab7e6` · match: severance
**stranger-things** — id: stranger-things · subject: Stranger Things — Upside Down red · accent: `#ff3b3b` · match: stranger-things
**rick-and-morty** — id: rick-and-morty · subject: Rick and Morty — portal green · accent: `#8ce651` · match: rick-and-morty
**spongebob** — id: spongebob · subject: SpongeBob — Bikini Bottom · accent: `#29c7d6` · match: spongebob
**futurama** — id: futurama · subject: Futurama — Planet Express · accent: `#46c8b0` · match: futurama
**avatar-last-airbender** — id: avatar-last-airbender · subject: Avatar — waterbender blue · accent: `#4bb1e0` · match: avatar-last-airbender
**game-of-thrones** — id: game-of-thrones · subject: Game of Thrones — winter · accent: `#7fb6d6` · match: game-of-thrones
**the-last-of-us** — id: the-last-of-us · subject: The Last of Us — cordyceps · accent: `#e0762f` · match: the-last-of-us
**mario-kart** — id: mario-kart · subject: Mario Kart — Mario red · accent: `#ff5a4d` · match: mario-kart

## Provenance

Accents were taken from each title's widely-recognized signature colors and then
checked for WCAG contrast (≥4.5:1) against the dark deck surfaces with a standalone
check (recorded in the PR). They are **accents**, not categorical slots — the shared
8-slot data palette is the one validated by `assets/tools/dv-tools validate-palette`
(see `assets/brand-palette.md`). Add a title's theme only from a signature palette the
show is actually known for, and confirm the accent clears AA before it ships.
