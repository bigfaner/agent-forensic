---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/prd/"
iteration: 2
target_score: "N/A"
scoring_mode: "Mode A"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 86/100** (target: N/A, mode: Mode A)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Background & Goals        │   13     │  15      │ ✅         │
│    Three elements            │   4/5    │          │            │
│    Goals quantified          │   4/4    │          │            │
│    Logical consistency       │   5/6    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Flow Diagrams             │   18     │  20      │ ✅         │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   5/7    │          │            │
│    Decision + error branches │   6/6    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3a. Functional Specs (A)     │   16     │  20      │ ⚠️         │
│    Placement & Interaction   │   5/7    │          │            │
│    Data Req & States clarity │   6/7    │          │            │
│    Validation Rules explicit │   5/6    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. User Stories              │   26     │  30      │ ⚠️         │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story (G/W/T)      │   5/6    │          │            │
│    AC verifiability          │   7/10   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Scope Clarity             │   13     │  15      │ ✅         │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/4    │          │            │
│    Consistent with specs     │   4/6    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │   86     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md: Background/Users | Only one user type defined ("使用 Claude Code 的独立开发者"), no segmentation by usage scenario or skill level | -1 pt |
| prd-spec.md: Goals table | "排查时间从 20-40 分钟降至 ≤2 分钟" is aspirational without feasibility evidence for varying session sizes | -1 pt |
| prd-spec.md: Mermaid diagram | Replay flow (n/p keys for Turn navigation) has no branch in the diagram — TreeAction decision node lacks n/p option | -1 pt |
| prd-spec.md: Mermaid diagram | Quit action (`q`) shown only from LeftPanel UserAction, not from TreeAction — pressing q from tree view should also quit | -1 pt |
| prd-ui-functions.md: UF-2 | Real-time monitoring described as "启动时自动激活" with no user control to toggle on/off — interaction flow does not describe how monitoring is disabled or if it is always-on | -1 pt |
| prd-ui-functions.md: UF-4 | Dashboard interaction step 3 says "用户切换会话" but does not specify HOW — the dashboard is a full-screen overlay; can the user access the session list from within it? | -1 pt |
| prd-ui-functions.md: UF-5 | "上下文调用链" data source listed as "该异常节点的父级路径" — vague about HOW the parent path is computed; no parsing algorithm described | -1 pt |
| prd-ui-functions.md: UF-4 | Validation "工具调用计数误差 0" and "耗时误差 ≤1 秒" state desired accuracy but do not specify a testable validation method | -1 pt |
| prd-user-stories.md: Story 4 AC3 | "Given 无匹配结果，Then 显示空状态提示" — missing explicit When clause; no user action triggers the empty state display | -1 pt |
| prd-user-stories.md: Story 2 AC1 vs Story 8 AC1 | Threshold inconsistency: Story 2 AC1 says "耗时 >30 秒" (strict greater than), Story 8 AC1 and prd-spec line 128 say "≥30s" (greater than or equal) | -1 pt |
| prd-user-stories.md: Story 7 | No AC verifies dashboard data accuracy — counts, percentages, and "最大耗时步骤" correctness are untested | -1 pt |
| prd-user-stories.md: Story 5 | Replay ACs (n/p keys) test navigation but do not verify what content is displayed at each Turn — the display behavior is untested | -1 pt |
| prd-spec.md: Scope | "越权行为检测（访问项目外文件等）" — the "等" (etc.) implies more violation types exist, but only one is specified. "项目外文件" is used across spec/ACs but never defines how "project directory" is determined (cwd? git root? config?) | -1 pt |
| prd-spec.md: Scope vs UF-5 | In-scope item "AI 证据提取（Phase 1 / MVP）" implies AI involvement, but UF-5 Diagnosis Summary is purely rule-based threshold detection — no AI component exists. Label is misleading relative to actual spec. | -1 pt |

---

## Attack Points

### Attack 1: User Stories — Threshold inconsistency and untested content accuracy

**Where**: prd-user-stories.md Story 2 AC1: "Given 会话中存在耗时 >30 秒的工具调用"; Story 8 AC1: "Given 会话中存在耗时恰好 30 秒的工具调用...该节点标黄色（≥30 秒阈值含边界值）"; Story 7 has no AC verifying dashboard data correctness
**Why it's weak**: Story 2 AC1 uses strict inequality (>30s) while Story 8 AC1 uses inclusive (≥30s). A tester following Story 2 would NOT flag a 30.0s call as anomalous, while Story 8 explicitly tests the boundary at exactly 30s. This creates conflicting pass/fail criteria for the same threshold. Additionally, Story 7 (Dashboard) has 3 ACs but none verify that the statistics displayed are actually correct — only open/close/refresh timing is tested. A dashboard showing wrong counts would pass all current ACs.
**What must improve**: Unify the threshold notation across all stories to use the same operator (recommend ≥30s consistently). Add at least one AC to Story 7 that verifies dashboard data accuracy: e.g., "Given a session with 5 Read calls and 3 Write calls, When the dashboard is displayed, Then the tool distribution shows Read:5 and Write:3."

