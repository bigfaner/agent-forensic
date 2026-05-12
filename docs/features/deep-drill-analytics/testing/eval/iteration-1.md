---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/testing/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 1

**Score: 56/100** (target: 80)

```
┌──────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                    │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  19      │  25      │ ⚠️         │
│    TC-to-AC mapping          │   8/9    │          │            │
│    Traceability table        │   6/8    │          │            │
│    Reverse coverage          │   5/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  18      │  25      │ ⚠️         │
│    Steps concrete            │   6/9    │          │            │
│    Expected results          │   7/9    │          │            │
│    Preconditions explicit    │   5/7    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │   9      │  20      │ ❌         │
│    Routes valid              │   5/7    │          │            │
│    Elements identifiable     │   0/7    │          │            │
│    Consistency               │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  15      │  20      │ ⚠️         │
│    Type coverage             │   6/7    │          │            │
│    Boundary cases            │   6/7    │          │            │
│    Integration scenarios     │   3/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │   5      │  10      │ ⚠️         │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   1/3    │          │            │
│    Summary matches actual    │   0/3    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  56      │  100     │ ❌         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-027 Source | Source is "PRD Spec Performance" -- too vague; should reference the specific section heading or bullet | -1 pt (traceability) |
| TC-028 Source | Source is "PRD Spec Performance" -- same vagueness | -1 pt (traceability) |
| Summary table | Summary says Total: 28 but there are 34 TCs (TC-001 through TC-034). Integration TCs are not counted in the total | -2 pts (traceability table) |
| PRD Spec | Missing TC for ">50 subagents auto-degradation to summary mode" (PRD Spec Performance Requirements) | -1.5 pts (reverse coverage) |
| PRD Spec | Missing TC for ">10MB JSONL only loads index header" (PRD Spec Performance Requirements) | -1.5 pts (reverse coverage) |
| PRD Spec | Missing TC for "sensitive data sanitization: API key/token/password masking" (PRD Spec Security Requirements) | -1 pt (reverse coverage) |
| TC-001 Step 2 | "Navigate to a SubAgent node" does not specify the keystroke (j/k? arrow keys? /-search?) | -1 pt (actionability) |
| TC-013 Step 2 | "View Detail panel in Turn Overview mode" is a state observation, not a user action | -1 pt (actionability) |
| TC-028 Step 2 | "Navigate through Call Tree, Detail, Dashboard, and SubAgent overlay" is impossibly vague as a step | -1 pt (actionability) |
| TC-028 Expected | "All new features render correctly without truncation or layout issues" is not objectively verifiable | -1 pt (actionability) |
| TC-027 Expected | "Session loads without parsing SubAgent JSONL files upfront" cannot be verified by external observation | -1 pt (actionability) |
| TC-005 Pre-conditions | "Session contains a SubAgent with >50 child tool calls" -- no guidance on how to create or source this test data | -1 pt (actionability) |
| All 34 TCs Element field | Every single TC has `Element: sitemap-missing`. Zero usable selectors for test automation. | -7 pts (route & element) |
| Summary table | UI count: 28 + Integration count: 6 = 34, but Total row says 28. Arithmetic error. | -3 pts (structure) |
| TC-029 to TC-034 classification | Labeled "Integration" in section header and summary table, but typed as "UI" in each TC and traceability table. Contradictory classification. | -2 pts (structure) |
| TC-029 through TC-034 | Integration TCs duplicate existing functional TCs (TC-029 ≈ TC-001, TC-030 ≈ TC-008, TC-031 ≈ TC-013, TC-032 ≈ TC-015, TC-033 ≈ TC-017, TC-034 ≈ TC-020). No added coverage for cross-component data flow. | -3 pts (integration) |

---

## Attack Points

### Attack 1: Route & Element Accuracy — Every Element field is `sitemap-missing`, making automation impossible

**Where**: All 34 TCs have `Element: sitemap-missing`. The document header states: `WARNING: sitemap.json not found -- Element set to sitemap-missing`.

**Why it's weak**: Zero out of seven points for element identifiability. A test automation engineer cannot target any UI component because no TC specifies a selector strategy. Even for a TUI app, Bubble Tea components can be identified by their model IDs, view identifiers, or coordinate positions. The blanket `sitemap-missing` placeholder was accepted without attempting any alternative identification strategy.

**What must improve**: Replace `sitemap-missing` with TUI-appropriate selectors. Options include: Bubble Tea model identifiers (e.g., `model:call-tree`, `model:subagent-overlay`), section coordinate positions (e.g., `panel:detail`, `section:files`), or text-based locators (e.g., `text:"File Operations (top 20)"`). If the sitemap skill is not applicable to TUI, define a TUI-specific selector convention and apply it consistently.

### Attack 2: Structure & ID Integrity — Summary table arithmetic is wrong and classification is contradictory

**Where**: Summary table shows `UI: 28, Integration: 6, Total: 28`. Traceability table lists TC-029 through TC-034 as type `UI`.

**Why it's weak**: The total should be 34 (28 + 6), not 28. Additionally, TC-029-034 are categorized as "Integration Test Cases" in their section header and in the summary table, but each individual TC declares `Type: UI` and the traceability table lists them as `UI`. This cross-section inconsistency triggers the rubric's -3 pt penalty per conflict. The "Integration" category does not exist in the rubric's classification scheme (UI/API/CLI).

**What must improve**: (1) Fix the summary table total to 34. (2) Decide on a single classification: either reclassify TC-029-034 as a valid type (they are all UI tests, so they should be `UI`) or add "Integration" as a formal type and use it consistently in both the TC bodies and traceability table. (3) Merge the Integration section into the main UI section since these are all terminal UI interactions, and use a tag or label (e.g., `[Integration]`) in the TC title to distinguish them.

### Attack 3: PRD Traceability — Missing TCs for PRD performance thresholds and security requirements

**Where**: PRD Spec states: "大会话降级：>50 个子会话时自动降级为摘要模式；>10MB JSONL 只加载索引头" and "敏感数据脱敏：沿用现有 sanitizer（API key、token、password 掩码）". No TCs cover these requirements.

**Why it's weak**: Two performance thresholds from the PRD have no corresponding test cases: the 50-subagent degradation boundary and the 10MB JSONL index-only loading. These are stated requirements in the PRD Performance Requirements section. Similarly, the security requirement for data sanitization is untested. This leaves three PRD requirements with zero test coverage, which weakens reverse coverage from 8/8 to 5/8.

**What must improve**: Add TCs for: (1) Session with >50 subagent nodes triggers summary/degradation mode. (2) SubAgent JSONL file >10MB triggers index-header-only loading, full content is not parsed. (3) File paths or tool outputs containing API keys/tokens/passwords are masked by the sanitizer in Detail panel and Dashboard. Each should have a Source field referencing the specific PRD Spec section (e.g., "PRD Spec / Performance Requirements -- >50 subagents").

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| N/A (iteration 1) | -- | -- |

---

## Verdict

- **Score**: 56/100
- **Target**: 80/100
- **Gap**: 24 points
- **Step Actionability**: 18/25 (above blocking threshold of 20 is NOT met -- 18 < 20)
- **Action**: Continue to iteration 2. Priority fixes: (1) Replace all `sitemap-missing` element fields with TUI-appropriate selectors (+7 pts potential). (2) Fix summary table arithmetic and resolve Integration/UI classification contradiction (+5 pts potential). (3) Add missing TCs for PRD performance thresholds and security requirements (+3 pts potential). (4) Improve step concreteness -- specify keystrokes, provide test data fixtures, make expected results machine-verifiable (+4 pts potential).
