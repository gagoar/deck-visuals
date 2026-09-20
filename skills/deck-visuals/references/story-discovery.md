# Story discovery — the intake that feeds every visual

Before any chart, diagram, icon, or GIF gets placed, this step builds a **concept
map**. Every later step in `SKILL.md` reads that map — it is the single input that
positions all visual placement, humor included, not just levity.

**Boundary.** Discovery *elicits* the concept map to place visuals. It does not
write the story, the argument, or the slide copy — that stays with the user,
`iceberg`, and `frontend-slides`. If the user wants the narrative itself drafted or
rewritten, hand that back; this step only extracts or asks for structure.

## Adaptive intake

Pick the path based on what the user brings.

### No material yet → guided interview

One short pass, not an interrogation. Ask:

1. **Core argument.** What's the one thing this deck has to convince someone of?
   (The conceptual arch.)
2. **Narrative arc.** Setup → tension → resolution/ask — what's the tension, and
   what's the ask at the end?
3. **Audience.** Internal team, or external (client, exec, someone outside the
   team)? This gates humor and reference legibility downstream.
4. **Section-by-section claims.** For each planned section or slide, what's the one
   claim it needs to land?
5. **Energy and humor.** Where does the user want a lift — an opener, a section
   break, the payoff — and where should the deck stay serious throughout?
6. **Subject world (visual through-line).** What is this deck *about*, and does the
   subject suggest a visual metaphor the diagrams can speak in? A space deck →
   orbits / a solar system; an architecture deck → blueprints, angles, isometric
   plates; logistics → route maps. **Ask about the subject, never about colors** (see
   "Never name colors" below).

Six questions, asked together or in a short back-and-forth — not six separate rounds.

### Two-tier idiom capture (question 6)

- **If the user already pictures the visuals a certain way, use their mapping
  verbatim** — "nodes as planets on orbital rings", "layers as blueprint plates".
  A user-supplied idiom always wins.
- **If they don't, don't force it.** Infer 1-2 candidate idioms from the subject and
  offer them for a quick yes/no; if none fit, fall back to the clean house idiom (no
  metaphor). A plain, metaphor-free deck is a fine outcome.

### Brings an outline or draft → extract

Read the material and infer the concept map directly: per-section claims, the
narrative beat each section plays, the likely audience, and the **subject world /
visual through-line** (from the title, domain, and recurring nouns). Ask only where a
claim, beat, or the idiom is genuinely unclear from the material — do not re-ask what
the outline already answers.

## Output: the concept map

It has a **deck-level header** and then one row per slide or section.

**Deck-level header** (carried once, applied to every visual):

| Field | What it captures |
|---|---|
| `visual through-line` | The subject idiom + a one-line motif — e.g. "space → nodes as planets on orbital rings; edges as orbits". Empty = clean house idiom. From question 6; the user-supplied mapping wins over an inferred one. |
| `palette` | The chosen palette preset id, if the user picked one from shown swatches (see "Never name colors"); default = the house palette. Never a color name. |

**Per-slide rows:**

| Field | What it captures |
|---|---|
| `claim` | The one idea this slide has to land — not the copy, the concept. |
| `beat` | Its narrative or emotional role: setup, tension, turn, resolution, ask. |
| `story-role` | How this visual *advances the arc*: sets the stakes, embodies the tension, marks the turn, lands the ask. Distinct from `claim` (what) and `visual-intent` (which flavor). |
| `audience` | Internal / external, and which of `assets/familiar-sources.md` this room shares (see below). |
| `visual-intent` | chart / diagram / icon / humor / none. |

### Worked example

A short internal pitch for funding a triage-automation project.

**Deck header** — `visual through-line`: *"triage as an ER — intake, sorting, and a
queue of patients"* (subject-driven idiom); `palette`: house default.

| Slide | claim | beat | story-role | visual-intent |
|---|---|---|---|---|
| 1. Opener | "Manual triage doesn't scale past our current volume." | setup / stakes | sets the stakes | icon |
| 3. Feature comparison | "Too many features isn't always fantastic — ours does one thing well." | tension, self-aware | embodies the tension | humor |
| 5. Reliability | "Moving fast with no brakes is how the last rollout broke prod." | tension / warning | embodies the tension | humor |
| 7. Architecture | "Three services coordinate through one queue." | explanation | shows how it works, in the idiom | diagram |
| 9. Ask | "Fund phase two." | resolution / ask | lands the ask | chart |

(`audience` is internal for every row here; it stays a per-row field, omitted from
this table for width.) Row 7's diagram is drawn in the deck's through-line — the queue
as an ER triage line — not a generic box-and-arrow. Rows 3 and 5 are exactly the two
claims `references/levity.md` walks through as its analogy-mode worked examples — this
map is what hands those claims to that step.

## Never name colors

Do not ask the user to name or describe colors — people rarely know their palette or
how to name it. The intake asks about the **subject world** (question 6), not colors.
The default is the validated house palette (`assets/brand-palette.md`). If a per-deck
palette variant fits the subject, it is chosen by **showing swatches to pick from**
(the `dv-onboard --form palette` chooser, presets in `assets/palette-presets.md`),
and only the picked preset id is recorded on the concept-map header — never a color
name.

## Audience familiarity, layered per deck

`assets/familiar-sources.md` is the user's own fluency — a fixed, checked-in list he
curates by PR. The interview's audience question can additionally capture, per
deck, **which of those sources this specific room shares** (a team that's watched
the same shows together vs. a client that hasn't). That's a note carried in the
concept map's `audience` field for this one deck — it does not change the checked-in
file. Use it downstream in `levity.md`'s audience-fit gate alongside each source's
`travels-externally` flag.
