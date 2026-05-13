---
name: Skill 加载路径必须验证版本一致性
description: Agent 从 plugin cache 的旧版本路径读取 skill 模板，导致 eval 用了错误的 100 分 rubric 而非当前 1000 分 rubric。Trust-without-verify 模式。
type: feedback
---

## Rule

当 skill base directory 包含 `cache` 或版本号路径时，必须验证路径指向的版本与用户当前安装版本一致。不一致时从用户实际安装目录读取。

## Why

执行 `/eval-proposal` 时，skill 加载的 base directory 指向 plugin cache 路径 `...cache/forge/forge/2.18.0/skills/eval-proposal`，而非用户的实际 forge 安装目录 `Z:\project\ai\forge\plugins\forge\skills\eval-proposal`（3.0.0-beta-4）。

Agent 信任了 skill 加载机制给出的路径，直接从 cache 读取了 2.18.0 版本的 rubric（100 分 / 6 维度），而用户当前版本是 1000 分 / 10 维度。全部 3 轮 eval 都基于错误的 rubric 运行，评分和 attack points 完全无效。

**决策链回溯：**

| 层级 | 问题 | 本应怎么做 |
|------|------|-----------|
| Symptom | 3 轮 eval 用了 100 分 rubric 而非 1000 分 | N/A |
| Direct cause | 从 `cache/forge/forge/2.18.0/` 路径读取模板 | 检查路径中的版本号是否匹配用户当前版本 |
| Root cause | **Trust-without-verify**: 信任 skill 加载机制提供的 base dir 是最新版本 | 主动对比 base dir 版本与用户声明的版本 |

**路径中的信号（当时被忽略）：**
- `cache` 关键词 → 这是缓存目录，不是主安装目录
- `2.18.0` 版本号 → 具体版本，应与用户版本对比

## How to apply

1. **Skill 模板路径检查**：读取任何 skill 模板文件前，检查 base directory 是否包含 `cache` 或版本号。如果是，优先从用户实际安装路径读取。
2. **版本对比**：如果用户提到了 forge 版本（或可通过 `ls`/`cat package.json` 确认），对比 base dir 中的版本号。
3. **Glob 搜索的陷阱**：不要用宽泛的 Glob 搜索定位模板文件。Skill 定义已给出相对路径，应从 base dir 直接拼接。即使 base dir 是错的，直接拼接也比 Glob 更容易暴露路径问题（读不到文件会报错）。
4. **适用范围**：此规则适用于所有通过 skill base directory 加载的模板文件（rubric、report template、task template 等）。
