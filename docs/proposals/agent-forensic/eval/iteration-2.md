---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/proposals/agent-forensic/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval -- Iteration 2

**Score: 93/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Problem Definition        │  18      │  20      │ ✅         │
│    Problem clarity           │  6/7     │          │            │
│    Evidence provided         │  7/7     │          │            │
│    Urgency justified         │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  19      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  7/7     │          │            │
│    Differentiated            │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  14      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  5/5     │          │            │
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
│ 6. Success Criteria          │  14      │  15      │ ✅         │
│    Measurable                │  5/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  93      │  100     │ ✅         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem, line 9 | Problem statement conflates two distinct concerns in one sentence: "无法直观观察...行为链路" (real-time observability) and "问题排查困难" (post-hoc debugging). These are related but separate needs; a reader could interpret the primary goal differently depending on which half they weight more heavily | -1 pts (Problem clarity) |
| Urgency, line 22 | Urgency relies solely on the proposal author's own incidents (config/production.yml deletion, rm -rf in sub-agent). No broader validation: no user survey, no community forum complaints, no issue tracker references. Compelling but single-source anecdotal evidence | -1 pts (Urgency justified) |
| Solution, alternatives rationale line 83 | "开发周期约 1.5-2 周，显著短于 VS Code extension 的 3-4 周" -- timeline speed is the only explicitly comparative differentiator. Other unique TUI capabilities (zero deps, terminal-native, lazygit UX pattern) are mentioned as properties but not argued as differentiators that alternatives fundamentally cannot match | -1 pts (Differentiated) |
| Scope, Phase 2 creep row line 131 | "Phase 2 需在 Phase 1 验证活跃用户 >=20 且证据提取功能周使用率 >=60%...后再启动" -- the usage rate metric is circular: it measures adoption of a Phase 1 feature to gate Phase 2 scope. This is a conditional trigger, not a scope boundary. A bounded scope should state what is definitively excluded, not what might be included later if metrics are met | -1 pts (Scope bounded) |
| Risk, AI evidence row line 130 | Mitigation "(2) 仅展示事实（调用链 + 耗时 + 参数），不做推断性诊断" -- this describes a design choice, not a risk mitigation. Showing selected facts can still mislead if the selection is incomplete. The risk is users drawing wrong conclusions from partial evidence; the mitigation does not address the selection/framing problem | -1 pts (Mitigations actionable) |
| Success Criteria, line 137 | "搜索结果 500ms 内返回，支持日期筛选" -- no specification of what fields are searchable (tool names? file paths? thinking content?), what matching behavior is used, or what the search scope is (current session or all). This is a testability gap: a tester cannot write test cases without knowing what search covers | -1 pts (Coverage complete) |

---

## Attack Points

### Attack 1: Solution Clarity -- Differentiation rationale is thin

**Where**: `TUI 方案选择理由：终端是 Claude Code 用户的核心工作环境；TUI 零外部依赖...开发周期约 1.5-2 周，显著短于 VS Code extension 的 3-4 周。`
**Why it's weak**: The rationale lists four points but only the timeline (1.5-2 weeks vs 3-4 weeks) is comparatively argued. "终端是核心工作环境" is asserted without evidence. "零外部依赖" is a property, not a comparative argument -- the Web UI alternative also offers zero-install ("打开 URL 即用"). The lazygit precedent is mentioned but not connected to a concrete design principle the alternatives cannot match. The differentiator should articulate what the TUI uniquely enables that alternatives fundamentally cannot.
**What must improve**: Rewrite the rationale to lead with a capability argument, not a timeline argument. Explicitly state what TUI uniquely provides (e.g., instant in-terminal context switching without breaking developer flow, seamless tmux integration, zero network surface for sensitive session data). Connect "零外部依赖" to a concrete security or deployment benefit.

### Attack 2: Success Criteria -- Search functionality is underspecified

**Where**: `搜索结果 500ms 内返回，支持日期筛选` (line 137)
**Why it's weak**: Every other success criterion is specific enough to write test cases: parsing has line-count thresholds and time bounds, anomaly detection has test corpus specs with seeded anomalies, the dashboard specifies chart types and error tolerances. Search alone lacks: (a) what content is indexed and searchable (session metadata only? tool call parameters? thinking text? file paths?), (b) what operators are supported (substring? regex? exact match?), (c) what the search scope is (current session or all sessions?), (d) what "日期筛选" means precisely (exact date? range? relative?). A tester cannot determine what to test from this criterion alone.
**What must improve**: Expand the search criterion to specify: (1) searchable fields (session date, tool name, file path, parameter substrings), (2) matching behavior (case-insensitive substring), (3) scope (all sessions by default, filterable), (4) date filter semantics (exact date match or range). Then a test case can verify: "given a session with tool call 'Write src/index.ts', searching 'index' returns it within 500ms."

### Attack 3: Risk Assessment -- AI evidence mitigation conflates design choice with risk control

**Where**: `AI 证据提取产生误导性判断 -- Medium likelihood, High impact -- (2) 仅展示事实（调用链 + 耗时 + 参数），不做推断性诊断`
**Why it's weak**: The risk is that users draw incorrect conclusions from AI-extracted evidence. Mitigation (2) says "only show facts, no inferential diagnosis." But showing facts selectively is itself a form of inference -- choosing which facts to surface determines what conclusion the user draws. If the evidence extractor shows a sequence of tool calls but omits a key intervening event (because it was not classified as "anomalous"), the user may attribute causality incorrectly. The mitigation addresses output format, not the selection/framing problem that creates the risk. Mitigation (3) ("显示免责声明") is a legal shield, not a UX risk reduction.
**What must improve**: Replace or augment mitigation (2) with a concrete safeguard against selective evidence bias: "Evidence extraction displays the full contiguous call chain (minimum 5 turns before and after each anomaly), not cherry-picked points, to preserve causal context and reduce misattribution risk." This addresses the actual risk mechanism (selective presentation) rather than asserting a design intent.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Web UI straw-man with ~20 chars of analysis) | ✅ Yes | Web UI row now has 3 substantive pros (D3/Recharts rich rendering, remote access via SSH tunnel, zero-install URL deployment) and 3 substantive cons (server-side parsing + WS hub complexity, auth/access control requirements, terminal workflow context-switch). Genuine evaluation, no longer a stub. |
| Attack 2 (Adoption risk mitigation bar -- 5 users too low for High/High) | ✅ Yes | Mitigation restructured: defines "active user" precisely (>=2 distinct launches + call tree browsing per week), targets >=20 active users with >=50% week-4 retention, adds explicit go/no-go decision (<10 users or <30% retention triggers pivot to VS Code extension or downgrade to script tool). Severity-appropriate. |
| Attack 3 (Timeline hidden in alternatives rationale, missing from Scope section) | ✅ Yes | Scope section now includes a phased timeline table (Phase 1a: Week 1 / 5 days for core data + navigation; Phase 1b: Week 2 / 5 days for analysis + value-add features) with explicit deliverables per phase and critical path identification (JSONL parser -> call tree -> anomaly marking -> AI evidence extraction). |

---

## Verdict

- **Score**: 93/100
- **Target**: 90/100
- **Gap**: -3 points (above target)
- **Action**: Target reached. No further iterations required.
