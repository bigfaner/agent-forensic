---
date: "2026-05-09"
doc_dir: "docs/proposals/agent-forensic/"
iteration: 3
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 87/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  19      │  20      │ ✅         │
│    Problem clarity           │  7/7     │          │            │
│    Evidence provided         │  6/7     │          │            │
│    Urgency justified         │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  19      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  6/7     │          │            │
│    Differentiated            │  6/6     │          │            │
├──────────────────────────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  14      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  5/5     │          │            │
│    Scope bounded             │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅         │
│    Risks identified (≥3)     │  5/5     │          │            │
│    Likelihood + impact rated │  5/5     │          │            │
│    Mitigations actionable    │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Success Criteria          │  8       │  15      │ ⚠️         │
│    Measurable                │  4/5     │          │            │
│    Coverage complete         │  2/5     │          │            │
│    Testable                  │  2/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  87      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Evidence (line 17) | Evidence remains self-reported and anecdotal: "手动翻阅 ~/.claude/ JSONL 效率极低" — no external user survey, community forum complaints, or third-party validation. The numbers are specific but come from a single author's experience. | -1 pt from Evidence provided |
| User-facing behavior (lines 39-72) | The ASCII mockup and workflow are strong additions. However, the mockup shows only a single happy-path state. No description of empty states (no sessions found), error states (corrupt JSONL), or search-no-results behavior. A complete user-facing description must cover edge-case UI. | -1 pt from User-facing behavior |
| Alternatives — Web UI (line 81) | "依赖浏览器，不符合终端工作流习惯" is an assertion about user preference, not an analyzed trade-off. No effort to estimate implementation cost, compare user reach, or acknowledge that many developers already use browser-based tools alongside terminals. | -1 pt from Pros/cons honest |
| Alternatives — Rationale (line 83) | "lazygit 已验证了 TUI 三面板交互模式的可行性和用户接受度" — this is a good argument but lacks specificity: what specifically did lazygit validate? User count? Retention rate? The "1.5-2 周" estimate vs "3-4 周" for VS Code is a relative comparison but neither estimate is backed by a work breakdown. | -1 pt from Rationale justified |
| Scope bounded (lines 96, 99-101) | The Phase 1 / Phase 2 split for AI features is a clear improvement. However, "AI 证据提取（Phase 1 / MVP）" is still ambitious — it requires "自动提取关键证据（调用链 + thinking 片段 + 越权操作）" with "每条证据标注 JSONL 行号，覆盖 100% 已标记异常点". The 100% coverage requirement for automated evidence extraction is itself scope-risky; no fallback is described if the extraction logic misses edge cases. The boundary between "evidence extraction" and "root cause analysis" is stated but inherently fuzzy — when does "展示调用链 + thinking 片段" cross into analysis? | -1 pt from Scope bounded |
| Risk — Mitigations actionable (line 119) | User adoption risk mitigation: "发布 2 周内收集 ≥5 位用户周活跃数据验证" — this is a measurement plan, not a mitigation. If adoption is low after 2 weeks, what action is taken? The mitigation should describe both the measurement AND the contingency (e.g., "if <5 users, simplify onboarding and add pipe mode"). | -1 pt from Mitigations actionable |
| Success Criteria — Measurable (line 134) | "覆盖 100% 已标记异常点" sounds measurable but "已标记异常点" is circular — the system itself does the marking. If the anomaly detection misses a case, it won't be marked, and therefore won't need to be covered. This criterion cannot fail as written. | -1 pt from Measurable |
| Success Criteria — Coverage (lines 126-134) | Multiple in-scope items lack corresponding success criteria: (1) "统计仪表盘" (in-scope item 4) has only an accuracy criterion but no criterion for what metrics/charts it must display; (2) "事后回放" (in-scope item 3) has no criterion for replay navigation (e.g., time-axis scrolling, step-forward/backward); (3) "敏感内容脱敏" is mentioned in the UI description (line 62) but has no success criterion. | -3 pts from Coverage complete |
| Success Criteria — Testable (lines 126-134) | "调用树展示 ≥3 层嵌套" is testable but vague on completeness — what constitutes a "layer"? The sub-agent line says "显示调用次数 + 总耗时概要" but no criterion verifies the sub-agent summary rendering works. "搜索结果 500ms 内返回" is testable but the search feature itself (what fields are searchable, what syntax) is underspecified, making the test ambiguous. The SHA256 criterion (line 133) is excellent but it is the only criterion with a fully automated test path. | -3 pts from Testable |

---

## Attack Points

### Attack 1: Success Criteria — Coverage is incomplete, leaving 3 in-scope features unverifiable

**Where**: In-scope items: "统计仪表盘：工具/Skill 调用次数、任务总耗时、各步骤耗时占比" (line 94), "事后回放：加载历史会话，按时间轴浏览" (line 93), Detail panel: "敏感内容（API_KEY / SECRET / TOKEN / PASSWORD）自动脱敏" (line 62) vs Success Criteria (lines 126-134) which lack criteria for these features.

