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

## Render in the deck's idiom

Once the flavor is chosen, render it in the concept map's **`visual through-line`**
(the user-supplied mapping when given, otherwise the inferred-and-approved one; empty
→ clean house idiom). The flavor stays the same primitive; the idiom changes its skin,
so the deck's visuals read as one designed set instead of one-offs.

Map each diagram recipe to the idiom:

| Primitive | House (default) | Space idiom | Architecture idiom |
|---|---|---|---|
| Hub-and-spoke | isometric plates — hub elevated + larger, satellites on a lower tier, struts between | a sun with orbiting planets | a keystone with radiating members |
| Flow A→B→C | chamfered flat-geometric nodes + bold chevron arrows | probes along a trajectory | stages on a blueprint plan line |
| Layered architecture | isometric plates, one elevation tier per band | orbital shells | building floors / blueprint elevations |
| Timeline | chamfered line + milestone markers | a flight path with waypoints | a construction schedule / gantt bars |
| Before/after | two side-by-side states, in whichever primitive above they demonstrate | two orbital configs | plan vs elevation |

Icons lean into the idiom's vocabulary before the generic set (a space deck reaches
for `rocket`/`sparkles` only when the *claim* is about launch/novelty — the idiom
guides, it does not license decoration).

**Guardrails — the idiom never wins over these:**
- Self-contained inline SVG with a scoped `<style>`; no external fonts/images.
- **Colors come from `assets/brand-palette.md` only** (or a picked preset on the
  concept-map header) — the idiom drives *shape and metaphor, never color*.
- ≥2px strokes, ≥14px labels, 3-7 nodes; a relationship diagram must still read as a
  relationship diagram even when drawn as orbits. If the skin hurts legibility, drop
  it and use the clean house primitive.

**Coherence rule:** reuse the same idiom primitives across the whole deck — one motif
family, not a new metaphor per slide. Two diagrams in the same deck should look like
siblings.

## Flavor 1: inline-SVG concept diagram

Use for relationship, flow, or architecture slides. Keep every diagram:

- **Self-contained.** One `<svg>` with a scoped `<style>` block. No external fonts,
  no external images, no external library (the isometric construction below is
  plain SVG polygons/transforms — never pull in an isometric/3D library).
- **On-brand.** Stroke and fill colors come from `assets/brand-palette.md`
  (the 8-slot categorical palette) plus the deck's `--dv-*` theme tokens, never
  ad hoc hex values, and never a derived/shaded tint of either — see "No new
  hex values" under the isometric construction below for why that constraint
  is easy to keep.
- **Legible at deck scale.** Minimum stroke width 2px at the deck's native
  resolution (the recipes below default to 3px for more visual weight — 2px is
  the floor, not the target). Label text at 14px or larger, and always
  horizontal — never rotated or skewed, even inside an isometric diagram.
- **Simple.** 3-7 nodes. Past 7, the diagram needs its own slide or a simplified
  summary view.

### The shape vocabulary (replaces plain rounded rectangles)

Two related constructions, chosen per recipe below — never mixed within one
diagram, and never mixed across diagrams in the same deck (coherence rule).

**Chamfered node (flat, no projection) — Flow, Timeline, and whichever
primitive Before/after is demonstrating.** A rectangle with its four corners
cut at 45°, instead of a rounded rectangle — reads as more geometric and
confident, and is one `<polygon>`:

```
points = (x+c,y) (x+w-c,y) (x+w,y+c) (x+w,y+h-c)
         (x+w-c,y+h) (x+c,y+h) (x,y+h-c) (x,y+c)
```
(`c` = chamfer size, roughly 10-15% of the shorter side.) Fill = the theme
panel token (`--dv-panel`); stroke = the node's categorical slot color, 3px.
Add a small solid **triangle** in the node's categorical slot color at one
corner — `(x,y) (x+t,y) (x,y+t)` with `t` somewhat larger than `c` (e.g. 22 vs
14) is fine; the extra reach just paints over part of the node's own fill,
which is harmless since it stays inside the polygon's real edges. This
replaces "colored stroke only" as the category signal. **Never a plain
axis-aligned `<rect>` here** — see "Mistakes to never repeat" below for why.
Arrowheads: one reused `<marker>` per diagram, a bolder solid chevron/triangle
(more visual weight than a thin arrow) — `M0 0L12 6L0 12L4 6Z`, with `refX`
set to the path's own tip x-coordinate exactly (12, not an approximate
nearby value — see "Mistakes to never repeat"). **`markerWidth`/`markerHeight`
of 4** (with the default `markerUnits="strokeWidth"` and a 3px stroke, that
renders at 4×3=12 user-space units) — check this against the actual gap
between nodes, not just against how the arrowhead looks in isolation; a
default marker size chosen only "does this look bold enough" can end up
larger than the whole gap it has to fit in (see "Mistakes to never repeat").
Connectors 3px.

