#!/usr/bin/env node
/**
 * validate_palette.js — categorical / ordinal palette validator for deck-visuals.
 *
 * PROVENANCE: this is an independent MIT reimplementation of the checks
 * performed by Claude Code's built-in `dataviz` skill validator. That
 * built-in skill is the authority — when it is reachable, run it first and
 * treat its result as the real one. This file exists for the cases where
 * it isn't: a standalone shell, CI, or a PR review outside a Claude Code
 * session. It is kept deliberately in agreement with the built-in
 * validator (same thresholds, same checks, same CLI shape) rather than
 * copied from it — the color-science constants below (the Machado,
 * Oliveira & Fernandes 2009 CVD matrices, Björn Ottosson's OKLab
 * coefficients, the WCAG relative-luminance/contrast formula, and the
 * sRGB transfer function) are all public, standard values; the assembly
 * into these five checks is written from scratch here.
 *
 * Checks, all computed in OKLCH / OKLab:
 *   1. Lightness band     — each color's OKLCH L sits inside the mode's
 *                           band: light [0.43, 0.77], dark [0.48, 0.67].
 *   2. Chroma floor        — each color's OKLCH C >= 0.10; below that a
 *                           hue reads as gray and stops doing identity work.
 *   3. CVD separation      — worst OKLab ΔE×100 (Euclidean distance in
 *                           OKLab, scaled by 100) across the active
 *                           pairlist (adjacent slots by default; every
 *                           combination with --pairs all), taking
 *                           min(protanopia, deuteranopia) per pair under
 *                           the Machado-Oliveira-Fernandes (2009)
 *                           severity-1.0 simulation. >=8.0 is a clean pass,
 *                           6.0-8.0 is a WARN legal only alongside a
 *                           secondary encoding (labels, gaps, texture),
 *                           below 6.0 is a hard FAIL. Tritanopia is
 *                           reported alongside for visibility but is not
 *                           gated.
 *   4. Normal-vision floor — the same ΔE×100 formula, unsimulated, on the
 *                           same pairlist; the worst pair must clear 15.0.
 *                           This is a hard gate — no WARN band, and a
 *                           secondary encoding does not excuse it.
 *   5. Contrast vs surface — WCAG contrast ratio of each color against the
 *                           given surface hex. >=3.0 is a pass; below that
 *                           is a WARN that obligates a relief channel
 *                           (visible direct labels or a table view), not a
 *                           hard FAIL.
 *
 * --ordinal switches to the ramp checks for a single-hue sequential or
 * diverging-arm ramp instead of the five checks above:
 *   - Lightness monotone   — OKLCH L must strictly sort light->dark or
 *                           dark->light in input order.
 *   - Adjacent ΔL           — each step must move OKLCH L by >= 0.06 from
 *                           its neighbor.
 *   - Light-end contrast    — the lightest step must clear 2.0:1 against
 *                           the surface.
 *   - Single hue             — OKLab hue-angle spread across all steps
 *                           must stay <= 40°, or it isn't a one-hue ramp.
 *
 * Usage:
 *   node validate_palette.js "#hex,#hex,..." --mode dark --surface "#1a1a19"
 *   node validate_palette.js "#hex,#hex,..." --mode light
 *   node validate_palette.js "#hex,#hex,..." --mode dark --pairs all
 *   node validate_palette.js "#hex,#hex,..." --mode dark --ordinal
 *
 * Flags:
 *   --mode light|dark      selects the lightness band and the default
 *                          surface. Defaults to "light".
 *   --surface HEX          surface color to contrast-check against.
 *                          Defaults to #fcfcfb (light) / #1a1a19 (dark).
 *   --pairs adjacent|all   adjacent (default) checks only consecutive
 *                          slots — how stacks/bars/lines are read. all
 *                          checks every combination — required for
 *                          scatter/bubble/map/small-multiples, where any
 *                          two marks can end up side by side.
 *   --ordinal              validate a one-hue sequential/diverging-arm
 *                          ramp instead of a categorical set.
 *
 * Exit code: 1 on any hard FAIL. WARN bands (CVD 6-8, sub-3:1 contrast)
 * still exit 0 — each requires a secondary/relief channel, but does not
 * fail the run by itself.
 */

'use strict';

// -- thresholds --------------------------------------------------------------
const LIGHTNESS_BAND = { light: [0.43, 0.77], dark: [0.48, 0.67] }; // OKLCH L
const CHROMA_FLOOR = 0.10; // OKLCH C
const CVD_TARGET = 8.0;    // OKLab ΔE×100, min(protan, deutan) — clean pass
const CVD_FLOOR = 6.0;     // OKLab ΔE×100 — below this is a hard FAIL
const NORMAL_FLOOR = 15.0; // OKLab ΔE×100, unsimulated, worst pair — hard gate
const CONTRAST_MIN = 3.0;  // WCAG ratio vs surface
const DEFAULT_SURFACE = { light: '#fcfcfb', dark: '#1a1a19' };
const ORDINAL_MIN_DELTA_L = 0.06; // OKLCH L step between ordinal ramp entries
const ORDINAL_LIGHT_END_FLOOR = 2.0; // WCAG ratio for the ramp's lightest step
const ORDINAL_MAX_HUE_SPREAD = 40; // degrees — beyond this it's not one hue

