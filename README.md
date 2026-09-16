# deck-visuals

A Claude Code plugin that decides where a visual belongs in an HTML deck and
which kind, then hands back a ready-to-inject fragment. It does not render
decks.

## The three pillars

1. **Impactful visuals.** Inline-SVG concept diagrams for relationships, flows,
   and architecture. Icon anchors for single ideas. AI illustrations for hero
   moments are deferred — this build has no image-generation dependency to
   satisfy that case yet.
2. **Data-viz principles.** Charts lean on the built-in `dataviz` skill for form,
   marks, and interaction, paired with a validated brand palette instead of
   `dataviz`'s brand-neutral default. Every chart ships with a reachable table
   view.
3. **Levity.** Real GIFs by direct media URL, at a deliberate cadence — openers,
   section breaks, the payoff — with alt text and a `prefers-reduced-motion`
   fallback on every one.

## Guidance and injection, not rendering

`frontend-slides` renders decks: the HTML shell, the fixed stage, the theme, the
animations. `deck-visuals` never touches that layer. It reads a deck's content,
decides where a visual belongs, builds that one visual as a self-contained
fragment with its own scoped styles, and hands it back for injection into the
slide that's already there. `iceberg` still owns the copy pass, run on prose
before it enters the deck — never on the finished `.html` file.

See `skills/deck-visuals/SKILL.md` for the full six-step pipeline and
`skills/deck-visuals/references/` for the per-flavor rules.

## Install

```
/plugin marketplace add gagoar/gago-plugins
/plugin install deck-visuals@gago-plugins
/reload-plugins
```

## Usage

Ask in natural language once a deck exists:

- "Add visuals to my deck."
- "Make this slide impactful."
- "Chart this number."
- "Add a GIF/levity slide here."
- "Inject a diagram for this architecture slide."

The skill reads the deck, applies the decision rule in
`skills/deck-visuals/references/visual-concepts.md`, and returns the fragment to
inject — it does not rewrite the deck's shell or theme.

## What's inside

```
skills/deck-visuals/
  SKILL.md                    orchestrator: pipeline + trigger phrases
  references/                 decision rules, per-flavor guidance
  assets/
    brand-palette.md          validated 8-slot categorical palette + ramps
    scripts/                  palette validator (Node + Python)
    icons/                    curated Lucide-style icon set
    fragments/                injectable HTML fragments
```

## License

MIT. See `LICENSE`. The icon set carries its own provenance note and license
text in `skills/deck-visuals/assets/icons/LICENSE`.