**Isometric plate (projected, dimensional) — Relationship and Architecture
only,** because both are already about physical/structural hierarchy (a hub
coordinating spokes; layers stacked on each other) — depth earns its keep
here in a way it wouldn't for a linear flow or timeline.

A plate is a flat top diamond **plus two extruded rim faces** beneath it —
not a bare outline, and not a solid shaded cube either. A bare diamond
outline does not read as dimensional; it just reads as a rotated 2D shape.
The rim is what sells the depth, and it costs no new color:

```
top diamond: (cx,cy-ph/2) (cx+pw/2,cy) (cx,cy+ph/2) (cx-pw/2,cy)
d = max(8, ph*0.2)                                      // rim depth
left rim:  (cx-pw/2,cy) (cx,cy+ph/2) (cx,cy+ph/2+d) (cx-pw/2,cy+d)
right rim: (cx,cy+ph/2) (cx+pw/2,cy) (cx+pw/2,cy+d) (cx,cy+ph/2+d)
```
Draw both rim faces first, then the top diamond on top of them (draw order
matters — the top diamond must cover the rims' upper edge). Fill each rim
with the plate's own categorical slot hex — **the exact same hex as the
diamond's stroke, just at reduced fill-opacity** (left rim ≈0.35, right rim
≈0.2, giving a lit-top/shadowed-side read from one color) — never a darker or
lighter *derived* hex. This is the whole trick for staying inside the
validated-palette constraint: opacity is not a new color, a mixed/blended hex
would be. Top diamond: fill = theme panel token, stroke = the categorical
slot, 3-4px. Add the same small accent-triangle tag as the chamfered node,
rotated 45° into a small diamond at the plate's top vertex.

(`pw`/`ph` at a fixed ~2:1 ratio, e.g. 120×60, the same for every plate in a
diagram.)

**Text stays flat.** A chamfered node is already unrotated, so its label sits
centered directly inside it, same as a plain rectangle would. An isometric
plate is a rotated diamond, too thin and angled to hold text directly — its
label instead sits in its own small, unrotated, rounded-rect tag (fill =
theme panel, 1.5px border in `--dv-border-strong`) positioned just below the
plate. Never skew or rotate text with its shape; this is what keeps an
isometric diagram exactly as legible as a flat one.

### Mistakes to never repeat

Five specific construction bugs shipped in the first pass of this shape
vocabulary, each one found by actually looking closely at a rendered
example — not by reasoning about the math in the abstract. Keep them
documented here precisely so the next generation doesn't re-derive them the
hard way:

1. **A tag is geometry-aware, never a floating rect.** The first version of
   the category tag was a plain axis-aligned `<rect>` dropped at a chamfered
   node's bounding-box corner — but that corner is a diagonal cut, so the
   rect's own corner poked outside the polygon's real edge, with the node's
   stroke cutting messily across it. Fix: the tag is a **triangle** whose
   points lie on the shape's own boundary (the corner + two points on the
   real edges), so it is either inside the polygon or in the exact notch the
   chamfer cut away — never floating past an edge into nothing.
2. **A bare outline is not depth.** The first isometric plates were just the
   top-diamond outline with no rim — which reads as "a rotated 2D shape,"
   not as a solid tile. Depth needs a visible side face (the extruded rim
   above), even though the projection angle alone is technically correct
   isometric math.