### Attack 2: Scope — Undefined "project directory" boundary and misleading "AI 证据提取" label

**Where**: prd-spec.md Scope: "越权行为检测（访问项目外文件等）"; prd-ui-functions.md UF-2: "访问路径在项目目录外标记 unauthorized"; prd-user-stories.md Story 2 AC2: "访问项目外路径的操作"
**Why it's weak**: The "越权行为" detection relies on determining whether a path is "inside" or "outside" the project directory, but no document defines what constitutes the project directory boundary. Is it the current working directory? The git repository root? An explicitly configured path? A hardcoded assumption? This is a critical business rule that affects the core anomaly detection feature — every spec document references it, none defines it. Furthermore, the scope labels this as "AI 证据提取" but the actual implementation (UF-5) is purely rule-based (threshold checks for duration and path). There is no AI/ML component. Calling it "AI 证据提取" sets misleading expectations for implementers and stakeholders.
**What must improve**: Add an explicit definition of "project directory" to prd-spec.md (e.g., "project directory = CWD at launch time" or "project directory = git repository root detected via `git rev-parse --show-toplevel`"). Rename "AI 证据提取" to "异常证据提取" or "规则化证据提取" in the scope to match the actual non-AI implementation.

### Attack 3: Functional Specs — Real-time monitoring lacks user control and Dashboard session switching is undefined

**Where**: prd-ui-functions.md UF-2 interaction flow step 2: "启动时自动激活实时监听：检测到活跃会话（JSONL 文件持续写入）→ 调用树自动追加新节点；无活跃会话时静默等待"; UF-4 interaction flow step 3: "用户切换会话 → 仪表盘数据在 500ms 内刷新"
**Why it's weak**: The real-time monitoring is described as automatically activated at startup with no user-facing control. There is no keybinding to toggle it, no way to pause/resume monitoring, and no indication of resource cost (file watchers consume system resources). If monitoring causes performance issues, the user has no escape hatch. The UF-2 interaction flow lists 8 steps but none describe monitoring management. Meanwhile, UF-4 Dashboard says "用户切换会话" triggers a refresh, but the dashboard is a full-screen overlay — there is no visible session list to switch within. The interaction flow does not explain whether the user can access sessions from the dashboard (via a shortcut? by exiting dashboard first?) or whether this is a passive behavior triggered only from outside the dashboard view.
**What must improve**: Add a monitoring control interaction to UF-2 (e.g., "用户按 `m` 切换实时监听开/关" or document that monitoring is always-on by design with rationale). Clarify UF-4 session switching: either add a mechanism to switch sessions from within the dashboard, or explicitly state that session switching is only available from the main view and the dashboard refreshes when the user returns.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Missing Dashboard story | ✅ | Story 7 added for Dashboard View (prd-user-stories.md lines 83-93) |
| Attack 1: "top 20%" ambiguity in Story 5 AC3 | ✅ | Replaced with concrete "耗时 ≥30 秒" threshold (Story 5 AC3, line 68) |
| Attack 1: No boundary ACs (30s, 200 chars) | ✅ | Story 8 now has ACs for exactly 30s (line 104), exactly 200 chars (line 105), and large files (line 106) |
| Attack 1: No error-path ACs | ✅ | Story 8 ACs 4-6 cover corrupted JSONL, missing directory, empty file |
| Attack 2: Dashboard missing from mermaid | ✅ | Dashboard branch added: TreeAction →|s| Dashboard → DashDisplay → DashAction (prd-spec.md lines 112-117) |
| Attack 2: Real-time monitoring missing from mermaid | ✅ | RealtimeMonitor subgraph added with Watcher → NewLine → ParseNew → UpdateTree (prd-spec.md lines 137-143) |
| Attack 2: Error paths absent from diagram | ✅ | FormatCheck and SizeCheck decision diamonds added (prd-spec.md lines 90-96) |
| Attack 3: "规则引擎" vague data source | ✅ | Replaced with "阈值比较计算" (UF-2 line 119) and "阈值比较" (UF-5 line 258) |

---

## Verdict

- **Score**: 86/100
- **Target**: N/A
- **Gap**: N/A
- **Action**: Report generated for iteration 2. Key remaining issues: threshold inconsistency across stories, undefined "project directory" boundary, misleading "AI 证据提取" scope label, and gaps in interaction flow completeness for real-time monitoring and dashboard session switching.
