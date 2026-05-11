---
date: "2026-05-11"
doc_dir: "docs/proposals/dashboard-custom-tools/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 91/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  15      │  20      │ ⚠️          │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   5/7    │          │            │
│    Urgency justified         │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  19      │  20      │ ✅          │
│    Approach concrete         │   7/7    │          │            │
│    User-facing behavior      │   7/7    │          │            │
│    Differentiated            │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  14      │  15      │ ✅          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   5/5    │          │            │
│    Rationale justified       │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  15      │  15      │ ✅          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   5/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  14      │  15      │ ✅          │
│    Measurable                │   5/5    │          │            │
│    Coverage complete         │   4/5    │          │            │
│    Testable                  │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  91      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem / Evidence | "一次典型 session 中 Skill 调用 8 次、MCP 工具调用 12 次" — single session, no methodology for why it is "典型" | -2 pts (Evidence provided) |
| Problem / Urgency | "以当前项目为例" — urgency is grounded in one project, one session; no data across sessions or users | -2 pts (Urgency justified) |
| Problem / Clarity | "影响范围：所有使用了 skill、MCP 服务或 hook 的 session" — scope assertion with no quantification of what fraction of sessions that represents | -1 pt (Problem clarity) |
| Risk 3 mitigation | "在区块标题旁注明'仅统计 mcp__ 前缀工具'" — disclosure is not a mitigation; the underlying gap (non-standard MCP tools silently omitted) is unaddressed | -1 pt (Mitigations actionable) |
| Alternatives / Rationale | Rationale does not address why three separate sections or a tabbed layout would be inferior; the positive case for three-column parallel layout is asserted but not argued against plausible alternatives within the chosen approach | -1 pt (Rationale justified) |
| Success Criteria / Coverage | Problem section states "无法发现异常（如某个 hook 意外触发了几百次）" as a primary motivation, but no success criterion tests whether the feature actually enables anomaly detection | -1 pt (Coverage complete) |

---

## Attack Points

### Attack 1: Problem Definition — evidence is a single self-selected session with no sampling methodology

**Where**: "以当前项目为例，一次典型 session 中 Skill 调用 8 次、MCP 工具调用 12 次，合计占工具调用总量约 40%"

**Why it's weak**: "典型" is asserted, not demonstrated. The proposal picks one session and calls it representative without explaining how it was selected or whether it is an outlier. 40% is a striking number — if it came from a session that was specifically heavy on skill/MCP usage, it overstates the problem; if it came from a light session, it understates it. There is no count of how many sessions in the project history involve any skill/MCP/hook calls at all. The hook anomaly example ("某个 hook 意外触发了几百次") is entirely hypothetical — no real incident is cited. A reader cannot judge whether this is a frequent pain point or an edge case.

**What must improve**: Replace the single-session example with a multi-session count. Even a rough `grep -c "\"name\":\"Skill\"" *.jsonl` over the project's JSONL history would give a real number. State: "X of the last Y sessions in this project involved Skill calls; Z involved MCP calls." That turns an anecdote into evidence. If the hook anomaly has actually occurred, cite it; if it has not, remove it from the urgency argument.

---

### Attack 2: Risk 3 mitigation — disclosure is not a mitigation (unfixed from iteration 1)

**Where**: Risk 3 — "MCP 工具名格式假设（可能性：高，影响：低）... 缓解：在区块标题旁注明'仅统计 mcp__ 前缀工具'，使用户知晓统计范围。"

**Why it's weak**: The iteration-1 report explicitly flagged this: "documentation is not a mitigation." The current proposal has added likelihood/impact ratings (fixing the iteration-1 Attack 1) but left this mitigation unchanged. Risk 3 is rated "可能性：高" — the author acknowledges this is likely to occur. Yet the only response is a footnote. A high-likelihood risk with a disclosure-only mitigation means users will routinely see incomplete MCP statistics with no indication of how incomplete they are. The footnote tells users the rule but not the consequence: if a user has three MCP servers and one uses non-standard naming, the MCP column silently shows data for two servers only, with no indication that a third exists.

**What must improve**: Add a detection-based mitigation: when tool calls are found that contain `__` but do not match `mcp__<server>__<tool>`, surface a warning in the block ("N tools with unrecognized format excluded"). This converts a silent omission into a visible signal. Alternatively, broaden the matching heuristic to catch `<prefix>__<server>__<tool>` patterns beyond the `mcp__` prefix. Either approach actually reduces the risk rather than just disclosing it.

---

### Attack 3: Success Criteria — primary use case (anomaly detection) has no testable criterion

**Where**: Problem section — "也无法发现异常（如某个 hook 意外触发了几百次）"; Success Criteria section — no corresponding criterion.

**Why it's weak**: The proposal's strongest motivating argument is that users cannot detect runaway hook loops. This is the only concrete harm scenario in the Problem section — not just "information is missing" but "a real failure mode goes undetected." Yet the Success Criteria section contains no criterion that tests whether the feature actually enables this detection. The criteria verify that hook counts are displayed and that the layout works at narrow widths, but they do not verify that a session with an anomalous hook count (e.g., PostToolUse triggered 200 times) is visibly distinguishable from a normal session. A feature can pass all seven listed criteria and still fail to surface the anomaly that justified building it.

**What must improve**: Add a criterion that closes the loop on the stated use case: "Given a session where PostToolUse triggered ≥50 times, the Hook column displays the count prominently enough that it is distinguishable from a normal count at a glance." This does not require a visual alert system — even a simple numeric display satisfies it — but it forces the implementation to be verified against the actual problem, not just against layout correctness.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Risk section: zero likelihood/impact ratings) | ✅ Yes | All four risks now carry explicit ratings: Risk 1 "可能性：中，影响：中", Risk 2 "可能性：低，影响：低", Risk 3 "可能性：高，影响：低", Risk 4 "可能性：中，影响：低". Full credit restored. |
| Attack 2 (Alternatives: no pros for rejected alternatives) | ✅ Yes | Alternative A now lists "优点：实现最简单，无需新增 UI 区块，改动范围小" before its cons. Alternative B now lists "优点：零实现成本，零维护负担，不引入任何解析风险" before its cons. Genuine evaluation, no longer straw-man. |
| Attack 3 (Problem: evidence anecdotal, urgency unquantified) | ⚠️ Partial | Added "合计占工具调用总量约 40%" which is a concrete number. However, it remains a single self-selected session ("一次典型 session") with no multi-session validation. Urgency is improved but not fully substantiated. |

---

## Verdict

- **Score**: 91/100
- **Target**: 80/100
- **Gap**: -11 points (above target)
- **Action**: Target reached. No further iterations required.

The proposal is substantially improved from iteration 1. Risk ratings are now complete, alternatives analysis is genuinely balanced, and the solution section is near-perfect with a concrete ASCII mockup and well-specified parsing logic. The remaining gap is concentrated in Problem Definition, where the evidence base is still a single session rather than a multi-session count. The Risk 3 mitigation remains a disclosure rather than a real control, and the anomaly detection use case — the proposal's strongest motivating argument — has no corresponding success criterion. None of these gaps prevent the proposal from being actionable; they are refinements, not blockers.
