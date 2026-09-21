# Brand — deck theme tokens and layer boundaries

The decks this plugin serves use one fixed theme: dark, high-contrast, geometric
display type. These tokens are read from a real shipped deck
(`notes/src/pharos-pitch.html`, its `:root` block) — they are not invented.

## Theme tokens

| Token | Value | Use |
|---|---|---|
| `--bg-primary` | `#0a1428` | Deck background |
| `--bg-panel` | `#101d38` | Card / panel surface |
| `--bg-panel-2` | `#152747` | Nested panel surface |
| `--border-subtle` | `rgba(147,168,201,0.18)` | Hairline dividers |
| `--border-strong` | `rgba(147,168,201,0.36)` | Emphasized dividers |
| `--text-primary` | `#eef3fb` | Body text on dark |
| `--text-secondary` | `#93a8c9` | Secondary text |
| `--text-dim` | `#56698a` | Muted text, captions |
| `--code-line` (accent) | `#3ea6ff` | Primary accent — brand-palette slot 1 anchors here |
| `--config-line` | `#c860f0` | Secondary accent (violet) |
| `--pharos` | `#ffb020` | Marker / "you are here" amber |
| `--drift` | `#ff5470` | Warning / regression red |
| `--success` | `#34d399` | Positive / success green |
| `--font-display` | `'Space Grotesk', sans-serif` | Headlines |
| `--font-body` | `'Inter', sans-serif` | Body copy |
| `--font-mono` | `'JetBrains Mono', monospace` | Data, tags, code |

`assets/brand-palette.md` derives its validated 8-slot categorical palette from
this same accent family and documents status colors that reuse `--pharos`,
`--drift`, and `--success` directly rather than inventing new ones.

## Two layers: fixed theme vs per-deck idiom

Being "on-brand" here has two layers, and they must not be confused:

- **Fixed house theme (unchanged).** The color and type tokens above. Every visual
  uses them; `deck-visuals` never invents ad-hoc colors.
- **Per-deck visual idiom (story-driven).** The `visual through-line` captured in
  `references/story-discovery.md` — a subject metaphor (space → orbits, architecture →
  blueprints) that drives the *forms, composition, and metaphor* of diagrams and icons
  via `references/visual-concepts.md`. It is a skin over the same primitives; it never
  overrides the tokens or the legibility bar.

The only case color varies per deck is a **validated accent theme** the user picked
from shown swatches (`dv-onboard --form palette`, themes in
`assets/palette-presets.md`). Every theme shares the same validated 8-slot data
palette; only the accent color — drawn from an already-validated hue family — varies
per theme. Colors are never named in prose — always picked from swatches. See
story-discovery's "Never name colors".

## Where each layer sits

**`frontend-slides` (renderer).** Owns the HTML shell, the fixed 16:9 stage,
layout, animation timing, and its Pillow image pipeline (`crop_circle`,
`resize_max`). It reads these same theme tokens to build the page. `deck-visuals`
never emits a shell, a stage, or theme CSS — that would duplicate the renderer's
job and drift out of sync with it.

**`iceberg` (copy).** Runs on slide prose before it enters the deck — drafts,
outlines, speaker notes. Never point `iceberg:edit` at a finished `.html` deck: it
treats markup as prose and will mangle inline SVG and `<img>` tags. The copy pass
happens in step 1 of the pipeline in `SKILL.md`, before render.

**`deck-visuals` (this plugin).** Sits between the two. It reads the rendered
deck's content, decides where a visual belongs, and injects a fragment scoped to
just that visual. It borrows the theme tokens above to stay on-brand; it does not
own or restate them anywhere but this reference file.
