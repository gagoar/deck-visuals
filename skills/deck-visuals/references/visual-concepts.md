# Visual concepts — the decision rule

Every concept slide gets one visual. Pick the flavor with this rule, in order.

The concept map from `references/story-discovery.md` feeds this rule directly: its
`claim` field is what you're evaluating in the decision rule below, and its
`visual-intent` field is a first pass at the answer (chart / diagram / icon /
humor / none) — this file's rule is what confirms or overrides that first pass, not
humor imagery alone.

## The decision rule

1. **Does the slide show a relationship, a flow, or an architecture?**
   Two or more things connect, transform, or depend on each other.
   → **Inline-SVG concept diagram.** See recipes below.

2. **Does the slide carry one single idea that needs a visual anchor?**
   A claim, a milestone, a warning — one noun, not a system.
   → **Icon** from `assets/icons/`. Pair it with the claim; do not build a diagram
     for one idea.

3. **Is the slide an emotional or hero moment** — an opener, a section title, a
   payoff — that would benefit from an illustration rather than a diagram or icon?
   → **Deferred.** See below. Do not fabricate a raster image.

If a slide fits none of these, it does not need a visual. Not every slide gets one.

## Flavor 1: inline-SVG concept diagram

Use for relationship, flow, or architecture slides. Keep every diagram:

- **Self-contained.** One `<svg>` with a scoped `<style>` block. No external fonts,
  no external images.
- **On-brand.** Stroke and fill colors come from `assets/brand-palette.md`, never
  ad hoc hex values.
- **Legible at deck scale.** Minimum stroke width 2px at the deck's native
  resolution. Label text at 14px or larger.
- **Simple.** 3-7 nodes. Past 7, the diagram needs its own slide or a simplified
  summary view.

### Recipes

**Flow (linear process, A → B → C):**
Row of rounded rectangles connected by arrows. Equal spacing. Arrowheads as a
`<marker>` reused across all edges so they render identically.

**Relationship (hub and spokes):**
One central node, 3-5 satellite nodes around it, straight or gently curved edges.
Use for "one system touches many things" slides.

**Architecture (layered boxes):**
Horizontal bands (e.g., client / service / data), boxes placed inside their band,
edges only between adjacent bands unless the slide is specifically about a
cross-layer shortcut.

**Before/after (two-state):**
Two side-by-side diagrams of the same shape, differing in one dimension (fewer
boxes, a rerouted edge, a color change on the changed node). Never redraw the
whole diagram from scratch for a one-node change — dim the unchanged nodes.

**Timeline:**
Single horizontal line, milestone dots, labels alternating above/below to avoid
collision. Use the amber/pharos-family status color for "you are here."

## Flavor 2: icon anchor

Use for single-idea slides. One icon from `assets/icons/`, placed near the claim,
sized 48-96px at deck scale. Do not stack more than one icon per idea — if the
slide needs two icons, it is carrying two ideas and should split.

Pick by literal meaning first (a growth claim gets `trending-up`, not `sparkles`).
Reach for `sparkles`, `zap`, `rocket`, or `lightbulb` only when the claim itself is
about novelty, speed, or ideas — not as generic decoration.

## Deferred: the hero-moment / AI-illustration case

Out of scope for v1. This environment has no image-generation tool, and an
illustration carries a dependency the other two flavors do not.

When a slide matches case 3 above:

1. Write a concrete image-prompt spec — subject, mood, composition, palette
   reference — as a comment in the fragment, not as prose in the deck.
2. Drop an on-brand placeholder (a bordered box with the prompt spec visible in
   dev mode, invisible in presentation mode) instead of a raster.
3. When AI illustrations ship in a later phase, the self-gating design is: generate
   the raster with an external image tool, save it into the renderer's `assets/`
   directory, and let `frontend-slides`' Pillow pipeline (`crop_circle`,
   `resize_max`) process it like any other user image. Do not reimplement that
   pipeline here.
