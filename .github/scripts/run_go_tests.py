#!/usr/bin/env python3
"""Run every Go test package in the repository and emit a results-only result.json.

Layout assumption: <lesson>/<test unit>/<go package>/..., for example
lesson-01-basics/test06/roster. Suites are attributed to their lesson and test
unit so the grading website can address a single exercise by route.

Two files are produced:

  result.json  the uploaded payload: verdicts only, no captured output
  score.md     a human-readable table for the GitHub Step Summary

Raw `go test -json` streams stay in --raw-dir and never enter result.json; they
are attached to the workflow run as an artifact for debugging.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import subprocess
import sys
import time
from collections import OrderedDict
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

TERMINAL_ACTIONS = {"pass", "fail", "skip"}
EXCLUDED_DIRS = {".git", ".github", "testdata", "vendor", "node_modules"}
LESSON_PATTERN = re.compile(r"^lesson-\d+")
TEST_UNIT_PATTERN = re.compile(r"^test\d+$")
MAX_RAW_OUTPUT_CHARS = 200_000


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--root",
        default=".",
        help="Directory to discover test packages in; defaults to the whole repository",
    )
    parser.add_argument("--result", default="result.json", help="Results-only payload")
    parser.add_argument("--summary", default="score.md", help="Markdown summary")
    parser.add_argument("--raw-dir", default="test-results", help="Raw go test -json streams")
    parser.add_argument(
        "--timeout",
        default="120s",
        help="Per-package go test timeout, passed through as -timeout",
    )
    return parser.parse_args()


def discover_test_packages(root: Path) -> list[Path]:
    """Every directory holding at least one *_test.go, excluding fixture trees."""
    packages: set[Path] = set()
    for test_file in root.rglob("*_test.go"):
        relative = test_file.relative_to(root)
        if any(part in EXCLUDED_DIRS for part in relative.parts):
            continue
        packages.add(test_file.parent)
    return sorted(packages, key=lambda item: item.as_posix())


def safe_name(relative_path: str) -> str:
    return re.sub(r"[^A-Za-z0-9._-]+", "__", relative_path).strip("_") or "root"


def attribute(relative_path: str) -> tuple[str, str]:
    """Map a package path to (lesson, test unit).

    lesson-01-basics/test06/roster -> ("lesson-01-basics", "test06")
    Anything outside that layout is grouped under "-" so it still gets graded.
    """
    parts = relative_path.split("/")
    lesson = parts[0] if parts and LESSON_PATTERN.match(parts[0]) else "-"
    test_unit = "-"
    if lesson != "-":
        for part in parts[1:]:
            if TEST_UNIT_PATTERN.match(part):
                test_unit = part
                break
    return lesson, test_unit


def parse_events(stdout: str) -> tuple[list[dict[str, Any]], str, bool]:
    """Collapse a `go test -json` stream into one verdict per test function.

    Returns (verdicts, package level output, build failed). A test that starts
    but never reaches a terminal action is counted as failed: that is what an
    unrecoverable runtime crash such as `concurrent map writes` looks like, and
    it must be attributed to the test that triggered it rather than to the
    whole package.
    """
    test_states: "OrderedDict[str, dict[str, Any]]" = OrderedDict()
    started: "OrderedDict[str, None]" = OrderedDict()
    package_output: list[str] = []
    build_failed = False

    for line in stdout.splitlines():
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            # Non-JSON lines are emitted straight through on toolchain errors.
            package_output.append(line + "\n")
            continue

        action = event.get("Action")
        if action in {"build-fail", "build-output"}:
            build_failed = build_failed or action == "build-fail"
            output = event.get("Output")
            if output:
                package_output.append(output)
            continue

        test_name = event.get("Test")
        if not test_name:
            output = event.get("Output")
            if output:
                package_output.append(output)
            continue

        # Subtests (Parent/Child) roll up into their parent's verdict.
        if "/" in test_name:
            continue
        if action == "run":
            started.setdefault(test_name, None)
        elif action in TERMINAL_ACTIONS:
            test_states[test_name] = {
                "name": test_name,
                "status": action,
                "ms": round(float(event.get("Elapsed", 0)) * 1000),
            }

    for test_name in started:
        if test_name not in test_states:
            test_states[test_name] = {"name": test_name, "status": "fail", "ms": 0}

    return list(test_states.values()), "".join(package_output), build_failed


def count(tests: list[dict[str, Any]]) -> dict[str, int]:
    return {
        "total": len(tests),
        "passed": sum(test["status"] == "pass" for test in tests),
        "failed": sum(test["status"] == "fail" for test in tests),
        "skipped": sum(test["status"] == "skip" for test in tests),
    }


def score_of(counts: dict[str, int]) -> int:
    return counts["passed"] * 100 // counts["total"] if counts["total"] else 0


def run_package(repo_root: Path, package_dir: Path, raw_dir: Path, timeout: str) -> dict[str, Any]:
    relative = package_dir.relative_to(repo_root).as_posix()
    command = ["go", "test", "-json", "-count=1", f"-timeout={timeout}", f"./{relative}"]

    started = time.monotonic()
    completed = subprocess.run(
        command,
        cwd=repo_root,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
        check=False,
    )
    duration_ms = round((time.monotonic() - started) * 1000)

    tests, package_output, build_failed = parse_events(completed.stdout)

    # Raw streams are debugging material only; they stay out of result.json.
    raw_file = raw_dir / f"{safe_name(relative)}.jsonl"
    raw_file.write_text(completed.stdout[:MAX_RAW_OUTPUT_CHARS], encoding="utf-8")
    if completed.stderr or package_output:
        log_file = raw_dir / f"{safe_name(relative)}.log"
        log_file.write_text((package_output + completed.stderr)[:MAX_RAW_OUTPUT_CHARS], encoding="utf-8")

    # A package that does not compile emits no test events at all; record it as
    # one synthetic failure so a single broken exercise cannot hide the others.
    if not tests:
        if build_failed or completed.returncode != 0:
            tests.append({"name": "__build__", "status": "fail", "ms": duration_ms})
        else:
            tests.append({"name": "__no_tests__", "status": "fail", "ms": duration_ms})

    lesson, test_unit = attribute(relative)
    counts = count(tests)
    return {
        "lesson": lesson,
        "test": test_unit,
        "package": relative,
        "status": "pass" if counts["failed"] == 0 and counts["total"] > 0 else "fail",
        "duration_ms": duration_ms,
        "summary": counts,
        "cases": tests,
    }


def group(packages: list[dict[str, Any]]) -> list[dict[str, Any]]:
    """Nest packages under lesson -> test unit, preserving discovery order."""
    lessons: "OrderedDict[str, OrderedDict[str, list[dict[str, Any]]]]" = OrderedDict()
    for package in packages:
        lessons.setdefault(package["lesson"], OrderedDict()).setdefault(package["test"], []).append(package)

    grouped: list[dict[str, Any]] = []
    for lesson_name, units in lessons.items():
        lesson_cases: list[dict[str, Any]] = []
        test_entries: list[dict[str, Any]] = []

        for unit_name, unit_packages in units.items():
            unit_cases = [case for package in unit_packages for case in package["cases"]]
            lesson_cases.extend(unit_cases)
            unit_counts = count(unit_cases)
            test_entries.append(
                {
                    "test": unit_name,
                    "status": "pass" if all(p["status"] == "pass" for p in unit_packages) else "fail",
                    "score": score_of(unit_counts),
                    "duration_ms": sum(p["duration_ms"] for p in unit_packages),
                    "summary": unit_counts,
                    "packages": [
                        {
                            "package": p["package"],
                            "status": p["status"],
                            "cases": p["cases"],
                        }
                        for p in unit_packages
                    ],
                }
            )

        lesson_counts = count(lesson_cases)
        grouped.append(
            {
                "lesson": lesson_name,
                "status": "pass" if all(entry["status"] == "pass" for entry in test_entries) else "fail",
                "score": score_of(lesson_counts),
                "summary": lesson_counts,
                "tests": test_entries,
            }
        )
    return grouped


def event_payload() -> dict[str, Any]:
    """The webhook payload for this run, or {} outside GitHub Actions."""
    event_path = os.getenv("GITHUB_EVENT_PATH")
    if not event_path or not Path(event_path).is_file():
        return {}
    try:
        payload = json.loads(Path(event_path).read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        return {}
    return payload if isinstance(payload, dict) else {}


def submission() -> dict[str, Any]:
    """Identity of this submission: who pushed what, where, from which fork."""
    repository = os.getenv("GITHUB_REPOSITORY", "local/local")
    commit = os.getenv("GITHUB_SHA", "local")
    run_id = os.getenv("GITHUB_RUN_ID", "local")
    run_attempt = os.getenv("GITHUB_RUN_ATTEMPT", "1")

    # On a pull_request run these come from the event payload: there is no
    # environment variable carrying the PR number or the fork's name.
    pull_request = event_payload().get("pull_request") or {}
    number = pull_request.get("number")
    head_repository = (pull_request.get("head") or {}).get("repo") or {}

    fingerprint = f"{repository}:{commit}:{run_id}:{run_attempt}"

    return {
        "result_id": hashlib.sha256(fingerprint.encode("utf-8")).hexdigest()[:16],
        "repository": repository,
        "head_repository": head_repository.get("full_name") or repository,
        "commit": commit,
        "ref": os.getenv("GITHUB_REF", ""),
        "branch": os.getenv("GITHUB_HEAD_REF") or os.getenv("GITHUB_REF_NAME", ""),
        "pull_request": number if isinstance(number, int) else None,
        "actor": os.getenv("GITHUB_ACTOR", ""),
        "event": os.getenv("GITHUB_EVENT_NAME", "local"),
        "run_id": run_id,
        "run_attempt": int(run_attempt) if run_attempt.isdigit() else 1,
        "run_url": (
            f"{os.getenv('GITHUB_SERVER_URL', 'https://github.com')}/{repository}"
            f"/actions/runs/{run_id}"
            if run_id != "local"
            else ""
        ),
    }


def write_markdown(path: Path, result: dict[str, Any]) -> None:
    summary = result["summary"]
    lines = [
        "### Go 练习评分",
        "",
        f"- 结果：**{result['status']}**",
        f"- 积分：**{result['score']}/100**",
        f"- 通过 {summary['passed']} / 失败 {summary['failed']} / 跳过 {summary['skipped']}"
        f"（共 {summary['total']} 个测试）",
        f"- result_id：`{result['submission']['result_id']}`",
        "",
    ]
    for lesson in result["lessons"]:
        lines += [
            f"#### {lesson['lesson']} — {lesson['status']} ({lesson['score']}/100)",
            "",
            "| 题目 | 结果 | 通过 | 失败 | 跳过 | 未通过的测试 |",
            "|---|---|---:|---:|---:|---|",
        ]
        for test in lesson["tests"]:
            counts = test["summary"]
            failures = [
                case["name"]
                for package in test["packages"]
                for case in package["cases"]
                if case["status"] == "fail"
            ]
            lines.append(
                f"| `{test['test']}` | {test['status']} | {counts['passed']} | "
                f"{counts['failed']} | {counts['skipped']} | "
                f"{', '.join(f'`{name}`' for name in failures) or '—'} |"
            )
        lines.append("")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_github_outputs(result: dict[str, Any], summary_path: Path) -> None:
    output_path = os.getenv("GITHUB_OUTPUT")
    if output_path:
        with Path(output_path).open("a", encoding="utf-8") as output:
            output.write(f"status={result['status']}\n")
            output.write(f"score={result['score']}\n")
            output.write(f"result_id={result['submission']['result_id']}\n")

    step_summary = os.getenv("GITHUB_STEP_SUMMARY")
    if step_summary and summary_path.exists():
        with Path(step_summary).open("a", encoding="utf-8") as handle:
            handle.write(summary_path.read_text(encoding="utf-8"))


def main() -> int:
    args = parse_args()
    repo_root = Path.cwd().resolve()
    test_root = (repo_root / args.root).resolve()

    if not test_root.is_dir():
        print(f"test root does not exist: {test_root}", file=sys.stderr)
        return 2
    if repo_root != test_root and repo_root not in test_root.parents:
        print("test root must stay inside the repository", file=sys.stderr)
        return 2

    raw_dir = Path(args.raw_dir)
    raw_dir.mkdir(parents=True, exist_ok=True)

    packages = discover_test_packages(test_root)
    if not packages:
        print(f"no *_test.go found under {test_root}", file=sys.stderr)
        return 2

    # Every package runs on every invocation; one failure never skips the rest.
    results = [run_package(repo_root, package, raw_dir, args.timeout) for package in packages]
    lessons = group(results)

    all_cases = [case for package in results for case in package["cases"]]
    counts = count(all_cases)
    status = "pass" if counts["failed"] == 0 and counts["total"] > 0 else "fail"

    result = {
        "schema_version": "2.0",
        "generated_at": datetime.now(timezone.utc).isoformat(timespec="seconds"),
        "status": status,
        "score": score_of(counts),
        "submission": submission(),
        "summary": {
            **counts,
            "packages": len(results),
            "duration_ms": sum(package["duration_ms"] for package in results),
        },
        "lessons": lessons,
    }

    result_path = Path(args.result)
    result_path.write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )
    summary_path = Path(args.summary)
    write_markdown(summary_path, result)
    write_github_outputs(result, summary_path)

    print(summary_path.read_text(encoding="utf-8"))
    return 0 if status == "pass" else 1


if __name__ == "__main__":
    raise SystemExit(main())
