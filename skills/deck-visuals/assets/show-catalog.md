# Show catalog — the swipe deck for onboarding

This file is **config, not documentation**. It is the fixed pool the onboarding
picker swipes through (see `references/onboarding.md`). Each entry becomes one card:
a still image plus a signature GIF, sourced and verified at generation time, with
the prompt "do you recognize this show?" A *yes* means the presenter gets the jokes
around it, so the show earns a place in `assets/familiar-sources.md`.

**This is a recognition survey, not the meme fallback.** Broadly-legible internet
memes live in `assets/meme-library.md` — an ungated, always-on fallback the levity
step reaches for when no familiar source fits. They are not surveyed here.

## Format

Domain is the `##` heading. Every entry carries parseable fields the picker and the
scripts read:

- `id:` — kebab-case, stable, unique. The picker and `results.json` key on it.
- `Shapes:` — semicolon-separated concept shapes this show can carry. **Use the
  exact vocabulary of the "Concept shape" column in `references/levity.md`'s seed
  library** — that coupling is what `assets/tools/dv-tools coverage` measures.
- `Scenes:` — semicolon-separated `scene — the concept it embodies` bullets, in the
  same shape as `familiar-sources.md`. These are copied verbatim when a recognized
  show is written into the profile, so mapping is deterministic, not re-invented.
- `GIF search:` — the baked query for the card's signature GIF.
- `Image search:` — optional; defaults to `"<title> still"` when omitted.
- `Travels externally:` — `yes` / `limited` / `no` and a one-line reason, seeded as
  the default `travels-externally` flag on any profile entry this show produces.

Curate by PR: add a title once it earns a place, remove one that has gone stale,
retune a `Travels externally` default once tested on a real audience. Keep the union
of `Shapes:` tags covering each `levity.md` seed shape at least three times so the
matcher always has options — `dv-tools coverage` checks this.

## TV comedy

**The Office**
- id: the-office
- Shapes: Small mistake, big visible mess; Declaring a fix that isn't one; Avoidance / retreat
- Scenes: Kevin's chili spilled across the floor — one small mistake becoming a visible, can't-look-away disaster in front of the whole room; Michael yelling "I declare bankruptcy" — declaring a problem solved by saying the words without doing the work; Michael driving into the lake because the GPS said to — trusting a broken process past the point it obviously stopped working
- GIF search: the office kevin chili floor giphy
- Travels externally: yes — long-running global syndication, broadly recognized

**Friends**
- id: friends
- Shapes: Momentum without control; Obvious problem nobody names; Unseen scale
- Scenes: Ross shouting "PIVOT!" wedging a couch up a stairwell — brute-forcing a change that ignores the constraints; "We were on a break!" — an argument that never ends because nobody wrote down what was agreed; Monica's one junk-stuffed closet behind a spotless apartment — a tidy surface with the mess pushed one door out of sight
- GIF search: friends ross pivot couch giphy
- Travels externally: yes — long-running global syndication, broadly recognized

**It's Always Sunny in Philadelphia**
- id: always-sunny
- Shapes: Accidental complexity; Fragmented understanding
- Scenes: Charlie at the "Pepe Silvia" mail wall — string, papers and arrows everywhere, frantically explaining a system nobody else can follow
- GIF search: charlie pepe silvia conspiracy board giphy
- Travels externally: yes — the conspiracy-board scene is a near-universal internet reference

**Parks and Recreation**
- id: parks-and-rec
- Shapes: Persistent thankless effort; Coordination difficulty
- Scenes: Leslie pushing a parks project through committee after committee — pushing the same rock uphill against an indifferent system; the whole team wrangled toward one goal none of them started aligned on — herding people who each want something different
- GIF search: parks and rec leslie knope giphy
- Travels externally: limited — well loved but less universally quoted outside the US

**Brooklyn Nine-Nine**
- id: brooklyn-99
- Shapes: Coordination difficulty; Recurring bugs
- Scenes: the precinct's endless Halloween heist — the same problem staged again every year with new rules; "cool cool cool, no doubt no doubt" — papering over a coordination gap with agreement nobody means
- GIF search: brooklyn 99 cool cool cool giphy
- Travels externally: yes — broad streaming recognition

