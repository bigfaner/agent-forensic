---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/prd/"
iteration: 3
target_score: "N/A"
scoring_mode: "Mode A"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 3

**Score: 94/100** (target: N/A, mode: Mode A)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │   15     │  15      │ ✅         │
│    Three elements            │   5/5    │          │            │
│    Goals quantified          │   4/4    │          │            │
│    Logical consistency       │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │   19     │  20      │ ✅         │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   6/7    │          │            │
│    Decision + error branches │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3a. Functional Specs (A)     │   17     │  20      │ ✅         │
│    Placement & Interaction   │   6/7    │          │            │
│    Data Req & States clarity │   6/7    │          │            │
│    Validation Rules explicit │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │   28     │  30      │ ✅         │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story (G/W/T)      │   6/6    │          │            │
│    AC verifiability          │   8/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │   15     │  15      │ ✅         │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/4    │          │            │
│    Consistent with specs     │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │   94     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md: Mermaid diagram | Quit action (`q`) only shown from UserAction in LeftPanel; no quit path visible from TreeAction — a user in TreeAction → Dashboard → DashAction cycle has no diagram-visible exit | -1 pt |
| prd-ui-functions.md: UF-5 interaction flow | No step describes what happens when a user presses `d` on a session with zero anomalies. The States table has "No Anomalies" state but the interaction flow jumps from "按 d → 弹出诊断摘要" directly to listing anomalies, skipping the empty case | -1 pt |
| prd-ui-functions.md: UF-4 Data Requirements | "各步骤耗时占比" source is "每个工具调用的耗时 / 总耗时" — ambiguous whether this is per-tool-type aggregate (Read total / total) or per-individual-call. "最大耗时步骤" does not specify tie-breaking behavior when two calls have identical max duration | -1 pt |
| prd-ui-functions.md: UF-4 Validation Rules | "工具调用计数误差 0" and "耗时误差 ≤1 秒" are accuracy assertions, not actionable validation rules — they describe what the output should be, not what inputs to validate or how to check them programmatically | -1 pt |
| prd-user-stories.md: Story 6 AC2 | "视觉标记（如闪烁或高亮边框）" uses "如" (such as) introducing ambiguity — a tester cannot verify whether "闪烁" or "高亮边框" is the correct implementation; the AC must specify one behavior or list both as acceptable alternatives | -1 pt |
| prd-user-stories.md: Story 4 AC3 | "Given 无匹配结果，Then 显示空状态提示" — missing explicit When clause; no user action triggers the empty state in the AC text, making it not independently executable by a tester | -1 pt |
| prd-user-stories.md: Story 3 AC2 | "这些值被脱敏替换为 ***" — does not define which values are sensitive within the AC itself; relies on cross-reference to prd-spec.md regex pattern, making the AC not standalone-verifiable | -0.5 pt |
| prd-user-stories.md: Story 7 AC3 | "仪表盘数据在 500ms 内刷新" — human tester cannot distinguish 500ms from 600ms without instrumentation; timing assertion should specify observable behavior or require automated verification | -0.5 pt |

---

## Attack Points

### Attack 1: User Stories — Story 4 AC3 missing When clause and Story 6 AC2 ambiguous visual spec

**Where**: prd-user-stories.md Story 4 AC3 (line 55): "Given 无匹配结果，Then 显示空状态提示"; Story 6 AC2 (line 81): "该节点有视觉标记（如闪烁或高亮边框）持续 3 秒"
**Why it's weak**: Story 4 AC3 has no When clause — a tester following Given/When/Then literally cannot execute this step. What action produces "无匹配结果"? Pressing `/`, typing a keyword, and pressing Enter? The AC must specify the triggering action. Story 6 AC2 uses "如" (such as) to list optional visual effects, but an AC must be deterministic — it should say exactly "闪烁高亮 3 秒" or "高亮边框持续 3 秒", not give a choice. Two testers reading "如闪烁或高亮边框" could expect different behaviors and both be "correct."
**What must improve**: Rewrite Story 4 AC3 as "Given 搜索关键词无匹配结果，When 按 Enter 确认搜索，Then 列表区域显示 '无匹配会话' 空状态提示". Rewrite Story 6 AC2 to specify a single deterministic behavior: "Then 该节点以高亮边框标记，持续 3 秒后恢复正常显示" (pick one, not "such as").

### Attack 2: Functional Specs — UF-4 Dashboard validation rules are accuracy goals, not validation rules; UF-5 interaction flow skips empty state

