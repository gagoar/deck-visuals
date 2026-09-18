#!/usr/bin/env python3
"""
coverage_report.py — does the presenter's recognized pool give the levity matcher
real choice?

The "3 options per levity offering" target means: for each concept shape in the
seed-library table of references/levity.md, at least three recognized sources in
assets/show-catalog.md carry that shape. This script measures it.

It parses the shape vocabulary from levity.md (the single source of truth), the
`Shapes:` tags from show-catalog.md, and — when given a results.json from the
onboarding picker — counts only the shows the presenter recognized. Without
--results it reports the catalog's full capacity (every title assumed known), which
is useful when curating the catalog itself.

Also flags drift: any catalog shape tag that does not match a levity shape, so the
two files can't silently fall out of sync.

Usage:
  python3 coverage_report.py --levity <levity.md> --catalog <show-catalog.md> \
      [--results <results.json>] [--target 3]

Exit code: 1 if any shape is under target (a GAP) or drift is found, else 0.
Standard library only.
"""

import argparse
import json
import re
import sys


def normalize(shape):
    """Lenient key: lowercase, drop a trailing '(alt)', punctuation → spaces."""
    s = shape.lower()
    s = re.sub(r"\(alt\)", " ", s)
    s = re.sub(r"[^a-z0-9]+", " ", s)
    return " ".join(s.split())


def parse_levity_shapes(path):
    """First column of every seed-library table row → {norm: display}.

    Scoped to the "## Seed library" section only, so the reaction-mode beat table
    (Facepalm, Relief, …) elsewhere in levity.md is not mistaken for a concept shape.
    """
    shapes = {}
    in_section = False
    with open(path, encoding="utf-8") as fh:
        for line in fh:
            line = line.rstrip("\n")
            if line.startswith("## "):
                in_section = "seed library" in line.lower()
                continue
            if not in_section:
                continue
            if not line.startswith("|"):
                continue
            cells = [c.strip() for c in line.strip().strip("|").split("|")]
            if not cells:
                continue
            first = cells[0]
            # Skip the header row and the |---|---| separator.
            if not first or set(first) <= set("-: "):
                continue
            if first.lower() == "concept shape":
                continue
            key = normalize(first)
            if not key:
                continue
            display = re.sub(r"\s*\(alt\)\s*", "", first).strip()
            shapes.setdefault(key, display)
    return shapes


def parse_catalog(path):
    """Return (entries, order): entries[id] = {'title','shapes':[norm,...]}."""
    entries = {}
    order = []
    cur_id = None
    with open(path, encoding="utf-8") as fh:
        title = None
        for line in fh:
            stripped = line.strip()
            m_title = re.match(r"^\*\*(.+?)\*\*$", stripped)
            if m_title:
                title = m_title.group(1).strip()
                cur_id = None
                continue
            m_id = re.match(r"^-\s*id:\s*(.+)$", stripped)
            if m_id:
                cur_id = m_id.group(1).strip()
                entries[cur_id] = {"title": title or cur_id, "shapes": []}
                order.append(cur_id)
                continue
            m_shapes = re.match(r"^-\s*Shapes:\s*(.+)$", stripped)
            if m_shapes and cur_id:
                tags = [t.strip() for t in m_shapes.group(1).split(";") if t.strip()]
                entries[cur_id]["shapes"] = [normalize(t) for t in tags]
                # keep raw for drift reporting
                entries[cur_id]["raw_shapes"] = tags
    return entries, order


def main(argv):
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--levity", required=True)
    ap.add_argument("--catalog", required=True)
    ap.add_argument("--results", default=None)
    ap.add_argument("--target", type=int, default=3)
    args = ap.parse_args(argv[1:])

    levity = parse_levity_shapes(args.levity)
    entries, order = parse_catalog(args.catalog)

    if args.results:
        with open(args.results, encoding="utf-8") as fh:
            data = json.load(fh)
        selected = set(data.get("recognized_ids", []))
        scope = "recognized"
    else:
        selected = set(order)
        scope = "catalog capacity (all titles)"

    # Count sources per shape, and remember which titles could close a gap.
    counts = {k: 0 for k in levity}
    carriers = {k: [] for k in levity}          # recognized titles carrying shape
    potential = {k: [] for k in levity}         # any catalog title carrying shape
    drift = {}                                  # raw tag -> [titles]

    for cid in order:
        e = entries[cid]
        for norm, raw in zip(e["shapes"], e.get("raw_shapes", e["shapes"])):
            if norm not in levity:
                drift.setdefault(raw, []).append(e["title"])
                continue
            potential[norm].append(e["title"])
            if cid in selected:
                counts[norm] += 1
                carriers[norm].append(e["title"])

    print(f"Coverage report — scope: {scope}, target: >= {args.target} per shape\n")
    gaps = 0
    for norm, display in sorted(levity.items(), key=lambda kv: kv[1].lower()):
        n = counts[norm]
        ok = n >= args.target
        if not ok:
            gaps += 1
        print(f"[{'OK ' if ok else 'GAP'}] {display}: {n}")
        if not ok:
            # Suggest titles that carry this shape but aren't in the selected set.
            missing = [t for t in potential[norm] if t not in carriers[norm]]
            if missing:
                print(f"        could add: {', '.join(sorted(set(missing)))}")
            else:
                print("        no catalog title carries this shape — extend show-catalog.md")

    if drift:
        print("\nDrift — catalog shape tags with no match in levity.md:")
        for raw, titles in sorted(drift.items()):
            print(f"  {raw!r} (in: {', '.join(sorted(set(titles)))})")

    print(f"\n{len(levity)} shapes, {gaps} gap(s), {len(drift)} drift tag(s).")
    return 1 if (gaps or drift) else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