**Arrested Development**
- id: arrested-development
- Shapes: Accidental complexity; Recurring bugs
- Scenes: "I've made a huge mistake" — a plan that folds the instant it meets reality; the family's tangle of schemes propping up other schemes — accidental complexity nobody can unwind
- GIF search: arrested development huge mistake giphy
- Travels externally: limited — a cult favorite, quoted mostly by fans

**Seinfeld**
- id: seinfeld
- Shapes: Obvious problem nobody names; Diminishing returns / breaking point
- Scenes: George's web of small lies collapsing at once — the obvious thing nobody would say out loud finally landing; "no soup for you!" — one more demand and the whole arrangement is off
- GIF search: seinfeld no soup for you giphy
- Travels externally: yes — decades of syndication, widely quoted

**The Good Place**
- id: the-good-place
- Shapes: Shifting requirements; Reinventing effort
- Scenes: the neighborhood rebooted for the hundredth time — the target moved so often nobody remembers the original spec; the point system re-derived from scratch when the old one turns out broken — rebuilding the rules rather than patching them
- GIF search: the good place reboot giphy
- Travels externally: yes — broad streaming recognition

**Community**
- id: community
- Shapes: Chaotic parallel effort; Over-engineering
- Scenes: the paintball war escalating into an all-out campaign — everyone trying everything at once until it's a full production; a simple study group turning into an elaborate multi-front operation — far more machinery than the task needed
- GIF search: community paintball giphy
- Travels externally: limited — a cult favorite with meme-level scenes

**Curb Your Enthusiasm**
- id: curb-your-enthusiasm
- Shapes: Small mistake, big visible mess; Obvious problem nobody names
- Scenes: Larry breaking an unspoken social rule and the room freezing — a tiny misstep becoming a visible standoff; Larry saying the thing everyone was thinking and nobody would say — the obvious problem named out loud, to everyone's discomfort
- GIF search: curb your enthusiasm larry david giphy
- Travels externally: limited — the "curb" freeze-frame meme travels further than the show

## TV drama

**Breaking Bad**
- id: breaking-bad
- Shapes: Momentum without control; Denial of consequence; Unseen scale
- Scenes: Gus Fring's spotless laundromat sitting on a hidden industrial lab — a clean front with all the complex machinery deliberately out of view; Walt sure he's in control right up until he isn't — confidence and momentum running well past the brakes; Walt's quiet "Say my name" — the moment the system finally names who owns the thing
- GIF search: breaking bad say my name giphy
- Travels externally: yes — a globally distributed, widely-referenced series

**Game of Thrones**
- id: game-of-thrones
- Shapes: Cascading failure; Early warning ignored; Shifting requirements
- Scenes: "Winter is coming" ignored for years until it arrives all at once — a warning everyone dismissed cascading into disaster; alliances flipping every season — the target moving faster than any plan could hold
- GIF search: game of thrones winter is coming giphy
- Travels externally: yes — near-universal cultural recognition at its peak

**The Wire**
- id: the-wire
- Shapes: Fragmented understanding; Persistent thankless effort
- Scenes: each unit seeing only its slice of the city — everyone describing a different system, none of them the whole; the same corner cleaned up and reoccupied season after season — thankless effort against a problem that regrows
- GIF search: the wire omar giphy
- Travels externally: limited — critically revered, less broadly quoted

**Succession**
- id: succession
- Shapes: Coordination breakdown under pressure; Obvious problem nobody names
- Scenes: the family scrambling as a decision comes due and no one is aligned — a drilled process falling apart under real pressure; the succession question everyone circles and no one will state plainly
- GIF search: succession logan roy giphy
- Travels externally: limited — strong among a prestige-TV audience

**Severance**
- id: severance
- Shapes: Fragmented understanding; Slow drift; Unseen scale
- Scenes: the severed floor with no memory of the outside — a system split so cleanly the two halves can't see each other; the innies slowly realizing how much sits above them — the real scope hidden one floor up
- GIF search: severance innie outie giphy
- Travels externally: limited — prestige streaming, stronger with a newer audience

**Stranger Things**
- id: stranger-things
- Shapes: Early warning ignored; Cascading failure
- Scenes: the small-town signs waved off until the Upside Down breaks through — early warnings ignored until they cascade; one gate opening and pulling the whole town in after it
- GIF search: stranger things upside down giphy
- Travels externally: yes — broad global streaming recognition

