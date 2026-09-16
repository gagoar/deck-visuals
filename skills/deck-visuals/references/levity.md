# Levity — the GIF slide pattern

A GIF slide is a tension release, not decoration. Use it sparingly and on purpose.

## Real GIFs, by external URL

Pin the **direct media URL**, not the page URL.

- Giphy: use the `media*.giphy.com/media/<id>/giphy.gif` form, not
  `giphy.com/gifs/<slug>`. The page URL renders a webpage, not an image.
- Tenor: use the `media.tenor.com/.../....gif` (or `.mp4`) form from the "Copy
  direct link," not the `tenor.com/view/...` page.
- Verify the direct link actually returns image bytes before shipping it. A page
  URL dropped into an `<img src>` renders broken.

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
still resolves. If a link has gone dead, replace it — do not leave a broken image
in a shipped deck. This is a real, recurring failure mode of external-URL embeds;
budget for it at delivery time, not just at build time.
