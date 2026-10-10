#!/usr/bin/env node
// Snapshot slides from a frontend-slides deck (a fixed 1920x1080 stage where the
// `.slide.active` slide is visible) by driving headless Chrome over the DevTools
// protocol. Doc tooling only; needs Google Chrome and Node 22+ (built-in WebSocket).
//
//   node docs/capture-slides.ts <deck.html> --meta
//   node docs/capture-slides.ts <deck.html> --out <dir> --name <prefix> --slides 0,7,12 [--width 1280] [--quality 80]
//
// --meta prints one JSON line per slide (index, title, svg/img counts, a text
// snippet, and whether the slide mentions a term from --avoid) so slides can be
// chosen without opening the deck.
import { spawn } from "node:child_process";
import type { ChildProcess } from "node:child_process";
import { mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";

const CHROME = process.env.CHROME ?? "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome";
const STAGE_WIDTH = 1920;
const STAGE_HEIGHT = 1080;
const DEFAULT_WIDTH = 1280;
const DEFAULT_QUALITY = 80;
const SNIPPET_CHARS = 160;
const POLL_MS = 100;
const POLL_ATTEMPTS = 100;
const SETTLE_MS = 1100;
const RM_RETRIES = 5;
const INDEX_DIGITS = 2; // slide numbers in file names: 00, 07, 12
const ARGV_START = 2; // argv[0] is node, argv[1] is this script

type Json = Record<string, unknown>;

interface Cdp {
  readonly send: (method: string, params?: Json) => Promise<Json>;
  readonly close: () => void;
}

interface SlideInfo {
  readonly index: number;
  readonly title: string;
  readonly svg: number;
  readonly img: number;
  readonly text: string;
}

const isRecord = (v: unknown): v is Json => typeof v === "object" && v !== null && !Array.isArray(v);
const sleep = (ms: number): Promise<void> => new Promise((done) => setTimeout(done, ms));

function flag(args: readonly string[], name: string): string | undefined {
  const at = args.indexOf(`--${name}`);
  return at === -1 ? undefined : args[at + 1];
}

// ---- Chrome ----------------------------------------------------------------

async function waitFor<T>(probe: () => T | undefined, what: string): Promise<T> {
  const attempt = async (left: number): Promise<T> => {
    const value = probe();
    if (value !== undefined) return value;
    if (left === 0) throw new Error(`timed out waiting for ${what}`);
    await sleep(POLL_MS);
    return attempt(left - 1);
  };
  return attempt(POLL_ATTEMPTS);
}

function readDebugPort(profile: string): number | undefined {
  try {
    const first = readFileSync(join(profile, "DevToolsActivePort"), "utf8").split("\n")[0] ?? "";
    return first === "" ? undefined : Number(first);
  } catch {
    return undefined;
  }
}

function launch(profile: string): ChildProcess {
  return spawn(
    CHROME,
    [
      "--headless=new",
      "--disable-gpu",
      "--hide-scrollbars",
      "--no-first-run",
      "--remote-debugging-port=0",
      `--user-data-dir=${profile}`,
      `--window-size=${STAGE_WIDTH},${STAGE_HEIGHT}`,
      "about:blank",
    ],
    { stdio: "ignore" },
  );
}

async function pageSocketUrl(port: number): Promise<string> {
  const res = await fetch(`http://127.0.0.1:${port}/json/list`);
  const targets: unknown = await res.json();
  const page = Array.isArray(targets) ? targets.find((t) => isRecord(t) && t["type"] === "page") : undefined;
  const url = isRecord(page) ? page["webSocketDebuggerUrl"] : undefined;
  if (typeof url !== "string") throw new Error("no debuggable page found");
  return url;
}

async function connect(url: string): Promise<Cdp> {
  const socket = new WebSocket(url);
  const pending = new Map<number, { resolve: (v: Json) => void; reject: (e: Error) => void }>();
  const ids = { next: 0 };
  socket.addEventListener("message", (event) => {
    const msg: unknown = JSON.parse(String(event.data));
    if (!isRecord(msg) || typeof msg["id"] !== "number") return;
    const waiter = pending.get(msg["id"]);
    pending.delete(msg["id"]);
    if (isRecord(msg["error"])) waiter?.reject(new Error(String(msg["error"]["message"])));
    else waiter?.resolve(isRecord(msg["result"]) ? msg["result"] : {});
  });
  await new Promise<void>((ready) => socket.addEventListener("open", () => ready(), { once: true }));
  return {
    send: (method, params = {}) =>
      new Promise<Json>((resolveCall, reject) => {
        ids.next += 1;
        pending.set(ids.next, { resolve: resolveCall, reject });
        socket.send(JSON.stringify({ id: ids.next, method, params }));
      }),
    close: () => socket.close(),
  };
}

async function evaluate(cdp: Cdp, expression: string): Promise<unknown> {
  const res = await cdp.send("Runtime.evaluate", { expression, returnByValue: true, awaitPromise: true });
  const result = res["result"];
  return isRecord(result) ? result["value"] : undefined;
}

// ---- slides ----------------------------------------------------------------

const LIST_SLIDES = `(() => [...document.querySelectorAll('.slide')].map((el, index) => ({
  index,
  title: (el.querySelector('h1,h2,h3')?.textContent || '').replace(/\\s+/g, ' ').trim(),
  svg: el.querySelectorAll('svg').length,
  img: el.querySelectorAll('img').length,
  text: el.textContent.replace(/\\s+/g, ' ').trim(),
})))()`;

const showSlide = (index: number): string => `((i) => {
  const slides = [...document.querySelectorAll('.slide')];
  slides.forEach((s) => s.classList.remove('active', 'visible'));
  slides[i].classList.add('active', 'visible');
  return true;
})(${index})`;

function toSlideInfo(raw: unknown): SlideInfo | undefined {
  if (!isRecord(raw)) return undefined;
  const { index, title, svg, img, text } = raw;
  const wellFormed =
    typeof index === "number" &&
    typeof title === "string" &&
    typeof svg === "number" &&
    typeof img === "number" &&
    typeof text === "string";
  return wellFormed ? { index, title, svg, img, text } : undefined;
}

async function listSlides(cdp: Cdp): Promise<readonly SlideInfo[]> {
  const raw = await evaluate(cdp, LIST_SLIDES);
  const infos = Array.isArray(raw) ? raw.map(toSlideInfo) : [];
  return infos.filter((s): s is SlideInfo => s !== undefined);
}

async function snapshot(cdp: Cdp, index: number, width: number, quality: number): Promise<Buffer> {
  await evaluate(cdp, showSlide(index));
  await sleep(SETTLE_MS);
  const clip = { x: 0, y: 0, width: STAGE_WIDTH, height: STAGE_HEIGHT, scale: width / STAGE_WIDTH };
  const shot = await cdp.send("Page.captureScreenshot", { format: "jpeg", quality, clip });
  return Buffer.from(String(shot["data"]), "base64");
}

// ---- main ------------------------------------------------------------------

async function withDeck<T>(deck: string, run: (cdp: Cdp) => Promise<T>): Promise<T> {
  const profile = mkdtempSync(join(tmpdir(), "capture-slides-"));
  const chrome = launch(profile);
  try {
    const port = await waitFor(() => readDebugPort(profile), "Chrome's debug port");
    const cdp = await connect(await pageSocketUrl(port));
    await cdp.send("Page.enable");
    await cdp.send("Emulation.setDeviceMetricsOverride", {
      width: STAGE_WIDTH,
      height: STAGE_HEIGHT,
      deviceScaleFactor: 1,
      mobile: false,
    });
    await cdp.send("Page.navigate", { url: pathToFileURL(resolve(deck)).href });
    await evaluate(cdp, "new Promise((r) => (document.readyState === 'complete' ? r(true) : addEventListener('load', () => r(true))))");
    await evaluate(cdp, "document.fonts.ready.then(() => true)");
    const result = await run(cdp);
    cdp.close();
    return result;
  } finally {
    const exited = new Promise<void>((done) => chrome.once("exit", () => done()));
    chrome.kill();
    await exited;
    rmSync(profile, { recursive: true, force: true, maxRetries: RM_RETRIES, retryDelay: POLL_MS });
  }
}

async function main(): Promise<void> {
  const args = process.argv.slice(ARGV_START);
  const deck = args[0];
  if (deck === undefined) throw new Error("usage: capture-slides.ts <deck.html> (--meta | --out <dir> --slides 0,3)");
  const avoid = (flag(args, "avoid") ?? "").split(",").filter((t) => t !== "");

  if (args.includes("--meta")) {
    const slides = await withDeck(deck, listSlides);
    slides.forEach((s) => {
      const lower = s.text.toLowerCase();
      const hits = avoid.filter((term) => lower.includes(term.toLowerCase()));
      console.log(JSON.stringify({ ...s, text: s.text.slice(0, SNIPPET_CHARS), avoid: hits }));
    });
    return;
  }

  const out = flag(args, "out");
  const name = flag(args, "name");
  const picks = (flag(args, "slides") ?? "").split(",").filter((n) => n !== "").map(Number);
  if (out === undefined || name === undefined || picks.length === 0) throw new Error("--out, --name and --slides are required");
  const width = Number(flag(args, "width") ?? DEFAULT_WIDTH);
  const quality = Number(flag(args, "quality") ?? DEFAULT_QUALITY);
  mkdirSync(out, { recursive: true });

  await withDeck(deck, (cdp) =>
    picks.reduce<Promise<void>>(async (previous, index) => {
      await previous;
      const file = join(out, `${name}-${String(index).padStart(INDEX_DIGITS, "0")}.jpg`);
      writeFileSync(file, await snapshot(cdp, index, width, quality));
      console.log(`wrote ${file}`);
    }, Promise.resolve()),
  );
}

main().catch((err: unknown) => {
  console.error(err instanceof Error ? err.message : String(err));
  process.exit(1);
});