**The Sopranos**
- id: the-sopranos
- Shapes: Denial of consequence; Recurring bugs
- Scenes: Tony insisting everything's handled while it quietly isn't — denial running well past the evidence; the same feuds resurfacing no matter how many are "settled" — a problem that keeps coming back
- GIF search: the sopranos tony giphy
- Travels externally: limited — revered, quoted mostly by fans

**Chernobyl**
- id: chernobyl
- Shapes: Denial of consequence; Early warning ignored; Cascading failure
- Scenes: officials insisting the reading is impossible while it climbs — denial in the face of a plain signal; one dismissed warning cascading into total failure; "not great, not terrible" — downplaying a reading that was off the scale
- GIF search: chernobyl not great not terrible giphy
- Travels externally: limited — acclaimed, one meme line travels wider than the show

**Mad Men**
- id: mad-men
- Shapes: Reinventing effort; Slow drift; Avoidance / retreat
- Scenes: Don pitching a whole new identity when the old one stops selling — rebuilding from scratch instead of fixing; Don slipping out when a problem gets close — retreating rather than addressing it
- GIF search: mad men don draper giphy
- Travels externally: limited — prestige recognition, fewer quotable meme moments

**The Last of Us**
- id: the-last-of-us
- Shapes: Fragility; Cascading failure
- Scenes: a single infection breaching a fortified zone — one weak point taking down a whole safe system; the network of quarantine zones falling one after another — a chain reaction once the first link goes
- GIF search: the last of us giphy
- Travels externally: yes — huge game-plus-show crossover recognition

## Animation

**The Simpsons**
- id: the-simpsons
- Shapes: Excess / feature bloat; Avoidance / retreat
- Scenes: "The Homer" — the car with every conceivable feature bolted on, unsellable as a result; Homer slowly reversing into the hedge to avoid a conversation — retreating from a problem instead of addressing it
- GIF search: simpsons homer bushes giphy
- Travels externally: yes — decades of global syndication, near-universal recognition

**Looney Tunes**
- id: looney-tunes
- Shapes: Momentum without control; Declaring a fix that isn't one
- Scenes: Wile E. Coyote running off the cliff and only falling once he looks down — momentum that outran the ground under it; the anvil dropping on Coyote right after his plan looked like it worked — a fix that immediately backfires
- GIF search: wile e coyote cliff giphy
- Travels externally: yes — near-universal recognition across generations

**Futurama**
- id: futurama
- Shapes: Reinventing effort; Diminishing returns / breaking point
- Scenes: Bender's "I'm gonna build my own, with blackjack and…" — walking off to rebuild from scratch instead of reusing; Fry's "shut up and take my money" — demand running ahead of the thing even being ready
- GIF search: futurama shut up and take my money giphy
- Travels externally: yes — widely syndicated with meme-level scene recognition

**Rick and Morty**
- id: rick-and-morty
- Shapes: Over-engineering; Chaotic parallel effort; Excess / feature bloat; Cascading failure
- Scenes: Rick building a galaxy-spanning contraption for a trivial errand — vastly more machinery than the task needed; a portal-hopping scheme spinning off into a dozen parallel disasters at once
- GIF search: rick and morty portal giphy
- Travels externally: yes — strong meme recognition among a broad audience

**SpongeBob SquarePants**
- id: spongebob
- Shapes: Recurring bugs; Persistent thankless effort
- Scenes: the same Krusty Krab crisis flaring up episode after episode — a problem that keeps coming back; SpongeBob flipping patties forever with unbroken cheer — thankless effort that never ends
- GIF search: spongebob krusty krab giphy
- Travels externally: yes — near-universal recognition, meme-heavy

**South Park**
- id: south-park
- Shapes: Over-engineering; Obvious problem nobody names
- Scenes: the underpants gnomes' "collect underpants → ??? → profit" — an elaborate plan with the crucial step missing; the town elaborately avoiding the obvious thing right in front of them
- GIF search: south park underpants gnomes giphy
- Travels externally: yes — broad recognition, one near-universal meme

