# Levity — the GIF slide pattern

A GIF slide is a tension release, not decoration. Use it sparingly and on purpose.

## Ideation — from claim to reference

The concept map from `references/story-discovery.md` is the input here. Every
levity slot starts from that map's `claim` for the slide, not from a blank search.

### Analogy mode (primary)

Do not start from "which GIF." Start from the slide's **claim** — the concept or
argument it's making, not the feeling it evokes. Method:

1. **Name the concept's shape.** Strip the claim to the pattern underneath it:
   excess, momentum-without-control, fragility, over-engineering, recurrence,
   coordination, slow drift, denial of consequence, cascading failure, and so on.
2. **Match a specific famous scene that embodies that shape.** Not a genre, not a
   vague "something funny" — one named scene from one named source.
3. **Prefer `assets/familiar-sources.md` first.** Read it before reaching anywhere
   else; it is the user's own curated fluency, and a match from it is more specific
   and lands harder than a generic reference. Be specific about the scene, not just
   the franchise — "The Homer," not just "The Simpsons."
4. **If nothing on that list fits, walk the fallback chain.** In order:
   `assets/familiar-sources.md` (above) → the **generic evergreen** rows in the seed
   library below → `assets/meme-library.md` (see "Meme fallback") → reaction mode →
   skip. Say "no familiar source fits" before stepping past the profile, and never
   reach for an unlisted franchise the user or audience may not know — a guessed
   reference reads worse than none.

Two worked examples, both from claims mapped by `story-discovery.md`:

- Claim: *"too many features isn't always fantastic."* Shape: excess /
  over-engineering. Match in `familiar-sources.md`: The Simpsons, **"The
  Homer"** — the car with every conceivable feature bolted on, unsellable as a
  result. Specific, on-list, lands the point without a word of explanation.
- Claim: *"moving fast with no brakes."* Shape: momentum without control. Match in
  `familiar-sources.md`: two candidates on the same list — **Wile E. Coyote**
  running off the cliff and only falling once he looks down, or a **Formula 1**
  crash. Either embodies the shape; pick by audience (Coyote travels further
  externally per the file's flags, F1 crash reads faster to a technical or
  motorsport-aware room).

The method is not limited to the seed library or to `familiar-sources.md` alone —
it generates new matches per claim — but it always checks `familiar-sources.md`
first and prefers what's there.

### Reaction mode (secondary)

Use this when the point of the slide is the **audience's feeling**, not a concept —
a pure release moment where there's no claim to embody, just a mood to land: relief,
dread, triumph, overwhelm. Start from the **beat** — the single feeling the audience
should have — and map it to a widely-legible reference the same way as before.

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

Prefer evergreen references over whatever is trending this month: a hyper-current
meme dates the deck the moment the trend passes, and an external audience often
hasn't seen it at all.

### Shared rules for both modes

**Audience fit gates the pool.** For an internal team deck, niche and edgy team
jokes are fair game. For anything external — a client, an exec outside the team —
cut anything niche or edgy, drop any `familiar-sources.md` entry flagged
`travels-externally: limited` or `no`, and stay with the most legible archetypes.
When in doubt about whether a reference travels, treat the deck as external.

**Offer options, not a verdict.** For each levity slot, propose 2-3 candidate
references with the shape/beat and search term for each. The human picks the one
with the right taste and timing — Claude proposes, the human approves. Never lock
in a single GIF choice unasked.

**Humor lives in the imagery, not the prose.** The caption names the joke in one
short line; the slide's own copy stays serious and still routes through `iceberg`
unchanged. Never rewrite slide text to be jokey so it can carry a reference — the
image carries the joke, the words carry the argument.

Once a reference is chosen, source and verify the actual URL below before it goes
anywhere near a fragment.

## Seed library — evergreen concept → scene analogies

