"""校验 config.json，解析出课次与中心站点地址并写入 GITHUB_ENV。

课次是计分的归属，必须明确：缺失或非法都直接报错并中止本 job，不进入测试与上传。
其余问题（姓名留空、站点地址缺失）只输出 ::warning::，不影响后面的步骤。

课次以**目录**的形式交给 run_go_tests.py 的 --root：那一课下的 testNN/ 全都会跑到，
但其它课次不受影响——报告只归属这一个课次。
"""

import json
import os
from pathlib import Path

CONFIG_FILE = "config.json"
REQUIRED_FIELDS = ("name", "center")
URL_PREFIXES = ("http://", "https://")


def available_lessons(root):
    return sorted(path.name for path in root.glob("lesson-*") if path.is_dir())


def warn(message):
    print("::warning::" + message)


def fail(message):
    print("::error::" + message)
    raise SystemExit(1)


def load_config(root):
    path = root / CONFIG_FILE
    if not path.is_file():
        return {}
    try:
        loaded = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as error:
        warn(CONFIG_FILE + " 不是合法 JSON：" + str(error))
        return {}
    if not isinstance(loaded, dict):
        warn(CONFIG_FILE + " 顶层必须是 JSON 对象。")
        return {}
    return loaded


def resolve_lesson(data, lessons):
    """课次必须明确。缺失或非法都直接报错，不做任何回退。"""
    requested = str(data.get("lesson", "")).strip()
    if not requested:
        fail(CONFIG_FILE + " 没有填写 lesson。请在 lesson 里填写本次完成的课次后重新 push；"
             "可选值：" + "、".join(lessons))
    if requested not in lessons:
        fail("课次「" + requested + "」不存在。请改成下列之一后重新 push：" + "、".join(lessons))
    return requested


def resolve_center(data):
    center = str(data.get("center", "")).strip()
    if center and not center.startswith(URL_PREFIXES):
        warn("center 应是完整 URL，例如 https://grade.example.com")
        return ""
    return center


def write_env(values):
    env_file = os.environ.get("GITHUB_ENV")
    if not env_file:
        print("（未检测到 GITHUB_ENV，本地调试模式，不写入环境变量）")
        return
    with open(env_file, "a", encoding="utf-8") as handle:
        for key, value in values.items():
            handle.write(key + "=" + value + "\n")


def main():
    root = Path.cwd()
    lessons = available_lessons(root)
    data = load_config(root)

    if not data:
        fail("读不到 " + CONFIG_FILE + " 的内容。请补上 " + CONFIG_FILE
             + "，填写 name 与 lesson 后重新 push。")

    missing = [field for field in REQUIRED_FIELDS if not str(data.get(field, "")).strip()]
    if missing:
        warn(CONFIG_FILE + " 缺少字段：" + "、".join(missing))

    lesson = resolve_lesson(data, lessons)
    center = resolve_center(data)

    print("本次自测课次：" + lesson)
    print("自测范围：" + lesson + "/ 下的全部 testNN")
    print("中心站点：" + (center if center else "未配置，本次不会上传报告"))

    write_env({"LESSON_PATH": lesson, "TEST_ROOT": lesson, "CENTER_API": center})


if __name__ == "__main__":
    main()