**Why it's weak**: Three in-scope deliverables have no success criterion that would allow a reviewer to declare them "done." The dashboard accuracy criterion ("工具调用计数误差 0，耗时误差 ≤1 秒") only checks data correctness, not that the dashboard actually renders the described charts (调用次数分布, 耗时占比). The replay feature has no criterion for the time-axis navigation experience. The data-sensitive redaction described in the UI section is completely absent from success criteria. A team could ship a dashboard with correct numbers but wrong visualizations, a replay mode with no navigation controls, and no redaction — and still pass every criterion.

**What must improve**: Add one criterion per missing in-scope item: (1) "统计仪表盘渲染 ≥3 种图表（调用次数柱状图、耗时占比饼图、时间线），每种图表数据与 JSONL 原文一致"; (2) "事后回放支持按 Turn 前进/后退，每次跳转 ≤200ms"; (3) "Detail 面板中匹配 API_KEY|SECRET|TOKEN|PASSWORD 的字符串自动替换为 \*\*\*"。

### Attack 2: Success Criteria — "100% coverage of marked anomalies" is circular and untestable

**Where**: "异常标记检出率 100%、误标率 ≤5%：耗时 >30 秒标黄色，访问项目外路径标红色" (line 130) and "AI 证据提取（Phase 1）：...覆盖 100% 已标记异常点" (line 134)

**Why it's weak**: The anomaly detection system itself determines what gets "marked." The "检出率 100%" criterion says every marked anomaly must be detected — but if the detection logic fails to mark an anomaly in the first place, it is invisible to this criterion. This is a tautology: "we detect 100% of what we detect." The "误标率 ≤5%" is better (it measures false positives against the total marked) but there is no criterion for false negatives — cases that should have been marked but were not. Without a ground-truth test corpus with known anomalies, these percentages are unverifiable.

**What must improve**: Define a test corpus with known anomalies: "提供包含 N 个已知异常的测试 JSONL 文件（含耗时 >30s 步骤 M 个、越权路径访问 K 个），检出率 = 检出数/已知异常数 ≥ 95%，误标率 = 误标数/总标记数 ≤ 5%." This grounds the criterion in an independently verifiable dataset.

### Attack 3: Alternatives — Web UI alternative analysis remains superficial

**Where**: "Web UI 方案" row: Pros "图表渲染能力强，可远程访问" / Cons "依赖浏览器，不符合终端工作流习惯" (line 81)

**Why it's weak**: After three iterations, the Web UI alternative still has one-sentence pros and cons with no substance. "图表渲染能力强" — what rendering capability does the TUI lack? "可远程访问" — this is actually a significant advantage for monitoring long-running agents remotely, which the TUI cannot do. The con "不符合终端工作流习惯" is an assumption, not evidence. Many terminal-heavy developers use browser-based tools (GitHub, Grafana, Jira) daily. The VS Code Extension alternative received detailed treatment (3 pros, 3 cons with specific estimates) — the Web UI deserves comparable depth, or it should be removed as a straw-man.

**What must improve**: Either (a) expand the Web UI row with concrete pros (D3.js/recharts charting, remote monitoring over SSH tunnel, responsive design for mobile alerts) and concrete cons (requires HTTP server process, latency for file watching over network, additional security surface for serving session data) or (b) remove it and replace with a more distinct alternative like "structured JSON export + jq" or "IDE-agnostic LSP-style protocol."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Solution lacks user-facing behavior walkthrough | ✅ Yes | Full ASCII mockup added (lines 39-57) with 3-panel layout, node content examples, status bar. Primary workflow described in 5 steps (lines 64-71). Keyboard shortcuts documented (line 72). |
| Scope contradiction on AI root cause analysis | ✅ Yes | Scope now explicitly splits: "AI 证据提取（Phase 1 / MVP）" in-scope (line 96) with clear boundary ("仅展示事实（调用链 + 耗时 + 参数），不做推断性诊断"), and "AI 根因分析（Phase 2）" in Post-MVP (lines 99-101). Risk table line 122 confirms: "Phase 1 严格限定为证据提取，不启动 agent 会话". |
| Replace Dashboard alternative with distinct option | ✅ Yes | "Agent Dashboard" alternative removed. Replaced with "Claude Code Hook 集成" (line 80) — a genuinely distinct alternative using Claude Code's built-in hook mechanism. VS Code Extension also expanded with 3 substantive pros and 3 cons. |
| Success Criteria unmeasurable | Partial | Evidence: Many criteria now have quantitative thresholds. But "100% 已标记异常点" remains circular (see Attack 2). |
| Evidence anecdotal / unsubstantiated | Partial | Evidence: Numbers are specific and concrete. But still single-author, self-reported. No community validation, user survey, or external source. |

---

## Verdict

- **Score**: 87/100
- **Target**: 80/100
- **Gap**: +7 points (target exceeded)
- **Action**: Target reached. The proposal is ready to proceed to `/write-prd`. Recommended pre-PRD improvements (non-blocking): (1) Add success criteria for dashboard rendering, replay navigation, and data redaction; (2) Replace circular anomaly detection criterion with ground-truth test corpus approach; (3) Expand or remove the Web UI alternative row.
