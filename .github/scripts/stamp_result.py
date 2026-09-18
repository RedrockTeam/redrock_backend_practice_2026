#!/usr/bin/env python3
"""Overwrite a result bundle's identity fields with values GitHub vouches for.

result.json is produced inside a fork's pull_request run, so every field in it
is attacker-controlled: a student can claim any score, any commit, any
classmate's name. This script runs in the base repository under workflow_run,
where the event payload is trustworthy, and:

  * replaces the whole `submission` block with values from the event,
  * recomputes `result_id` from those trusted values,
  * verifies the claimed pull request really points at the same commit,
  * recomputes `status`, `score` and every `summary` from the per-test verdicts
    so a hand-edited total cannot disagree with the verdicts it is derived from.

What stays student-controlled is the list of test verdicts itself, which is
inherent: the tests run on their machine image. Anything derived from those
verdicts is recomputed here rather than trusted.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import subprocess
import sys
from pathlib import Path
from typing import Any

VALID_STATUSES = {"pass", "fail", "skip"}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--bundle", required=True, help="Directory the artifact was extracted to")
    parser.add_argument("--output", required=True, help="Where to write the stamped payload")
    return parser.parse_args()


def count(cases: list[dict[str, Any]]) -> dict[str, int]:
    return {
        "total": len(cases),
        "passed": sum(case["status"] == "pass" for case in cases),
        "failed": sum(case["status"] == "fail" for case in cases),
        "skipped": sum(case["status"] == "skip" for case in cases),
    }


def score_of(counts: dict[str, int]) -> int:
    return counts["passed"] * 100 // counts["total"] if counts["total"] else 0


def verified_pull_request(bundle: Path, trusted_commit: str) -> int | None:
    """The PR number from the bundle, kept only if its head commit matches."""
    number_file = bundle / "pull-request-number.txt"
    if not number_file.is_file():
        return None
    raw = number_file.read_text(encoding="utf-8").strip()
    if not raw.isdigit():
        return None
    number = int(raw)

    repository = os.environ["TRUSTED_REPOSITORY"]
    try:
        head_sha = subprocess.run(
            [
                "gh", "api",
                f"repos/{repository}/pulls/{number}",
                "--jq", ".head.sha",
            ],
            capture_output=True,
            text=True,
            check=True,
            timeout=30,
        ).stdout.strip()
    except (subprocess.CalledProcessError, subprocess.TimeoutExpired, FileNotFoundError) as error:
        print(f"::warning::无法核验 PR #{number}：{error}", file=sys.stderr)
        return None

    if head_sha != trusted_commit:
        print(
            f"::warning::PR #{number} 的 head 是 {head_sha}，与本次运行的 {trusted_commit} 不一致，已丢弃该 PR 号",
            file=sys.stderr,
        )
        return None
    return number


def rebuild_totals(result: dict[str, Any]) -> None:
    """Recompute every aggregate from the leaf verdicts, in place."""
    all_cases: list[dict[str, Any]] = []

    for lesson in result.get("lessons", []):
        lesson_cases: list[dict[str, Any]] = []

        for test in lesson.get("tests", []):
            test_cases: list[dict[str, Any]] = []
            for package in test.get("packages", []):
                cases = [
                    case
                    for case in package.get("cases", [])
                    if isinstance(case, dict) and case.get("status") in VALID_STATUSES
                ]
                package["cases"] = cases
                package["status"] = "pass" if cases and all(c["status"] != "fail" for c in cases) else "fail"
                test_cases.extend(cases)

            test["summary"] = count(test_cases)
            test["score"] = score_of(test["summary"])
            test["status"] = (
                "pass"
                if test.get("packages") and all(p["status"] == "pass" for p in test["packages"])
                else "fail"
            )
            lesson_cases.extend(test_cases)

        lesson["summary"] = count(lesson_cases)
        lesson["score"] = score_of(lesson["summary"])
        lesson["status"] = (
            "pass"
            if lesson.get("tests") and all(t["status"] == "pass" for t in lesson["tests"])
            else "fail"
        )
        all_cases.extend(lesson_cases)

    counts = count(all_cases)
    result["summary"] = {**counts, "packages": result.get("summary", {}).get("packages", 0)}
    result["score"] = score_of(counts)
    result["status"] = "pass" if counts["total"] > 0 and counts["failed"] == 0 else "fail"


def main() -> int:
    args = parse_args()
    bundle = Path(args.bundle)
    result_file = bundle / "result.json"

    if not result_file.is_file():
        print(f"::error::artifact 里没有 result.json（{result_file}）", file=sys.stderr)
        return 1
    try:
        result = json.loads(result_file.read_text(encoding="utf-8"))
    except json.JSONDecodeError as error:
        print(f"::error::result.json 不是合法 JSON：{error}", file=sys.stderr)
        return 1
    if not isinstance(result, dict) or not isinstance(result.get("lessons"), list):
        print("::error::result.json 结构不符合预期", file=sys.stderr)
        return 1

    repository = os.environ["TRUSTED_REPOSITORY"]
    commit = os.environ["TRUSTED_COMMIT"]
    run_id = os.environ["TRUSTED_RUN_ID"]
    run_attempt = os.environ.get("TRUSTED_RUN_ATTEMPT", "1")
    fingerprint = f"{repository}:{commit}:{run_id}:{run_attempt}"

    result["submission"] = {
        "result_id": hashlib.sha256(fingerprint.encode("utf-8")).hexdigest()[:16],
        "repository": repository,
        "head_repository": os.environ.get("TRUSTED_HEAD_REPOSITORY") or repository,
        "commit": commit,
        "branch": os.environ.get("TRUSTED_BRANCH", ""),
        "pull_request": verified_pull_request(bundle, commit),
        "actor": os.environ.get("TRUSTED_ACTOR", ""),
        "event": os.environ.get("TRUSTED_EVENT", ""),
        "run_id": run_id,
        "run_attempt": int(run_attempt) if run_attempt.isdigit() else 1,
        "run_url": os.environ.get("TRUSTED_RUN_URL", ""),
        "workflow_conclusion": os.environ.get("TRUSTED_CONCLUSION", ""),
        "verified": True,
    }
    rebuild_totals(result)

    Path(args.output).write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )

    github_output = os.getenv("GITHUB_OUTPUT")
    if github_output:
        with Path(github_output).open("a", encoding="utf-8") as handle:
            handle.write(f"result_id={result['submission']['result_id']}\n")
            handle.write(f"status={result['status']}\n")
            handle.write(f"score={result['score']}\n")

    print(
        f"stamped result_id={result['submission']['result_id']} "
        f"status={result['status']} score={result['score']}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
