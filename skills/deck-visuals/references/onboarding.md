# Onboarding — building the presenter profile by swipe

`assets/familiar-sources.md` is the pool the analogy matcher reads first. This flow
rebuilds it from a fast visual pass instead of by hand: the presenter swipes through
`assets/show-catalog.md`, marks what he recognizes, states the talk length, and
describes the room in his own words. Recognizing a show means he gets the jokes
around it, so it earns a place in the profile.

**Boundary.** Onboarding maintains the one profile. It does not build decks, write
copy, or place visuals — that is the pipeline in `SKILL.md`. Run it out-of-band, on
its own, whenever the profile needs a refresh.

**When to run.** On "onboard me", "set up / update my familiar sources", "build my
pop-culture profile" — independent of any deck.

## The flow

1. **Load the catalog.** Read `assets/show-catalog.md`. Each entry is one card.
2. **Source and verify media.** For each entry, `WebSearch` a still and the signature
   GIF (use the baked `GIF search`; the still defaults to `"<title> still"`). Resolve
   to **direct** media URLs per `references/levity.md`'s rules — never a guessed or
   hand-built URL. Verify in batches with
   `assets/tools/dv-tools check-media-url <url> …`; only PASS URLs are used. This
   is the same gate levity uses, no new verifier.
   - **Prefer a small rendition.** Use a lightweight variant (e.g. Giphy's
     `200w.gif`), not the original `giphy.gif` — originals can be tens of MB and read
     as an empty card while they load. A swipe deck of ~50 needs fast tiles.
   - **Degrade gracefully.** GIF fails → image-only card. Both fail → `degraded: true`
     title-only card. Never write an unverified URL. Write the verified set to a
     `cards.json` in an OS temp dir — that file is the per-run cache.
3. **Launch the binary.** Pick the prebuilt binary matching the platform
   (`uname -s` / `uname -m`) from `assets/onboarding/bin/`, e.g.
   `dv-onboard-darwin-arm64`. Run it:
   `dv-onboard --form onboard --data <cards.json> --out <results.json> --port 0`. It
   needs no Go and no interpreter. Hand the user the `http://127.0.0.1:<port>/` it
   prints. (The same binary serves the deck-time levity picker with
   `--form levity` — see `references/levity.md`.)
4. **The user swipes.** Yes / No / Skip per card, then suggests anything missed. The
   binary takes one POST, writes `results.json`, and shuts down on its own.
5. **Read results.** Parse `results.json` (shape below).
6. **Merge into the profile.** Rewrite `assets/familiar-sources.md` — merge, preserve,
   flag (see below). Keep its exact format.
7. **Check coverage.** Run
   `assets/tools/dv-tools coverage --levity references/levity.md
   --catalog assets/show-catalog.md --results <results.json>`. Relay any GAP and name
   the catalog titles that would close it.
8. **Review suggestions.** Walk the suggestions with the user one by one; promote only
   the confirmed ones.

## results.json

```json
{
  "submitted_at": "…",
  "talk": { "duration_minutes": 30 },
  "audience": { "description": "verbatim free text, uncategorized" },
  "recognized_ids": ["the-office"],
  "not_recognized_ids": ["formula-1"],
  "skipped_ids": [],
  "suggestions": [ { "title": "The Bear", "note": "kitchen chaos" } ]
}
```

## Merge, preserve, flag

- **Each `recognized_id` → a profile entry** in `familiar-sources.md`'s exact format:
  `**Title**` under its `## Domain`, the catalog's `Scenes` bullets (already in
  `scene — the concept it embodies` shape), and a `Travels externally: <flag> —
  reason` line seeded from the catalog's `Travels externally` default.
- **Preserve hand-tuned entries.** If a recognized show is already in the file with a
  hand-tuned `travels-externally` flag or edited scenes, keep them — do not clobber a
  flag the user tested on a real audience.
- **Flag removals, never delete silently.** A show currently in the file but *not*
  recognized this round is surfaced for the user's review, not dropped on its own.
- **Suggestions go to a review queue.** Never blind-write a suggested title. Each
  needs a confirmed scene bullet and travels flag before it earns a place — promoted
  one at a time.

## Talk length → cadence

Record `duration_minutes` and carry it into the concept map. A short slot gets fewer
levity beats; keep the rough cadence in `references/levity.md` (openers, section
breaks, the payoff — never every slide).

## Do not stereotype the audience

The picker captures the audience as free text on purpose. Take it **verbatim**. Do
not infer taste, technical level, or seniority from it, and do not map it to an
audience "type". Use it only as the user's own words feeding the audience-fit gate in
`levity.md`; ask a follow-up rather than assume.

## Hard rules

- The binary binds `127.0.0.1` only, uses no auth, and makes no outbound calls.
  Nothing leaves the machine.
- Never write an unverified media URL into `cards.json`.
- Never blind-delete a profile entry or blind-write a suggestion.
- The picker (`assets/onboarding/picker.html`, embedded in the binary) is local
  tooling — never injected into a deck or delivered.
