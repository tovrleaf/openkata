"""Write the agent-workflow diagrams as SVG files.

This script made the first versions of the four SVGs in
web/docs/architecture/agent-workflow/. The SVG files are the source:
after an SVG is edited by hand, running this script would overwrite the
edit. It is kept as a starting point for new diagrams, and spec 0030
removes it once Mermaid renders the diagrams.

Run from the repository root:

    python3 -I specs/0025-static-assets-and-diagrams/generate_diagrams.py \\
      web/docs/architecture/agent-workflow
"""
import pathlib
import sys

LIGHT = {
    "surface": "#fcfcfc", "raised": "#fff", "ink": "#0a0a0a",
    "muted": "#6b6b6b", "accent": "#4a7fa5",
}
DARK = {
    "surface": "#1a1a2e", "raised": "#2a2a3e", "ink": "#f0eeeb",
    "muted": "#a0a0a0", "accent": "#5a9fd4",
}


def rules(p, c):
    return f"""\
    #{p} {{ font-family: var(--font-body, system-ui, sans-serif); font-size: 14px; }}
    #{p} .box {{ fill: var(--color-surface-raised, {c['raised']}); stroke: var(--color-ink-muted, {c['muted']}); stroke-width: 1.25; }}
    #{p} .end {{ fill: var(--color-surface, {c['surface']}); stroke: var(--color-ink, {c['ink']}); stroke-width: 1.25; }}
    #{p} .choice {{ fill: var(--color-surface-raised, {c['raised']}); stroke: var(--color-accent-muted, {c['accent']}); stroke-width: 1.5; }}
    #{p} .edge {{ fill: none; stroke: var(--color-ink-muted, {c['muted']}); stroke-width: 1.25; }}
    #{p} .head {{ fill: var(--color-ink-muted, {c['muted']}); }}
    #{p} text {{ fill: var(--color-ink, {c['ink']}); text-anchor: middle; }}
    #{p} .note {{ fill: var(--color-ink-muted, {c['muted']}); font-size: 12px; }}"""


def svg(p, title, w, h, body):
    return f"""\
<svg xmlns="http://www.w3.org/2000/svg" id="{p}" viewBox="0 0 {w} {h}" width="{round(w * 8 / 7)}" role="img" aria-labelledby="{p}-title">
  <title id="{p}-title">{title}</title>
  <style>
{rules(p, LIGHT)}
    #{p} .strong {{ font-weight: 600; }}
    #{p} .dashed {{ stroke-dasharray: 5 4; }}
    #{p} .start {{ text-anchor: start; }}
    @media (prefers-color-scheme: dark) {{
{rules(p, DARK)}
    }}
  </style>
  <defs>
    <marker id="{p}-arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
      <path class="head" d="M0,0 L10,5 L0,10 z"/>
    </marker>
  </defs>
{body}
</svg>
"""


def arrow(p, d, extra=""):
    cls = "edge" + (" " + extra if extra else "")
    return f'  <path class="{cls}" d="{d}" marker-end="url(#{p}-arrow)"/>'


def line(d):
    return f'  <path class="edge" d="{d}"/>'


def box(x, y, w, h, lines, cls="box", rx=3):
    """A rectangle with one or two centred lines. Two lines put the
    first in bold."""
    cx = x + w / 2
    out = [f'  <rect class="{cls}" x="{x}" y="{y}" width="{w}" height="{h}" rx="{rx}"/>']
    if len(lines) == 1:
        out.append(f'  <text x="{cx:g}" y="{y + h / 2 + 5:g}">{lines[0]}</text>')
    else:
        mid = y + h / 2
        out.append(f'  <text class="strong" x="{cx:g}" y="{mid - 4:g}">{lines[0]}</text>')
        out.append(f'  <text class="note" x="{cx:g}" y="{mid + 14:g}">{lines[1]}</text>')
    return "\n".join(out)


def pill(x, y, w, h, label):
    return box(x, y, w, h, [label], cls="end", rx=h / 2)


def diamond(cx, cy, hw, hh, lines):
    pts = f"{cx},{cy - hh} {cx + hw},{cy} {cx},{cy + hh} {cx - hw},{cy}"
    out = [f'  <polygon class="choice" points="{pts}"/>']
    if len(lines) == 1:
        out.append(f'  <text x="{cx}" y="{cy + 5}">{lines[0]}</text>')
    else:
        out.append(f'  <text x="{cx}" y="{cy - 4}">{lines[0]}</text>')
        out.append(f'  <text x="{cx}" y="{cy + 14}">{lines[1]}</text>')
    return "\n".join(out)


def note(x, y, label, start=False):
    cls = "note start" if start else "note"
    return f'  <text class="{cls}" x="{x}" y="{y}">{label}</text>'