A starter set for analogy mode, and a worked example of the method above. Entries
marked **familiar-sources** are on the user's checked-in list and should be reached
for first; entries marked **generic evergreen** are broadly legible references
outside that list, useful when nothing on the user's list fits but a very
widely-known scene still does. The method is not limited to this table — it
generates new matches per claim the same way these were built.

| Concept shape | Example claim | Scene analogy | Source |
|---|---|---|---|
| Excess / feature bloat | "Too many features isn't always fantastic." | The Simpsons — "The Homer" car | familiar-sources |
| Momentum without control | "Moving fast with no brakes." | Wile E. Coyote running off the cliff | familiar-sources |
| Momentum without control (alt) | "We shipped too fast and it broke in prod." | Formula 1 crash | familiar-sources |
| Denial of consequence | "We knew and shipped anyway." | Coyote looking down after he's already past the edge | familiar-sources |
| Coordination breakdown under pressure | "The release process fell apart live." | Formula 1 pit-stop chaos | familiar-sources |
| Avoidance / retreat | "We backed away instead of fixing it." | Homer reversing into the hedge | familiar-sources |
| Small mistake, big visible mess | "One skipped step and the whole room noticed." | The Office — Kevin's chili on the floor | familiar-sources |
| Declaring a fix that isn't one | "We said it was handled and it wasn't." | The Office — Michael Scott "I declare bankruptcy" | familiar-sources |
| Over-engineering | "We built ten times what was needed." | Rube Goldberg machine | generic evergreen |
| Fragility | "One dependency down and it all falls." | House of cards | generic evergreen |
| Recurring bugs | "We fix it and it comes back." | Whack-a-mole | generic evergreen |
| Coordination difficulty | "Getting five teams to agree." | Herding cats | generic evergreen |
| Slow drift | "Nobody noticed until it was too late." | Boiling frog | generic evergreen |
| Precarious balance | "One more addition and it tips over." | A Jenga tower mid-collapse | generic evergreen |
| Diminishing returns / breaking point | "One more ask and it's too much." | The straw that broke the camel's back | generic evergreen |
| Reinventing effort | "We built this again from scratch." | Reinventing the wheel | generic evergreen |
| Unseen scale | "The real problem is under the surface." | Tip of the iceberg | generic evergreen |
| Early warning ignored | "We had the signal and ignored it." | Canary in the coal mine | generic evergreen |
| Cascading failure | "One part fails and the rest follows." | Toppling dominoes | generic evergreen |
| Chaotic parallel effort | "Everyone trying everything at once." | Throwing spaghetti at the wall | generic evergreen |
| Accidental complexity | "The diagram looks like a plate of noodles." | A spaghetti diagram | generic evergreen |
| Fragmented understanding | "Everyone describes a different system." | The blind men and the elephant | generic evergreen |
| Shifting requirements | "The target keeps moving." | Moving the goalposts | generic evergreen |
| Obvious problem nobody names | "Nobody wants to say it out loud." | The elephant in the room | generic evergreen |
| Persistent thankless effort | "We push the same rock uphill every quarter." | Sisyphus and the boulder | generic evergreen |

If a claim's concept doesn't map cleanly to anything above, run the method fresh —
name the shape, then search for a scene that fits it — but check
`familiar-sources.md` first, and say "no familiar source fits" rather than guessing
at a franchise outside it.

## Meme fallback

When no familiar source and no generic-evergreen row fits, but the slide still wants
an image, reach for `assets/meme-library.md` — a small pool of broadly-legible
internet references (this-is-fine dog, distracted boyfriend, toppling dominoes). It
is ungated by recognition, unlike `assets/familiar-sources.md`: these land without a
shared back-catalog. It stores **search terms, not URLs**, so source and verify the
direct media URL at use exactly as below. Apply the same audience-fit gate — drop any
entry flagged `travels-externally: limited` or `no` for an external deck. Prefer a
familiar-source or evergreen match when one exists; the meme library is the step
before reaction mode, not the first reach.

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
