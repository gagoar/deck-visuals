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

  It also owns onboarding — a one-off, out-of-band flow that rebuilds the presenter's
  pop-culture profile by swipe. Trigger phrases: "onboard me", "set up / update my
  familiar sources", "build my pop-culture profile".

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

## Onboarding (out-of-band)

Separate from the pipeline below, this skill maintains the presenter's analogy
profile (`assets/familiar-sources.md`). A swipe picker collects what he recognizes;
recognizing a show means he gets the jokes around it, so it earns a place. Run it on
its own — "onboard me", "update my familiar sources" — not as part of building a
deck. The pipeline's step 0 (audience familiarity) and step 5 (levity) then *read*
the refreshed profile. Full flow: `references/onboarding.md`.

## The pipeline

0. **Discovery.** Before placing any visual, run `references/story-discovery.md` —
   interview if the user has no material yet, extract if he brings an outline or
   draft. Output is a concept map: a deck-level header (`visual through-line` +
   `palette`) plus per-slide `claim` / `beat` / `story-role` / `audience` /
   `visual-intent`. Every step below reads this map; the `visual through-line` is the
   subject idiom every visual renders in. Never ask the user to name colors.
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
   Render the chosen flavor in the concept map's `visual through-line` (subject idiom)
   and let each visual carry its `story-role`, so the deck's visuals read as one set —
   idiom drives form/metaphor, brand tokens still drive color.
5. **Levity injection.** Read `references/levity.md`. Match the mapped `claim` to a
   scene (analogy-first), preferring `assets/familiar-sources.md`; fall back to
   reaction mode only for slides where the point is the audience's feeling, not a
   concept. Source a verified direct media URL with `WebSearch`/`WebFetch` plus
   `assets/tools/dv-tools check-media-url` — never a guessed URL. Insert GIF fragments
   at a tasteful cadence — openers, section breaks, the payoff. Never every slide.
   Pin the direct media URL, add alt text, respect `prefers-reduced-motion`. Inject
   `assets/fragments/levity-slide.html`. Let the human pick among the 2-3 candidates
   per slot either in chat or via the visual levity picker
   (`dv-onboard --form levity`) — see `references/levity.md`.
6. **Deliver (existing).** Ship through the user's `deliver-doc` flow. This skill
   does not deliver.

## Hard rules

- Never emit a full HTML page, a stage, or theme CSS. Fragments only.
- Every fragment carries its own scoped `<style>` and assumes nothing about the
  host page's classes or globals.
- Every chart ships with a reachable table-view twin.
- Every GIF has alt text and a `prefers-reduced-motion` fallback.
- Never embed a GIF/media URL that hasn't passed `dv-tools check-media-url` (or an
  equivalent manual check) — see `references/levity.md`.
- Colors come from `assets/brand-palette.md`, not ad hoc hex values.
- Visuals render in the deck's subject idiom (concept-map `visual through-line`) so
  they read as one set — but the idiom drives form/metaphor only, never color or
  legibility. Never ask the user to name colors; a per-deck palette is a validated
  preset picked from swatches (`dv-onboard --form palette`).
- Never run `iceberg:edit` on a rendered `.html` deck.
- The picker pages (`assets/onboarding/picker.html` and `levity-picker.html`,
  compiled into the local binary) are the full HTML pages this skill owns — local
  tools, never deck artifacts, never injected or delivered. The "fragments only" rule
  governs deck output, not these tools.

## Reference map

| Need | Read |
|---|---|
| Rebuilding the presenter profile by swipe | `references/onboarding.md` |
| The swipe catalog (config) | `assets/show-catalog.md` |
| Meme fallback for when no familiar source fits | `assets/meme-library.md` |
| The onboarding + levity picker pages (local tools) | `assets/onboarding/picker.html`, `assets/onboarding/levity-picker.html` |
| Shared picker styling (both pages) | `assets/onboarding/picker.css` |
| The runtime-free picker server (`--form onboard`/`levity`) + prebuilt binaries | `assets/onboarding/server/`, `assets/onboarding/bin/` |
| Coverage check — 3 options per concept shape | `assets/tools/dv-tools coverage` |
| Building the concept map (interview vs extract) | `references/story-discovery.md` |
| Chart vs diagram vs icon vs placeholder | `references/visual-concepts.md` |
| Chart rules, `dataviz` handoff | `references/charts.md` |
| Analogy-first levity, reaction mode, WebSearch sourcing, cadence, accessibility, taste | `references/levity.md` |
| The user's curated analogy source list | `assets/familiar-sources.md` |
| Deck theme tokens, layer boundaries, subject-idiom layer | `references/brand.md` |
| Validated 8-slot palette + ramps | `assets/brand-palette.md` |
| Per-deck palette presets (swatch-picked, validated) | `assets/palette-presets.md`, `dv-onboard --form palette` |
| Icon set | `assets/icons/` |
| Injectable fragments | `assets/fragments/` |
| Verify a GIF/media URL resolves | `assets/tools/dv-tools check-media-url` |