// Machado, Oliveira & Fernandes (2009), full (severity 1.0) dichromacy
// simulation matrices, applied in linear sRGB.
const CVD_MATRIX = {
  protan: [
    [0.152286, 1.052583, -0.204868],
    [0.114503, 0.786281, 0.099216],
    [-0.003882, -0.048116, 1.051998],
  ],
  deutan: [
    [0.367322, 0.860646, -0.227968],
    [0.280085, 0.672501, 0.047413],
    [-0.011820, 0.042940, 0.968881],
  ],
  tritan: [
    [1.255528, -0.076749, -0.178779],
    [-0.078411, 0.930809, 0.147602],
    [0.004733, 0.691367, 0.303900],
  ],
};

// -- input normalization ------------------------------------------------------
const HEX_RE = /^#?[0-9a-fA-F]{6}$/;
const isHex = (v) => HEX_RE.test(v);
const splitColors = (raw) => String(raw || '').split(',').map((c) => c.trim()).filter(Boolean);

// -- color math ----------------------------------------------------------------
function hexToSrgb(hex) {
  const h = hex.trim().replace(/^#/, '');
  return [0, 2, 4].map((i) => parseInt(h.slice(i, i + 2), 16) / 255);
}

function srgbChannelToLinear(c) {
  return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
}

function hexToLinear(hex) {
  return hexToSrgb(hex).map(srgbChannelToLinear);
}

function relativeLuminance(hex) {
  const [r, g, b] = hexToLinear(hex);
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrastRatio(hexA, hexB) {
  const [hi, lo] = [relativeLuminance(hexA), relativeLuminance(hexB)].sort((a, b) => b - a);
  return (hi + 0.05) / (lo + 0.05);
}

function linearToOklab([r, g, b]) {
  const l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b;
  const m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b;
  const s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b;
  const l3 = Math.cbrt(l), m3 = Math.cbrt(m), s3 = Math.cbrt(s);
  return [
    0.2104542553 * l3 + 0.7936177850 * m3 - 0.0040720468 * s3, // L
    1.9779984951 * l3 - 2.4285922050 * m3 + 0.4505937099 * s3, // a
    0.0259040371 * l3 + 0.7827717662 * m3 - 0.8086757660 * s3, // b
  ];
}

const oklabOf = (hex) => linearToOklab(hexToLinear(hex));
const oklchOf = (hex) => { const [L, a, b] = oklabOf(hex); return [L, Math.hypot(a, b)]; };
const hueAngleOf = (hex) => {
  const [, a, b] = oklabOf(hex);
  return ((Math.atan2(b, a) * 180 / Math.PI) % 360 + 360) % 360;
};

function simulateCvd(hex, kind) {
  const [r, g, b] = hexToLinear(hex);
  const M = CVD_MATRIX[kind];
  const clamp = (c) => Math.max(0, Math.min(1, c));
  return M.map(([m0, m1, m2]) => clamp(m0 * r + m1 * g + m2 * b));
}

function deltaEOklab(hexA, hexB, kind) {
  const a = kind ? linearToOklab(simulateCvd(hexA, kind)) : oklabOf(hexA);
  const b = kind ? linearToOklab(simulateCvd(hexB, kind)) : oklabOf(hexB);
  return 100 * Math.hypot(a[0] - b[0], a[1] - b[1], a[2] - b[2]);
}

// -- pairlist ------------------------------------------------------------------
function buildPairlist(n, mode) {
  if (mode === 'all') {
    const pairs = [];
    for (let i = 0; i < n; i++) for (let j = i + 1; j < n; j++) pairs.push([i, j]);
    return pairs;
  }
  const pairs = [];
  for (let i = 0; i < n - 1; i++) pairs.push([i, i + 1]);
  return pairs;
}

// -- categorical checks ----------------------------------------------------------
function validateCategorical(palette, { mode, surface, pairs }) {
  const [lo, hi] = LIGHTNESS_BAND[mode];
  const report = [];
  let hardFail = false;

  const offBand = palette.filter((c) => { const [L] = oklchOf(c); return L < lo || L > hi; });
  if (offBand.length) hardFail = true;
  report.push({
    check: 'Lightness band', state: offBand.length ? 'fail' : 'pass',
    detail: offBand.length
      ? `outside [${lo}, ${hi}]: ${offBand.map((c) => `${c} (L=${oklchOf(c)[0].toFixed(3)})`).join(', ')}`
      : `all ${palette.length} inside L ${lo}–${hi}`,
  });

  const belowChroma = palette.filter((c) => oklchOf(c)[1] < CHROMA_FLOOR);
  if (belowChroma.length) hardFail = true;
  report.push({
    check: 'Chroma floor', state: belowChroma.length ? 'fail' : 'pass',
    detail: belowChroma.length
      ? `below ${CHROMA_FLOOR}: ${belowChroma.map((c) => `${c} (C=${oklchOf(c)[1].toFixed(3)})`).join(', ')}`
      : `all ${palette.length} >= ${CHROMA_FLOOR}`,
  });

  const pairlist = buildPairlist(palette.length, pairs);
  const pairLabel = pairs === 'all' ? 'all-pairs' : 'adjacent';

  let worstCvd = null;
  for (const kind of ['protan', 'deutan']) {
    for (const [i, j] of pairlist) {
      const d = deltaEOklab(palette[i], palette[j], kind);
      if (!worstCvd || d < worstCvd.d) worstCvd = { d, kind, a: palette[i], b: palette[j] };
    }
  }
  const worstTritan = pairlist.length
    ? Math.min(...pairlist.map(([i, j]) => deltaEOklab(palette[i], palette[j], 'tritan')))
    : null;
  const cvdD = worstCvd ? worstCvd.d : Infinity;
  const cvdState = cvdD >= CVD_TARGET ? 'pass' : cvdD >= CVD_FLOOR ? 'warn' : 'fail';
  if (cvdState === 'fail') hardFail = true;
  report.push({
    check: 'CVD separation', state: cvdState,
    detail: worstCvd
      ? `worst ${pairLabel} ${worstCvd.a}↔${worstCvd.b} ΔE ${cvdD.toFixed(1)} (${worstCvd.kind})`
        + (worstTritan != null ? ` · tritan ${worstTritan.toFixed(1)}` : '')
      : 'n/a (single color)',
  });

  let worstNormal = null;
  for (const [i, j] of pairlist) {
    const d = deltaEOklab(palette[i], palette[j]);
    if (!worstNormal || d < worstNormal.d) worstNormal = { d, a: palette[i], b: palette[j] };
  }
  const normD = worstNormal ? worstNormal.d : Infinity;
  const normState = normD >= NORMAL_FLOOR ? 'pass' : 'fail';
  if (normState === 'fail') hardFail = true;
  report.push({
    check: 'Normal-vision floor', state: normState,
    detail: worstNormal
      ? `worst ${pairLabel} ${worstNormal.a}↔${worstNormal.b} ΔE ${normD.toFixed(1)}`
        + (normState === 'fail' ? ` — below ${NORMAL_FLOOR} floor` : '')
      : 'n/a (single color)',
  });

  const belowContrast = palette.filter((c) => contrastRatio(c, surface) < CONTRAST_MIN);
  report.push({
    check: 'Contrast vs surface', state: belowContrast.length ? 'warn' : 'pass',
    detail: belowContrast.length
      ? `below ${CONTRAST_MIN}:1, relief required (labels or table view): `
        + belowContrast.map((c) => `${c} (${contrastRatio(c, surface).toFixed(2)}:1)`).join(', ')
      : `all ${palette.length} >= ${CONTRAST_MIN}:1`,
  });

  return { report, ok: !hardFail };
}

// -- ordinal ramp check ----------------------------------------------------------
function validateOrdinal(palette, { mode, surface }) {
  const report = [];
  let hardFail = false;
  const Ls = palette.map((c) => oklchOf(c)[0]);

  const ascending = Ls.every((v, i) => i === 0 || v >= Ls[i - 1]);
  const descending = Ls.every((v, i) => i === 0 || v <= Ls[i - 1]);
  const monotone = ascending || descending;
  if (!monotone) hardFail = true;
  report.push({
    check: 'Lightness monotone', state: monotone ? 'pass' : 'fail',
    detail: monotone ? 'reads light→dark (or reverse)'
      : `not monotone: L = ${Ls.map((l) => l.toFixed(3)).join(', ')}`,
  });

  const gaps = Ls.slice(1).map((l, i) => Math.abs(l - Ls[i]));
  const thin = gaps.filter((g) => g < ORDINAL_MIN_DELTA_L);
  if (thin.length) hardFail = true;
  report.push({
    check: 'Adjacent ΔL', state: thin.length ? 'fail' : 'pass',
    detail: thin.length
      ? `${thin.length} step(s) below ${ORDINAL_MIN_DELTA_L}: gaps = ${gaps.map((g) => g.toFixed(3)).join(', ')}`
      : `all steps >= ${ORDINAL_MIN_DELTA_L} (gaps = ${gaps.map((g) => g.toFixed(3)).join(', ')})`,
  });

  const sortedByL = [...palette].sort((a, b) => oklchOf(a)[0] - oklchOf(b)[0]);
  const lightestStep = mode === 'light' ? sortedByL[sortedByL.length - 1] : sortedByL[0];
  const lightContrast = contrastRatio(lightestStep, surface);
  const lightOk = lightContrast >= ORDINAL_LIGHT_END_FLOOR;
  if (!lightOk) hardFail = true;
  report.push({
    check: 'Light-end contrast', state: lightOk ? 'pass' : 'fail',
    detail: `${lightestStep} at ${lightContrast.toFixed(2)}:1 vs surface`
      + (lightOk ? '' : ` — below ${ORDINAL_LIGHT_END_FLOOR}:1 floor`),
  });

  const hues = palette.map(hueAngleOf);
  let spread = hues.length ? Math.max(...hues) - Math.min(...hues) : 0;
  if (spread > 180) spread = 360 - spread;
  const oneHue = spread <= ORDINAL_MAX_HUE_SPREAD;
  if (!oneHue) hardFail = true;
  report.push({
    check: 'Single hue', state: oneHue ? 'pass' : 'fail',
    detail: `hue spread ${spread.toFixed(0)}°`
      + (oneHue ? '' : ` — exceeds ${ORDINAL_MAX_HUE_SPREAD}°, not a one-hue ramp`),
  });

  return { report, ok: !hardFail };
}

// -- CLI -------------------------------------------------------------------------
const STATE_GLYPH = { pass: 'PASS', warn: 'WARN', fail: 'FAIL' };

function printReport({ report, ok }, { mode, surface, ordinal, count }) {
  const kind = ordinal ? 'ordinal ramp' : 'categorical';
  console.log(`\nPalette (${mode}, surface ${surface}, ${kind}): ${count} slot(s)`);
  for (const { check, state, detail } of report) {
    console.log(`  [${STATE_GLYPH[state].padEnd(4)}] ${check.padEnd(22)} ${detail}`);
  }
  console.log(`\n  → ${ok ? 'ALL CHECKS PASS' : 'FAILED — fix the marked checks'}`);
  if (!ordinal) {
    console.log('  WARN on CVD (6–8 floor band) is legal only with a secondary encoding.');
    console.log('  WARN on contrast (sub-3:1) obligates a relief channel (labels or table view).');
  }
}

function parseArgs(argv) {
  const VALUE_FLAGS = new Set(['--mode', '--surface', '--pairs']);
  const opts = {};
  let positional = null;
  for (let i = 0; i < argv.length; i++) {
    let a = argv[i], val;
    const eq = a.indexOf('=');
    if (eq > 0) { val = a.slice(eq + 1); a = a.slice(0, eq); }
    if (VALUE_FLAGS.has(a)) opts[a.slice(2)] = val ?? argv[++i];
    else if (a === '--ordinal') opts.ordinal = true;
    else if (a.startsWith('--')) { console.error(`unknown flag: ${a}`); process.exit(2); }
    else if (positional === null) positional = a;
    else { console.error(`unexpected extra positional: ${a}`); process.exit(2); }
  }
  return { opts, positional };
}

function runCli(argv) {
  const { opts, positional } = parseArgs(argv);
  const palette = splitColors(positional);
  if (!palette.length) {
    console.error('usage: node validate_palette.js "#hex,#hex,..." [--mode light|dark] [--surface #hex] [--pairs adjacent|all] [--ordinal]');
    process.exit(2);
  }
  const mode = opts.mode || 'light';
  if (!['light', 'dark'].includes(mode)) {
    console.error(`--mode must be "light" or "dark" (got ${JSON.stringify(mode)})`);
    process.exit(2);
  }
  const pairs = opts.pairs || 'adjacent';
  if (!['adjacent', 'all'].includes(pairs)) {
    console.error(`--pairs must be "adjacent" or "all" (got ${JSON.stringify(pairs)})`);
    process.exit(2);
  }
  const surface = (opts.surface && opts.surface.trim()) || DEFAULT_SURFACE[mode];
  const badHex = [...palette, surface].filter((c) => !isHex(c));
  if (badHex.length) {
    console.error(`invalid hex value(s): ${badHex.join(', ')} — expected #rrggbb`);
    process.exit(2);
  }

  const result = opts.ordinal
    ? validateOrdinal(palette, { mode, surface })
    : validateCategorical(palette, { mode, surface, pairs });
  printReport(result, { mode, surface, ordinal: !!opts.ordinal, count: palette.length });
  process.exit(result.ok ? 0 : 1);
}

if (require.main === module) {
  runCli(process.argv.slice(2));
}

module.exports = {
  validateCategorical, validateOrdinal, contrastRatio, oklchOf, deltaEOklab,
};