**Where**: prd-ui-functions.md UF-4 Validation Rules (line 231-232): "工具调用计数误差 0（与 JSONL 原文一致）" and "耗时误差 ≤1 秒"; UF-5 interaction flow (lines 250-254) jumps from pressing `d` to listing anomalies without handling zero anomalies
**Why it's weak**: UF-4's validation rules describe desired output accuracy ("计数误差 0") but do not specify what to validate as input or how to verify programmatically. A proper validation rule would be: "When dashboard loads, count of each tool type must equal count of tool_use messages in JSONL for the selected session, verified by re-parsing the source file." The current phrasing is a quality assertion, not a testable rule. Meanwhile, UF-5's interaction flow lists 5 steps starting from "按 d → 弹出诊断摘要" but never addresses the "No Anomalies" state that exists in the States table. A developer implementing the flow would not know what to display or what keybindings are available when there are zero anomalies.
**What must improve**: Replace UF-4 validation rules with concrete, testable validations: "1) For each tool type T, displayed count must equal `grep -c '"tool_use".*"name":"T"' session.jsonl`. 2) Total displayed duration must equal sum of (tool_result.timestamp - tool_use.timestamp) for all tool calls, within ±1 second." For UF-5, add an interaction step between current steps 1 and 2: "1b. If no anomalies detected → display '该会话未检测到异常行为' message; user presses Esc or q to close; flow ends."

### Attack 3: Flow Diagrams — No quit path visible from TreeAction/Dashboard cycle

**Where**: prd-spec.md Mermaid diagram (lines 102-118): UserAction has `q` → Quit, but TreeAction has no `q` branch; Dashboard → DashAction has no quit option either
**Why it's weak**: The Navigation Architecture (prd-ui-functions.md line 39) states "q 在主视图退出应用", and both TreeAction and Dashboard are part of the main view. But the Mermaid diagram only shows the quit path from UserAction (LeftPanel). A user navigating the tree (TreeAction) or dashboard (DashAction → Dashboard cycle) has no diagram-visible way to quit. The diagram implies the user must first navigate back to LeftPanel (via `1` key, which is documented in Navigation Architecture but absent from the diagram) to access the quit action. This is a diagram completeness issue — the main path should show all exit routes from every decision node.
**What must improve**: Add `q` → Quit branch from TreeAction decision node, and optionally from DashAction. Alternatively, add `1` → LeftPanel transition from TreeAction to show the indirect quit path. The diagram should make it clear that a user can quit from any main-view state without needing to navigate back to the sessions panel first.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Threshold inconsistency (>30s vs ≥30s) | ✅ | Story 2 AC1 now reads "耗时 ≥30 秒" (line 27), matching Story 8 AC1 "耗时恰好 30 秒...≥30 秒阈值含边界值" (line 105). Unified to ≥. |
| Attack 1: No AC for dashboard data accuracy | ✅ | Story 7 AC2 added (line 92): "Given 一个会话包含 5 次 Read 调用和 3 次 Write 调用...工具调用次数分布显示 Read:5、Write:3，与 JSONL 原文计数一致" |
| Attack 2: Undefined "project directory" boundary | ✅ | prd-spec.md scope item 7 (line 48) now has explicit definition: "启动时通过 `git rev-parse --show-toplevel` 检测 git 仓库根目录...若不在 git 仓库内，则回退到当前工作目录（cwd）。工具参数中的路径经绝对路径规范化后与项目目录前缀比较" |
| Attack 2: Misleading "AI 证据提取" label | ✅ | Renamed to "规则化证据提取（Phase 1 / MVP）" in prd-spec.md line 49; UF-5 description explicitly states "纯规则化阈值检测...不涉及 AI/ML 推理" (line 246) |
| Attack 3: Real-time monitoring lacks user control | ✅ | UF-2 interaction step 3 (line 104) now has `m` toggle: "用户按 m → 切换实时监听开/关...状态栏显示 '监听:关'/'监听:开'" |
| Attack 3: Dashboard session switching undefined | ✅ | UF-4 interaction step 3 (line 209) now specifies: "仪表盘内按 1 → 左侧弹出会话列表面板（不退出仪表盘）...选择新会话 → 仪表盘数据在 500ms 内刷新" |

---

## Verdict

- **Score**: 94/100
- **Target**: N/A
- **Gap**: N/A
- **Action**: Report generated for iteration 3. All 6 previous attack points addressed. Remaining issues are minor: one missing When clause in AC, one ambiguous visual spec in AC, Dashboard validation rules phrased as assertions rather than testable rules, UF-5 interaction flow skips empty state, and Mermaid diagram missing quit path from TreeAction. These are polish-level fixes, not structural gaps.
