# Levity — the GIF slide pattern

A GIF slide is a tension release, not decoration. Use it sparingly and on purpose.

## Ideation — from beat to reference

Do not start from "which GIF." Start from the **beat** — the single feeling the
audience should have at that moment. Facepalm. Relief. Dread. "Here we go again."
Triumph / mic-drop. Overwhelm / chaos. Skepticism. Waiting / loading. Small-vs-big.
Name the beat first, in one word or phrase, before searching for anything.

Then map the beat to a **reference archetype** — a widely-legible bit that ages
slowly and travels across audiences. Prefer evergreen references over whatever is
trending this month: a hyper-current meme dates the deck the moment the trend
passes, and an external audience often hasn't seen it at all. Evergreen sources
that hold up: The Office, Parks and Rec, The Simpsons, the "this is fine" dog,
Drake, distracted boyfriend, Spider-Man pointing at Spider-Man, Picard facepalm,
Kermit sipping tea.

| Beat | Example archetype | Search term to try |
|---|---|---|
| Facepalm | Picard facepalm | `picard facepalm giphy` |
| Relief | Michael Scott exhale / "thank god" | `michael scott relief giphy` |
| Dread / here we go again | The Office "here we go again" | `the office here we go again giphy` |
| Triumph / mic-drop | Ron Swanson approval, or a literal mic drop | `ron swanson nod giphy` |
| Overwhelm / chaos | Kermit arms flailing, or the "this is fine" dog | `this is fine dog giphy` |
| Skepticism | Kermit sipping tea ("but that's none of my business") | `kermit sipping tea giphy` |
| Waiting / loading | The Simpsons "waiting" bench scene | `simpsons waiting bench giphy` |
| Small-vs-big / mismatch | Distracted boyfriend | `distracted boyfriend giphy` |

**Audience fit gates the pool.** For an internal team deck, niche and edgy team
jokes are fair game. For anything external — a client, an exec outside the team —
cut anything niche or edgy and stay with the most legible archetypes on the list
above (Picard, Kermit, The Office). When in doubt about whether a reference
travels, treat the deck as external.

**Offer options, not a verdict.** For each levity slot, propose 2-3 candidate
references with the beat and search term for each. The human picks the one with
the right taste and timing — Claude proposes, the human approves. Never lock in a
single GIF choice unasked.

Once a reference is chosen, source and verify the actual URL below before it goes
anywhere near a fragment.

## Real GIFs, by external URL

Pin the **direct media URL**, not the page URL.

- Giphy: use the `media*.giphy.com/media/<id>/giphy.gif` form, not
  `giphy.com/gifs/<slug>`. The page URL renders a webpage, not an image.
- Tenor: use the `media.tenor.com/.../....gif` (or `.mp4`) form from the "Copy
  direct link," not the `tenor.com/view/...` page.
- Verify the direct link actually returns image bytes before shipping it. A page
  URL dropped into an `<img src>` renders broken.

## Sourcing — WebSearch + verify

Once a reference is chosen (see Ideation above), get a real, resolving direct
media URL using the built-in `WebSearch` and `WebFetch` tools. A made-up or
guessed URL is not a shortcut — it returns 403 or 404 and ships a broken image.
Follow every step; do not skip the verify step because a URL "looks right."

1. **Search.** Run `WebSearch` for the reference plus `giphy` (or `tenor`), e.g.
   `"this is fine" dog giphy`. Prefer a result on `giphy.com` or `tenor.com`
   itself over a third-party aggregator.
2. **Resolve to a direct media URL.** Open a candidate result. If the result is
   already a direct media URL (matches the forms in the section above), use it.
   Otherwise `WebFetch` the page and extract the direct URL:
   - **Giphy:** a gif/embed page exposes its direct form as
     `https://media*.giphy.com/media/<id>/giphy.gif` (or `/200.gif` for a smaller
     variant). The `<id>` is the same token in the page's `giphy.com/gifs/<slug>-<id>`
     URL — take the trailing ID segment and drop it into the media-host template.
   - **Tenor:** a view page (`tenor.com/view/...-<id>`) exposes its direct form as
     `https://media*.tenor.com/.../....gif` or `.mp4`. `WebFetch` the page and ask
     for the `media.tenor.com` URL in the page source or the "Copy direct link"
     target — do not hand-construct it from the view-page slug, Tenor's media
     paths are not derivable from the slug alone.
3. **Verify before embedding.** Never embed a URL that hasn't been checked. Run
   the checker: `python3 assets/scripts/check_media_url.py <url>`. It requires
   HTTP 200 and a `Content-Type` of `image/*` or `video/*` and a non-trivial body
   size, and exits 1 if any URL fails. No script handy? `curl -sIL <url>` and
   confirm the same three things by eye. If it fails, go back to step 1 — do not
   ship the broken URL.
4. **Taste-check the actual content.** Look at what the verified GIF shows and
   confirm it matches the intended beat and reads as appropriate — no NSFW.
   `WebSearch` results are not rating-filtered; this judgment call is on the
   human or Claude, not the search tool.
5. **Record the URL and the search term.** Put the verified direct URL in the
   fragment's `src` (see `assets/fragments/levity-slide.html`), and keep the
   search term from step 1 in an HTML comment next to it. When the URL link-rots
   (see Link-rot note below), that comment is the fastest way back to a
   replacement — re-run this procedure from step 1 with the same term.

This is the only sourcing mechanism this skill uses. Do not fabricate a Giphy/Tenor
ID or guess at a slug-to-URL mapping — every direct URL in a shipped deck must
have come from a search result that was fetched and verified, not assembled from
memory.

## Cadence

Default cadence, tune per deck length:

- **Openers.** One GIF on the title or first-section slide to set tone.
- **Section breaks.** One GIF between major sections, not between every slide.
- **The payoff.** One GIF at the climax or ask — the slide the audience should
  remember.

**Never every slide.** A GIF on every slide reads as noise, not levity. If a deck
has 12 slides, 2-3 GIF slides is the right order of magnitude, not 6.

## Accessibility

- **Alt text is required**, not optional. Describe what the GIF shows and why it's
  there in one sentence — "Trombone dropping, for the failed migration."
- **`prefers-reduced-motion` must pause or replace the GIF with a still.** Use the
  poster-frame pattern: an `<img>` with the GIF as `src` and a static frame swapped
  in via `@media (prefers-reduced-motion: reduce)`, or pause a `<video>` tag if
  using the `.mp4` form. Never leave an animating GIF running for a user who asked
  for reduced motion.
- Caption the joke in one short line under the GIF. Do not rely on the audience
  getting the reference from the image alone.

## Taste rules

- Match the reaction to the actual content. A stretch reference reads worse than
  no GIF at all.
- Punch up or punch at the problem, never at a named person in the room.
- One clean reaction beats a chain of three GIFs making the same joke.
- If the deck is going external (a client, a exec outside the team), cut the
  cadence in half and raise the taste bar — internal-team humor does not always
  travel.

## Link-rot note

External GIF hosts remove content. Before delivery, re-check every pinned URL
still resolves — run `check_media_url.py` again on every URL in the deck. If a
link has gone dead, re-source it: the HTML comment left by step 5 of Sourcing
above has the original search term, so re-sourcing starts from step 1, not from
scratch. Do not leave a broken image in a shipped deck. This is a real, recurring
failure mode of external-URL embeds; budget for it at delivery time, not just at
build time.
