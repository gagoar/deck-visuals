# Brand — deck theme tokens and layer boundaries

The decks this plugin serves use one fixed **house** theme: dark, high-contrast,
geometric display type. The token *shape* below — surface/panel/text roles, type,
the categorical palette — is fixed; the house theme's specific hex values are read
from a real shipped deck (`notes/src/pharos-pitch.html`, its `:root` block) — they
are not invented. Every other theme in `assets/palette-presets.md` carries the
same token shape at coordinated, re-derived values (see below) — it is not a
single hardcoded hex list.

## Theme tokens (house — the default; see "Two layers" below for other themes)

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

## Two layers: fixed theme *shape* vs per-deck coordinated theme

Being "on-brand" here has two layers, and they must not be confused:

- **Fixed theme shape (unchanged).** The token *roles* above — surface, panel,
  border, text scale, type, the 8-slot categorical palette — and the legibility
  bar every value must clear. Every visual uses this shape; `deck-visuals` never
  invents ad-hoc colors or roles.
- **Per-deck visual idiom (story-driven).** The `visual through-line` captured in
  `references/story-discovery.md` — a subject metaphor (space → orbits, architecture →
  blueprints) that drives the *forms, composition, and metaphor* of diagrams and icons
  via `references/visual-concepts.md`. It is a skin over the same primitives; it never
  overrides the tokens or the legibility bar.

**What varies per deck: accent *and* a coordinated surface.** The user picks a
**validated theme** from shown swatches (`dv-onboard --form palette`, themes in
`assets/palette-presets.md`) — not just an accent color anymore. Each theme
carries `accent` **and** matching `surface`/`panel` tokens: the same
`--bg-primary`/`--bg-panel` values, re-rendered at the accent's hue while holding
the exact OKLCH lightness and chroma of the house values fixed (a monochromatic
surface-to-accent pairing — see `assets/palette-presets.md`'s derivation note).
This is why a warm theme (Simpsons yellow, Wes Anderson pink) gets its own warm
near-black instead of the cold house navy underneath it — the anti-pattern of a
fixed blue floor under every accent is gone. The **8-slot categorical data
palette never changes** regardless of theme — only the surface/panel and the
accent vary, together. Colors are never named in prose — always picked from
swatches. See story-discovery's "Never name colors".

**Emission boundary, updated.** `deck-visuals` still never emits a shell, a
stage, or theme CSS (see "Where each layer sits" below) — but it now hands the
renderer a **surface token set** (`surface`, `panel`, plus derived
border/text-secondary/text-dim at the same hue) alongside the accent, not just
the accent alone. Whether the rendered deck's actual background adopts these
per-theme surface tokens depends on `frontend-slides` accepting them instead of
a single fixed `:root` — that is a renderer-side change, tracked separately, not
owned by this plugin.

## Where each layer sits

**`frontend-slides` (renderer).** Owns the HTML shell, the fixed 16:9 stage,
layout, animation timing, and its Pillow image pipeline (`crop_circle`,
`resize_max`). It reads these same theme tokens to build the page. `deck-visuals`
never emits a shell, a stage, or theme CSS — that would duplicate the renderer's
job and drift out of sync with it. **Pending handoff:** today it reads a single
fixed `:root` (the house values above); adopting a per-theme `surface`/`panel`
token set (this plugin now hands over, see "Emission boundary, updated" above)
is a `frontend-slides`-side change, not yet made.

**`iceberg` (copy).** Runs on slide prose before it enters the deck — drafts,
outlines, speaker notes. Never point `iceberg:edit` at a finished `.html` deck: it
treats markup as prose and will mangle inline SVG and `<img>` tags. The copy pass
happens in step 1 of the pipeline in `SKILL.md`, before render.

**`deck-visuals` (this plugin).** Sits between the two. It reads the rendered
deck's content, decides where a visual belongs, and injects a fragment scoped to
just that visual. It borrows the theme tokens above to stay on-brand; it does not
own or restate them anywhere but this reference file.