**Bob's Burgers**
- id: bobs-burgers
- Shapes: Small mistake, big visible mess; Persistent thankless effort
- Scenes: one small kitchen slip spiraling into a visible restaurant disaster; Bob reopening and grinding through the same lunch rush every single day — thankless, repeated effort
- GIF search: bobs burgers giphy
- Travels externally: limited — beloved but less universally quoted

**Avatar: The Last Airbender**
- id: avatar-last-airbender
- Shapes: Precarious balance; Slow drift
- Scenes: keeping the four nations in balance — one move and the whole thing tips; the world drifting toward war over a hundred years nobody stopped — slow drift no one arrested in time
- GIF search: avatar last airbender aang giphy
- Travels externally: yes — strong cross-generational streaming recognition

## Film

**Star Wars**
- id: star-wars
- Shapes: Fragility; Early warning ignored; Declaring a fix that isn't one
- Scenes: the Death Star's unshielded exhaust port — one small weakness nobody reviewed taking down the whole system; Admiral Ackbar's "It's a trap!" — recognizing a known failure mode a beat too late; "these aren't the droids you're looking for" — waved through a checkpoint without the scrutiny it deserved
- GIF search: star wars its a trap ackbar giphy
- Travels externally: yes — near-universal cultural recognition

**Star Trek**
- id: star-trek
- Shapes: Coordination difficulty; Diminishing returns / breaking point
- Scenes: Picard's "Make it so" — a decision handed off to be executed without relitigating; Scotty padding a repair estimate to look like a miracle worker — headroom baked into the number
- GIF search: star trek make it so giphy
- Travels externally: yes — marquee lines are near-universal

**Jurassic Park**
- id: jurassic-park
- Shapes: Early warning ignored; Denial of consequence; Cascading failure
- Scenes: "life finds a way" — the warning waved off right before it proves true; one power failure cascading into every enclosure at once — a chain reaction from a single point
- GIF search: jurassic park life finds a way giphy
- Travels externally: yes — near-universal recognition

**The Lord of the Rings**
- id: lord-of-the-rings
- Shapes: Persistent thankless effort; Precarious balance
- Scenes: Frodo and Sam grinding step after step toward Mordor — a thankless march that has to be finished; "one does not simply walk into Mordor" — a plan balanced on a single impossible-looking step
- GIF search: one does not simply walk into mordor giphy
- Travels externally: yes — near-universal, heavily meme'd

**Inception**
- id: inception
- Shapes: Accidental complexity; Fragmented understanding
- Scenes: dreams nested inside dreams inside dreams — layers of complexity no one can hold at once; each character trusting a different version of what's real — fragmented understanding of the same system
- GIF search: inception we need to go deeper giphy
- Travels externally: yes — broad recognition, "we need to go deeper" is a meme

**The Matrix**
- id: the-matrix
- Shapes: Obvious problem nobody names; Unseen scale
- Scenes: the red-pill choice — naming the thing everyone lived inside and never questioned; the reveal that the world is vastly larger and stranger than anyone saw — the real scale hidden under the surface
- GIF search: the matrix red pill blue pill giphy
- Travels externally: yes — near-universal cultural recognition

**Titanic**
- id: titanic
- Shapes: Early warning ignored; Denial of consequence; Precarious balance
- Scenes: iceberg warnings brushed aside at full speed — a plain signal ignored until impact; "unsinkable" right up to the moment it wasn't — denial balanced on false confidence
- GIF search: titanic iceberg giphy
- Travels externally: yes — near-universal recognition

**Jaws**
- id: jaws
- Shapes: Early warning ignored; Denial of consequence
- Scenes: the mayor keeping the beaches open despite the warnings — a known risk denied for convenience; "you're gonna need a bigger boat" — realizing the problem is far larger than the plan
- GIF search: jaws bigger boat giphy
- Travels externally: yes — the "bigger boat" line is near-universal

**Ocean's Eleven**
- id: oceans-eleven
- Shapes: Coordination difficulty; Over-engineering
- Scenes: eleven specialists timed to the second — coordination that only works if everyone hits their mark; a heist plan with ten moving parts where two would nearly do — elaborate machinery as the point
- GIF search: oceans eleven giphy
- Travels externally: yes — broad recognition

**Back to the Future**
- id: back-to-the-future
- Shapes: Cascading failure; Shifting requirements
- Scenes: one changed moment in the past rippling through the whole timeline — a small edit cascading downstream; the photo fading and re-forming as the plan's target keeps moving
- GIF search: back to the future giphy
- Travels externally: yes — near-universal recognition

