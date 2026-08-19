"""Programmatic usage of the Disassembly.AI toolkit.

    pip install disassembly
    DISASSEMBLY_API_KEY=... python scan.py https://staging.example.com
"""
import os
import sys

from disassembly.engine import Engine
from disassembly.report import to_sarif, to_markdown


def main() -> int:
    target = sys.argv[1] if len(sys.argv) > 1 else "https://staging.example.com"
    if not os.environ.get("DISASSEMBLY_API_KEY"):
        print("set DISASSEMBLY_API_KEY", file=sys.stderr)
        return 2

    engine = Engine()  # defaults to the Claude-backed provider
    findings = engine.scan(target)

    with open("report.sarif", "w") as f:
        f.write(to_sarif(findings))
    print(to_markdown(findings))
    print(
        f"usage: tokens={engine.meter.total_tokens} runs={engine.meter.total_runs}",
        file=sys.stderr,
    )
    # Non-zero exit if any high-severity finding — fail the build.
    return 1 if any(f.severity == "error" for f in findings) else 0


if __name__ == "__main__":
    raise SystemExit(main())
