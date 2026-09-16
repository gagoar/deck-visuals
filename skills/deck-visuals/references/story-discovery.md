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

Five questions, asked together or in a short back-and-forth — not five separate
rounds.

### Brings an outline or draft → extract

Read the material and infer the concept map directly: per-section claims, the
narrative beat each section plays, the likely audience. Ask only where a claim or
beat is genuinely unclear from the material — do not re-ask what the outline
already answers.

## Output: the concept map

One row per slide or section:

| Field | What it captures |
|---|---|
| `claim` | The one idea this slide has to land — not the copy, the concept. |
| `beat` | Its narrative or emotional role: setup, tension, turn, resolution, ask. |
| `audience` | Internal / external, and which of `assets/familiar-sources.md` this room shares (see below). |
| `visual-intent` | chart / diagram / icon / humor / none. |

### Worked example

A short internal pitch for funding a triage-automation project:

| Slide | claim | beat | audience | visual-intent |
|---|---|---|---|---|
| 1. Opener | "Manual triage doesn't scale past our current volume." | setup / stakes | internal | icon |
| 3. Feature comparison | "Too many features isn't always fantastic — ours does one thing well." | tension, self-aware | internal | humor |
| 5. Reliability | "Moving fast with no brakes is how the last rollout broke prod." | tension / warning | internal | humor |
| 7. Architecture | "Three services coordinate through one queue." | explanation | internal | diagram |
| 9. Ask | "Fund phase two." | resolution / ask | internal | chart |

Rows 3 and 5 are exactly the two claims `references/levity.md` walks through as its
analogy-mode worked examples — this map is what hands those claims to that step.

## Audience familiarity, layered per deck

`assets/familiar-sources.md` is the user's own fluency — a fixed, checked-in list he
curates by PR. The interview's audience question can additionally capture, per
deck, **which of those sources this specific room shares** (a team that's watched
the same shows together vs. a client that hasn't). That's a note carried in the
concept map's `audience` field for this one deck — it does not change the checked-in
file. Use it downstream in `levity.md`'s audience-fit gate alongside each source's
`travels-externally` flag.