3. **Connectors attach to a shape's true outer boundary, never its center —
   and "outer boundary" means the boundary *including* any depth cue.** Once
   plates gained rim faces, struts still ran center-to-center; the opaque top
   diamond hid the strut's path through the *diamond*, but the strut also
   passed through the translucent rim below the diamond, which doesn't fully
   hide it — a faint line visibly bled through right at each plate's bottom
   tip. Fix: struts depart from the rim's own bottom vertex
   (`cy+ph/2+d`), not the plate's center.
4. **A marker's `refX` must equal its path's own tip coordinate, exactly —
   not "close to it."** The bolder chevron arrowhead (`M0 0L12 6L0 12L4 6Z`,
   tip at x=12) originally used `refX="10"`, two units short of the tip.
   Marker units scale by the referencing element's stroke-width, so that
   2-unit gap becomes a several-pixel overhang past the line's actual
   endpoint — landing the tip inside the *next* node, which then draws on
   top and hides it. `refX` must exactly match the tip's own x-coordinate in
   the path data, so nothing overhangs past the endpoint.
5. **A marker's *rendered* size must be checked against the actual gap it
   sits in, not chosen by "does this look bold enough."** Fixing #4 above
   (correct `refX`) immediately exposed a second bug: the marker's
   `markerWidth`/`markerHeight` of 9 (default `markerUnits="strokeWidth"`, 3px
   stroke) renders at 9×3=27 user-space units — *larger than the 24-unit gap
   between nodes*. Once the tip was correctly anchored at the line's
   endpoint, the whole oversized chevron body extended backward past the
   entire gap, swallowing both the connector line and part of the previous
   node — "just an arrowhead, no visible line" is exactly what an
   over-large marker looks like. Fixed by dropping to `markerWidth="4"`
   (rendered 12 units — half the gap, leaving the other half for a visible
   line). Always compute the marker's actual rendered size
   (`markerWidth × stroke-width` under the default `markerUnits`) and
   compare it to the real gap the connector spans — never eyeball it in
   isolation from the diagram's actual node spacing.

The common thread: verify by rendering and visually inspecting (or, better,
by measuring actual element bounding boxes / overlap programmatically) — not
by reasoning through the coordinate math and assuming it holds.

### Recipes

**Flow (linear process, A → B → C):**
Row of chamfered nodes connected by bold chevron arrows, equal spacing, each
node's category-square tag a distinct categorical slot. Same reused `<marker>`
across all edges.

**Relationship (hub and spokes):**
One hub plate — elevated (drawn higher on the canvas) and larger/thicker-
stroked than its satellites — with 3-5 satellite plates on a shared lower
tier, arranged in a fan beneath it. Straight struts from the hub's rim-bottom
tip (`hub_cy+hub_ph/2+d` — **not the hub's center**, see "Mistakes to never
repeat") to each satellite's center (a strut arriving at a satellite's center
from above is already fully hidden by that satellite's own opaque top
diamond, since a plate's rim sits below its center, not above — only the
departure end needs the rim-tip fix). Isometric projection keeps parallel
lines parallel, so straight struts read correctly without curve math. Use for
"one system touches many things" slides.

**Architecture (layered boxes):**
Each band (e.g., client / service / data) is an elevation tier: same plate
size within a tier, a consistent per-tier offset (shift each tier up and
slightly sideways relative to the one below, suggesting a rising stack), and
a faint dashed horizontal guideline per tier (`--dv-text-secondary` at ~12%
opacity) tying a tier's plates and its label together without competing with
the content. Struts only between adjacent tiers, **from the upper tier's
rim-bottom tip down to the lower tier's center** (same rule as Relationship,
same reason), unless the slide is specifically about a cross-layer shortcut.

**Before/after (two-state):**
Two side-by-side diagrams of the same shape (chamfered or isometric,
whichever the underlying recipe uses), differing in one dimension (fewer
nodes, a rerouted edge/strut, a color change on the changed node). Never
redraw the whole diagram from scratch for a one-node change — dim the
unchanged nodes.

**Timeline:**
Single horizontal line, chamfered (not circular) milestone markers, labels
alternating above/below to avoid collision. Use the amber/pharos-family
status color for "you are here." Stays flat — a timeline has no structural
hierarchy for depth to communicate.

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
