#!/usr/bin/env python3
"""将 go test -json 原文提交到评分服务；没有配置 URL 时只打印本地结果。"""
import json
import os
import sys
import urllib.request
from pathlib import Path

if len(sys.argv) != 2:
    raise SystemExit("usage: publish-result.py test-results.json")
result = Path(sys.argv[1]).read_text(encoding="utf-8")
url = os.environ.get("SCORE_API_URL", "").strip()
if not url:
    print("SCORE_API_URL is not configured; keeping local test artifact")
    raise SystemExit(0)
payload = {
    "student": os.environ.get("STUDENT", "unknown"),
    "exercise": os.environ.get("EXERCISE", "unknown"),
    "test_output": result,
}
request = urllib.request.Request(url, data=json.dumps(payload).encode(), headers={"Content-Type": "application/json"})
token = os.environ.get("SCORE_API_TOKEN", "")
if token:
    request.add_header("Authorization", "Bearer " + token)
with urllib.request.urlopen(request, timeout=15) as response:
    print(response.read().decode())