def feature_development():
    p = "aw-feature"
    body = "\n".join([
        box(100, 10, 180, 52, ["1. Plan", "spec-planner"]),
        box(100, 102, 180, 52, ["2. Implement", "frontend + backend"]),
        box(100, 194, 180, 52, ["3. Validate", "spec-validator"]),
        pill(130, 296, 120, 40, "Done"),
        arrow(p, "M190,62 V102"),
        arrow(p, "M190,154 V194"),
        arrow(p, "M190,246 V296"),
        note(198, 275, "pass", start=True),
        arrow(p, "M280,220 H330 V128 H280"),
        note(338, 178, "fail", start=True),
    ])
    title = ("Feature development: plan with spec-planner, implement with "
             "the frontend and backend developers, then validate with "
             "spec-validator. A pass is done; a fail goes back to implement.")
    return svg(p, title, 380, 346, body)


def parallel_implementation():
    p = "aw-parallel"
    body = "\n".join([
        box(150, 10, 160, 40, ["dojo-master"]),
        box(35, 100, 170, 40, ["frontend-developer"]),
        box(255, 100, 170, 40, ["backend-developer"]),
        arrow(p, "M230,50 V72 H120 V100"),
        arrow(p, "M230,72 H340 V100"),
        line("M120,140 V162"),
        note(120, 180, "templates, CSS, JS"),
        arrow(p, "M120,192 V240 H160"),
        line("M340,140 V162"),
        note(340, 180, "handlers, types, tests"),
        arrow(p, "M340,192 V240 H300"),
        box(160, 220, 140, 40, ["Converge"]),
        arrow(p, "M230,260 V300"),
        box(115, 300, 230, 52, ["Integration tasks", "both agents sequential"]),
    ])
    title = ("Parallel implementation: dojo-master hands templates, CSS and "
             "JS to frontend-developer, and handlers, types and tests to "
             "backend-developer. Their work converges into integration tasks "
             "that both agents do in sequence.")
    return svg(p, title, 460, 362, body)


def agent_responsibilities():
    p = "aw-agents"
    parts = [
        pill(170, 10, 100, 36, "User"),
        arrow(p, "M200,46 V70 H120 V96"),
        arrow(p, "M240,46 V70 H330 V96", "dashed"),
        note(285, 62, "direct 1-1"),
        box(40, 96, 160, 40, ["dojo-master"]),
        box(255, 96, 150, 40, ["spec-planner"]),
        line("M70,136 V362"),
    ]
    agents = ["frontend-developer", "backend-developer", "spec-validator",
              "adr-author", "kata-author"]
    for i, name in enumerate(agents):
        y = 160 + 46 * i
        parts.append(arrow(p, f"M70,{y + 18} H100"))
        parts.append(box(100, y, 180, 36, [name]))
    title = ("Agent responsibilities: the user works through dojo-master, "
             "or one to one with spec-planner. dojo-master delegates to "
             "frontend-developer, backend-developer, spec-validator, "
             "adr-author and kata-author.")
    return svg(p, title, 440, 390, "\n".join(parts))


def routing():
    p = "aw-routing"
    parts = [
        pill(45, 10, 110, 36, "Request"),
        arrow(p, "M100,46 V72"),
        diamond(100, 100, 60, 28, ["Planning?"]),
        arrow(p, "M160,100 H230"),
        note(195, 93, "yes"),
        box(230, 80, 220, 40, ["Switch to spec-planner"]),
        arrow(p, "M100,128 V166"),
        note(108, 152, "no", start=True),
        box(20, 166, 160, 40, ["dojo-master"]),
        line("M30,206 V570"),
    ]
    rows = [
        (["Templates,", "CSS, JS?"], "frontend-developer"),
        (["Handlers,", "types, tests?"], "backend-developer"),
        (["Decision?"], "adr-author"),
        (["Convention?"], "kata-author"),
        (["Validate?"], "spec-validator"),
    ]
    for i, (question, agent) in enumerate(rows):
        cy = 250 + 80 * i
        parts.append(arrow(p, f"M30,{cy} H55"))
        parts.append(diamond(150, cy, 95, 34, question))
        parts.append(arrow(p, f"M245,{cy} H300"))
        parts.append(note(272, cy - 7, "yes"))
        parts.append(box(300, cy - 20, 170, 40, [agent]))
    title = ("Routing: a planning request switches to spec-planner. "
             "Anything else goes to dojo-master, which sends templates, CSS "
             "and JS to frontend-developer, handlers, types and tests to "
             "backend-developer, decisions to adr-author, conventions to "
             "kata-author and validation to spec-validator.")
    return svg(p, title, 480, 612, "\n".join(parts))


def main():
    out = pathlib.Path(sys.argv[1])
    out.mkdir(parents=True, exist_ok=True)
    for name, make in [
        ("feature-development.svg", feature_development),
        ("parallel-implementation.svg", parallel_implementation),
        ("agent-responsibilities.svg", agent_responsibilities),
        ("routing.svg", routing),
    ]:
        (out / name).write_text(make())
        print(out / name)


if __name__ == "__main__":
    main()
