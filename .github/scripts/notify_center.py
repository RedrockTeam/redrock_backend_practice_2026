"""把本次自测报告推送到中心站点。

站点只接收与展示报告、不做复算，所以这里是「尽力而为」的上传：
失败会退避重试，但无论成败都不影响流水线结论（报告另有 artifact 兜底）。

站点开启上报鉴权后，请求需携带 GitHub Actions OIDC 令牌（audience
homework-center）；取不到令牌时降级为不带令牌上报，语义不变。
"""

import json
import os
import time
import urllib.error
import urllib.request
from pathlib import Path

TIMEOUT_SECONDS = 10
ATTEMPTS = 3
BACKOFF_SECONDS = (1, 3)
REPORTS_PATH = "/api/v1/reports"
OIDC_AUDIENCE = "homework-center"


def optional_json(path):
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError:
        return None


def repo_url():
    """仓库地址由 CI 环境推导，不依赖学生自填。"""
    repo = os.environ.get("GITHUB_REPOSITORY", "")
    if not repo:
        return ""
    server = os.environ.get("GITHUB_SERVER_URL", "https://github.com")
    return server.rstrip("/") + "/" + repo


def build_config():
    """config.json 去掉 center 后原样上报。

    center 是中心站点自己的地址，把它回传给站点没有意义。
    """
    config = optional_json(Path("config.json"))
    if config is None:
        return None
    config.pop("center", None)
    return config


def build_result():
    """把 result.json（schema 2.0）压成上报用的课次结果。

    看板只需要两件事：这一课的提交、以及每道题有没有做完。
    所以只留题目编号与结论——分数、用例名、包名、耗时都不进载荷，
    它们留在 result.json 与 artifact 里，需要时去那里取。

    课次由 LESSON_PATH 指定，取 lessons 里同名的那一项。
    """
    data = optional_json(Path("result.json"))
    if data is None:
        return None

    lesson = os.environ.get("LESSON_PATH", "").strip()
    entries = data.get("lessons") or []
    chosen = next((item for item in entries if item.get("lesson") == lesson), None)
    if chosen is None:
        chosen = entries[0] if entries else None
    if chosen is None:
        return None

    return {
        "lesson": chosen.get("lesson", lesson),
        "tests": [
            {"test": entry.get("test", "-"), "status": entry.get("status", "fail")}
            for entry in chosen.get("tests", [])
        ],
    }


def build_payload():
    payload = {
        "repo_url": repo_url(),
        "commit": os.environ.get("GITHUB_SHA", ""),
        "ref": os.environ.get("GITHUB_REF", ""),
        "event": os.environ.get("GITHUB_EVENT_NAME", ""),
    }

    config = build_config()
    if config is not None:
        payload["config"] = config

    result = build_result()
    if result is not None:
        payload["result"] = result

    return payload


def announce(line):
    print(line)
    summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if summary:
        with open(summary, "a", encoding="utf-8") as handle:
            handle.write(line + "\n")


def get_oidc_token():
    """换取 GitHub Actions OIDC 令牌；不在 Actions 环境或缺权限时返回 None。"""
    url = os.environ.get("ACTIONS_ID_TOKEN_REQUEST_URL")
    req_token = os.environ.get("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
    if not url or not req_token:
        print("[warn] runner 未提供 OIDC 环境变量，上报将不带令牌")
        return None
    try:
        request = urllib.request.Request(
            url + "&audience=" + OIDC_AUDIENCE,
            headers={"Authorization": "Bearer " + req_token},
        )
        with urllib.request.urlopen(request, timeout=TIMEOUT_SECONDS) as response:
            value = json.load(response)["value"]
            print("[info] 已取得 OIDC 令牌（audience=homework-center）")
            return value
    except Exception as error:
        print(f"[warn] 取 OIDC 令牌失败: {error}")
        return None


def upload(url, payload):
    for attempt in range(1, ATTEMPTS + 1):
        # 每次尝试都取一枚新令牌：令牌约 5 分钟时效、不值得缓存；
        # 收到 401 后下一轮自然会用重新换取的令牌重试。
        headers = {"Content-Type": "application/json"}
        token = get_oidc_token()
        if token:
            headers["Authorization"] = "Bearer " + token
        request = urllib.request.Request(
            url,
            data=json.dumps(payload, ensure_ascii=False).encode("utf-8"),
            headers=headers,
            method="POST",
        )
        try:
            with urllib.request.urlopen(request, timeout=TIMEOUT_SECONDS) as response:
                return f"报告上传成功：HTTP {response.status}（第 {attempt} 次尝试）"
        except (urllib.error.URLError, OSError) as error:
            last_error = error
            if getattr(error, "code", None) == 401:
                print("[warn] 站点拒绝了令牌（HTTP 401），重试前将重新换取 OIDC 令牌")
            if attempt < ATTEMPTS:
                time.sleep(BACKOFF_SECONDS[attempt - 1])
    return (f"报告上传失败（已尝试 {ATTEMPTS} 次）：{last_error}"
            "；报告已随 artifact 留存，可事后补取。")


def main():
    base = os.environ.get("CENTER_API", "").strip()
    if not base:
        announce("未配置 CENTER_API，跳过报告上传。")
        return

    announce(upload(base.rstrip("/") + REPORTS_PATH, build_payload()))


if __name__ == "__main__":
    main()