**Mission: Impossible**
- id: mission-impossible
- Shapes: Precarious balance; Coordination breakdown under pressure
- Scenes: Ethan suspended an inch above the pressure-sensor floor — the whole job balanced on not touching anything; the drilled team improvising live when one step of the plan goes wrong
- GIF search: mission impossible vault giphy
- Travels externally: yes — broad global recognition

## Sports

**Formula 1**
- id: formula-1
- Shapes: Momentum without control; Coordination breakdown under pressure
- Scenes: a car losing control and crashing at speed — moving fast with no margin until there's none left; pit-stop chaos — a slick, drilled process falling apart the moment one step goes wrong
- GIF search: formula 1 pit stop fail giphy
- Travels externally: limited — a crash reads universally; pit-stop nuance needs some F1 familiarity

**Olympic relay**
- id: olympic-relay
- Shapes: Coordination breakdown under pressure; Cascading failure
- Scenes: the dropped baton at the handoff — a whole race lost at the one coordination point; one fumble at the exchange cascading down the order
- GIF search: olympic relay dropped baton giphy
- Travels externally: yes — the dropped baton reads universally

**Marathon running**
- id: marathon
- Shapes: Persistent thankless effort; Slow drift
- Scenes: grinding out mile after mile with the finish nowhere in sight — thankless, sustained effort; "hitting the wall" late in the race — a slow drain nobody noticed until it stopped them
- GIF search: marathon runner hitting the wall giphy
- Travels externally: yes — broadly understood

## Games

**Jenga**
- id: jenga
- Shapes: Precarious balance; Fragility
- Scenes: the tower leaning one block from collapse — one more addition and it tips; the whole stack coming down from pulling a single wrong piece — fragility hiding in plain sight
- GIF search: jenga tower collapse giphy
- Travels externally: yes — near-universal

**Tetris**
- id: tetris
- Shapes: Chaotic parallel effort; Cascading failure
- Scenes: pieces falling faster than they can be placed — work arriving faster than it clears; the stack reaching the top and ending the game — a backlog cascading into failure
- GIF search: tetris game over giphy
- Travels externally: yes — near-universal

**Whack-a-Mole**
- id: whack-a-mole
- Shapes: Recurring bugs; Persistent thankless effort
- Scenes: knocking one mole down as two more pop up — fixes that spawn the next problem; hammering endlessly with the moles never stopping — thankless, unwinnable repetition
- GIF search: whack a mole arcade giphy
- Travels externally: yes — near-universal metaphor

**Mario Kart**
- id: mario-kart
- Shapes: Momentum without control; Shifting requirements
- Scenes: leading the whole race until the blue shell hits — momentum wiped out at the worst moment; the track and rules shifting under you every lap — the target never holding still
- GIF search: mario kart blue shell giphy
- Travels externally: yes — broad recognition

**The Sims**
- id: the-sims
- Shapes: Over-engineering; Excess / feature bloat
- Scenes: a dream house with every gadget and no door to the bathroom — features bolted on past the point of use; an elaborate build that forgets the one thing the Sim actually needed
- GIF search: the sims giphy
- Travels externally: yes — broad recognition

**Chess**
- id: chess
- Shapes: Coordination difficulty; Precarious balance
- Scenes: every piece dependent on the others to hold the position — coordination where one bad move unravels the rest; a position balanced so finely a single tempo loses it
- GIF search: chess checkmate giphy
- Travels externally: yes — universal

**Minecraft**
- id: minecraft
- Shapes: Reinventing effort; Chaotic parallel effort
- Scenes: tearing down a working build to remake it "properly" — rebuilding from scratch instead of reusing; a dozen half-finished projects sprawling at once — everyone building everything in parallel
- GIF search: minecraft giphy
- Travels externally: yes — near-universal among a broad audience

**Pac-Man**
- id: pac-man
- Shapes: Recurring bugs; Momentum without control
- Scenes: the ghosts respawning the instant they're cleared — a problem that keeps coming back; racing the maze until one wrong turn ends the run — speed with no room to recover
- GIF search: pac man giphy
- Travels externally: yes — universal
