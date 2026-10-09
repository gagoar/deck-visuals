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
   fallback on every one. A swipe onboarding flow rebuilds the presenter's
   pop-culture profile so the matcher always has options, backed by a meme library
   for when nothing else fits.

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
                              (incl. onboarding.md)
  assets/
    brand-palette.md          validated 8-slot categorical palette + ramps
    show-catalog.md           the ~50-title swipe catalog (onboarding)
    meme-library.md           evergreen meme fallback for levity
    familiar-sources.md       the presenter's curated analogy pool
    VERSION                   release tag the wrappers download from
    run-binary.sh             shared fetch + verify + cache + build-fallback
    tools/                    dv-tools: Go source for the palette validator,
                              media-URL check, and coverage report
      dv-tools                POSIX wrapper: fetches the pinned binary for this
                              machine on first use, then runs it
      SHA256SUMS              pinned hashes of the released dv-tools binaries
      build.sh                reproducible cross-compile; rewrites SHA256SUMS
    icons/                    curated Lucide-style icon set
    fragments/                injectable HTML fragments
    onboarding/
      picker.html             onboarding swipe page (local tool)
      levity-picker.html      deck-time levity GIF chooser (local tool)
      picker.css              shared premium dark styling for both
      dv-onboard              POSIX wrapper for the picker binary (same scheme)
      SHA256SUMS              pinned hashes of the released dv-onboard binaries
      server/                 Go source for the picker binary
```

The same picker server also powers a **deck-time levity picker**: once Claude
proposes 2-3 candidate GIFs per slide, run it with `--form levity` to choose each one
visually instead of in chat.

## Helper binaries

`dv-tools` (palette validator, media-URL check, coverage) and `dv-onboard` (the picker
server) are Go programs. Their binaries are **not committed**. The wrappers
(`assets/tools/dv-tools`, `assets/onboarding/dv-onboard`) download the binary for the
installed release on first use, run it only if it matches the pinned `SHA256SUMS`, and
cache it. They need **network access to GitHub, or Go installed** the first time; after
that they work offline. Maintainers: bump `assets/VERSION`, `.claude-plugin/plugin.json`
and the `SHA256SUMS` files (run both `build.sh` scripts) together, then push a `v<version>`
tag so `release-binaries.yml` publishes the assets.

## Onboard your references

```
"Onboard me" / "update my familiar sources"
```

Swipe through ~50 shows in a local page — each a still plus a signature GIF —
marking what you recognize, and give your talk length and audience. The picks feed
back through a small local-only binary (nothing leaves the machine) and rebuild
`familiar-sources.md`, with a coverage check that proves the matcher has ~3 options
per concept shape.

## License

MIT. See `LICENSE`. The icon set carries its own provenance note and license
text in `skills/deck-visuals/assets/icons/LICENSE`.

## Automated gate

CI runs `scripts/check-plugin-safety.ts` and `claude plugin validate --strict` on every PR. Run both before you push:

```
node scripts/check-plugin-safety.ts .
claude plugin validate . --strict
```

The gate fails on three things the plugin-directory upload scanner holds for review:

- A skill, agent, or command that pre-approves `Write`, `Edit`, or `Bash` with no path scope. Write `Write(./docs/**)`, not `Write`.
- Shipped code that copies the whole environment (`...process.env`, `Object.keys(process.env)`). Read the named variables you need.
- A committed native executable. The scanner cannot read it.
