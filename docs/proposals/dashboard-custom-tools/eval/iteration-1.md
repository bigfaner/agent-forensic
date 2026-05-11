---
date: "2026-05-11"
doc_dir: "docs/proposals/dashboard-custom-tools/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 1

**Score: 74/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  12      │  20      │ ⚠️          │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   3/7    │          │            │
│    Urgency justified         │   3/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  16      │  20      │ ✅          │
│    Approach concrete         │   6/7    │          │            │
│    User-facing behavior      │   6/7    │          │            │
│    Differentiated            │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  10      │  15      │ ⚠️          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   2/5    │          │            │
│    Rationale justified       │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  14      │  15      │ ✅          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │   9      │  15      │ ❌          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   0/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  13      │  15      │ ✅          │
│    Measurable                │   4/5    │          │            │
│    Coverage complete         │   5/5    │          │            │
│    Testable                  │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  74      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem / Urgency | "参考价值下降" — vague, unquantified | -2 pts (Evidence) |
| Problem / Urgency | "skill 调用越来越多" — asserted trend with no data | -2 pts (Urgency) |
| Solution / Hook column | "等系统消息" — open-ended enumeration, hook types not fully specified | -1 pt (Approach concrete) |
| Risk #3 mitigation | "记录为已知限制，文档说明" — documentation is not a mitigation | -1 pt (Mitigations actionable) |
| Alternatives A | No pros listed for Alternative A despite "实现简单" being a real advantage | -3 pts (Pros/cons honest) |
| Risk section | Zero likelihood or impact ratings across all 4 risks | -5 pts (Likelihood + impact rated) |

---

## Attack Points

### Attack 1: Risk Assessment — likelihood and impact ratings are entirely absent

**Where**: The entire Risks section — "1. Hook 识别不准确 ... 2. Skill input 格式不稳定 ... 3. MCP 工具名格式假设 ... 4. 三列布局在窄终端下溢出"

**Why it's weak**: All four risks are listed with mitigations but zero probability or severity assessment. The rubric requires honest likelihood + impact ratings. Without them, a reader cannot prioritize which risk to address first, cannot judge whether the mitigations are proportionate, and cannot tell if the author has actually thought about probability. Risk 3 ("非标准 MCP 工具会被漏掉") could be high-likelihood in practice — the proposal treats it identically to Risk 4 (a cosmetic layout issue) with no differentiation.

**What must improve**: Add explicit likelihood (high/medium/low) and impact (high/medium/low) ratings for each risk. At minimum, explain why Risk 3 is accepted as a known limitation rather than mitigated — "document it" is not a mitigation, it is an acknowledgment of defeat.

---

### Attack 2: Alternatives Analysis — straw-man arguments, no pros for rejected alternatives

**Where**: "A. 在现有工具列表内展开（未选）：实现简单，但列表会变得很长，且三类信息没有视觉区分。" and "B. 什么都不做：仪表盘继续只展示内置工具。代价是 skill/MCP/hook 的使用情况完全不可见，复盘价值有限。"

**Why it's weak**: Alternative A's only stated pro is buried in the rejection sentence ("实现简单") and is never weighed against the cons. Alternative B lists only costs with no acknowledgment that it has zero implementation risk and zero maintenance burden — which are real advantages. The analysis reads as a post-hoc justification for a decision already made, not an honest trade-off evaluation. The chosen approach's positive case is never stated: why is a separate block with three columns better than, say, three separate sections or a tabbed view?

**What must improve**: For each alternative, list at least one genuine pro before listing cons. State the chosen approach's positive case explicitly — what does it offer that A and B do not, beyond avoiding their downsides?

---

### Attack 3: Problem Definition — evidence is a single anecdote, urgency is asserted not demonstrated

**Where**: "Skill 工具只显示总次数（如 8），看不出具体调用了哪 8 个 skill" and "随着 forge 插件体系扩展，skill 调用越来越多，缺失这部分信息让仪表盘的参考价值下降"

**Why it's weak**: The only concrete example is a parenthetical "(如 8)" — one number from presumably one session. There is no data on how frequently users encounter this gap, no user feedback, no session count showing how many sessions involve skill/MCP/hook calls. The urgency claim "越来越多" is a trend assertion with no supporting numbers. "参考价值下降" is unquantified — by how much? Compared to what baseline? A reader cannot judge whether this is a top-priority problem or a minor annoyance.

**What must improve**: Provide at least one concrete data point: how many sessions in the last N days involved skill calls? What percentage of sessions use MCP tools? Even a rough count from `grep` over the JSONL history would be more convincing than "越来越多". Replace "参考价值下降" with a specific consequence: "users cannot detect runaway hook loops" or "skill usage is invisible in post-session review."

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 74/100
- **Target**: 80/100
- **Gap**: 6 points
- **Action**: Continue to iteration 2

The proposal is structurally sound — scope and success criteria are well-defined, and the solution has a concrete ASCII mockup that makes the intent clear. The three gaps holding it below target are: (1) risk assessment missing all likelihood/impact ratings, (2) alternatives analysis that only argues against alternatives rather than for the chosen approach, and (3) a problem section that asserts urgency without any supporting data. Fixing the risk ratings alone recovers 5 points.
