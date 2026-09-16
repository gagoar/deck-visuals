---
name: deck-visuals
description: >
  Guide and inject impactful visuals into an HTML deck — inline-SVG concept diagrams,
  dataviz charts, icon anchors, and GIF levity. This skill does not render decks.

  Use when the user asks to add visuals to a deck, make a slide impactful, chart a
  number, add a diagram, inject an icon, or add a GIF/levity slide. Trigger phrases:
  "add visuals to my deck", "make this slide impactful", "chart this", "add a
  gif/levity slide", "inject a diagram", "this slide needs an icon", "visualize this
  data in the deck".

  Division of labor: `frontend-slides` renders the deck — shell, stage, theme,
  animations. This skill never touches that. It reads the deck content, decides
  where a visual belongs and which kind, builds that one visual as a small
  self-contained HTML fragment, and hands it back for injection into the slide
  already on the page. It leans on the built-in `dataviz` skill for chart rules and
  on `iceberg` for copy — it does not re-implement either.

  Do not use this skill to build a deck from scratch, write slide copy, or design
  the stage/theme CSS — that is `frontend-slides` and `iceberg`.
---

# deck-visuals

This skill decides **where** a visual belongs on a slide and **which kind**, then
hands back an injectable fragment. It does not render. `frontend-slides` owns the
shell, the stage, and the animations — never edit that layer from here.

## The pipeline

1. **Copy.** Draft slide text. Run `iceberg` on the prose only, never on the
   finished `.html` file — `iceberg` reads markup as prose and will corrupt SVG or
   `<img>` tags.
2. **Render (existing).** Open or generate the deck with `frontend-slides`. Do not
   touch this output directly; only inject into it.
3. **Chart injection.** For each number that earns a chart, read
   `references/charts.md`. It points to the built-in `dataviz` skill for form,
   marks, interaction, and color. Build an inline-SVG chart plus a table-view twin.
   Inject `assets/fragments/chart-figure.html`.
4. **Concept-visual injection.** For each concept slide, read
   `references/visual-concepts.md` and apply its decision rule:
   - relationship / flow / architecture → inline-SVG **concept diagram**
     (`assets/fragments/concept-diagram.html`)
   - single idea needing an anchor → **icon** from `assets/icons/`
   - emotional / hero moment → **deferred**. Drop a placeholder and an image-prompt
     spec; do not fabricate a raster. See `references/visual-concepts.md`.
5. **Levity injection.** Read `references/levity.md`. Start from the slide's beat,
   propose 2-3 reference options, and source a verified direct media URL with
   `WebSearch`/`WebFetch` plus `assets/scripts/check_media_url.py` — never a
   guessed URL. Insert GIF fragments at a tasteful cadence — openers, section
   breaks, the payoff. Never every slide. Pin the direct media URL, add alt text,
   respect `prefers-reduced-motion`. Inject `assets/fragments/levity-slide.html`.
6. **Deliver (existing).** Ship through the user's `deliver-doc` flow. This skill
   does not deliver.

## Hard rules

- Never emit a full HTML page, a stage, or theme CSS. Fragments only.
- Every fragment carries its own scoped `<style>` and assumes nothing about the
  host page's classes or globals.
- Every chart ships with a reachable table-view twin.
- Every GIF has alt text and a `prefers-reduced-motion` fallback.
- Never embed a GIF/media URL that hasn't passed `check_media_url.py` (or an
  equivalent manual check) — see `references/levity.md`.
- Colors come from `assets/brand-palette.md`, not ad hoc hex values.
- Never run `iceberg:edit` on a rendered `.html` deck.

## Reference map

| Need | Read |
|---|---|
| Chart vs diagram vs icon vs placeholder | `references/visual-concepts.md` |
| Chart rules, `dataviz` handoff | `references/charts.md` |
| GIF ideation, WebSearch sourcing, cadence, accessibility, taste | `references/levity.md` |
| Deck theme tokens, layer boundaries | `references/brand.md` |
| Validated 8-slot palette + ramps | `assets/brand-palette.md` |
| Icon set | `assets/icons/` |
| Injectable fragments | `assets/fragments/` |
| Verify a GIF/media URL resolves | `assets/scripts/check_media_url.py` |
