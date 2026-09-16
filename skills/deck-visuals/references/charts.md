# Charts — hand off to `dataviz`, inject the result

This file does not restate chart rules. The built-in `dataviz` skill owns form
heuristics, mark specs, interaction rules, and color validation. Read it before
building any chart.

## What to do

1. **Pick the form.** Use `dataviz`'s form heuristic (`references/choosing-a-form.md`
   in that skill) to pick the chart type from the shape of the data — do not
   default to a bar chart.
2. **Pick the marks.** Use `dataviz`'s `references/marks-and-anatomy.md` for axis,
   gridline, and label specs.
3. **Pick the color.** Use `assets/brand-palette.md` from this plugin, not
   `dataviz`'s brand-neutral default palette. The brand palette is the
   deck-visuals equivalent of swapping in your own brand — validated with the
   built-in `dataviz` skill's validator, which is the authority for every
   check on that page.
4. **Build inline-SVG.** Render the chart as inline `<svg>`, not a rasterized
   image. This keeps it crisp at any deck zoom and keeps the fragment
   self-contained.
5. **Build the table-view twin.** Every chart ships with a `<table>` covering the
   same data, toggled by a visible control (not a hidden screen-reader-only table —
   a real, reachable view). This is non-negotiable: a picture of a chart loses
   hover, inspection, and per-value detail that the table view restores.
6. **Inject.** Drop both into `assets/fragments/chart-figure.html`'s structure and
   inject that fragment into the slide. Do not touch the slide shell around it.

## What `deck-visuals` adds on top of `dataviz`

- The palette swap (step 3) — `dataviz` ships brand-neutral by design; this plugin
  is the brand.
- The fragment shape — a `<figure>` wrapping chart + table + a toggle, scoped
  styles, no assumption about the host page.
- The injection point — `dataviz` does not know about slides; this skill decides
  which slide gets a chart and drops the fragment there.

## What stays in `dataviz`, do not duplicate here

- Chart form selection logic.
- Mark and axis specs.
- Interaction rules (hover, focus, tooltip behavior).
- The color-formula math. Validate a palette with the **built-in `dataviz`
  skill's validator first** — it is the authority. `assets/scripts/
  validate_palette.js` / `.py` are an independent MIT reimplementation of the
  same checks (same thresholds, same CLI shape), kept in agreement with the
  built-in validator, for use only when that skill isn't reachable — a
  standalone shell, CI, or a PR review outside a Claude Code session. See
  either script's header for detail.
