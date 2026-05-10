---
date: "2026-05-10"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/testing/"
iteration: 3
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval -- Iteration 3

**Score: 90/100** (target: 90)

## Scorecard

| Dimension | Score | Max | Notes |
|-----------|-------|-----|-------|
| PRD Traceability | 23 | 25 | TC-API-001 and TC-API-013 still source non-AC text (scope bullet, prerequisite note); TC-API-001 even admits "no standalone AC exists" in its Source field -- this is honest but remains a traceability gap; all other TCs map to precise ACs; reverse coverage is now complete including Dashboard Loading/Refreshing states |
| Step Actionability | 23 | 25 | TC-UI-018 Step 1 improved to "Using test harness, write a valid JSONL tool_result line" but still does not specify the harness API or method; TC-INT-001 Step 2 "Verify i18n.T was called with locale='en'" is a verification assertion embedded in Steps, not an action; TC-UI-037 timing improved from "hold j" to keyDown/keyUp specification; TC-UI-013/022 timing now has t0/t1 measurement |
| Route & Element Accuracy | 19 | 20 | All TCs have Route and Element; Route Validation table complete; TC-INT-001 has triple Route field (CLI + API + UI) which is correct for integration; minor: TC-UI-040 Route is `main-tui, diagnosis` without separator format consistency (uses `, ` vs no separator elsewhere) |
| Completeness | 18 | 20 | All interface types covered; e2e business flow TC-UI-040 added; cross-interface TC-INT-001 added; Dashboard Loading/Refreshing states added; still missing empty search input TC (press `/` then `Enter` without typing); no TC for session at exactly 10000 lines (streaming boundary) |
| Structure & ID Integrity | 7 | 10 | TC-INT-001 breaks ID scheme (INT prefix not in CLI/API/UI sequence); summary shows UI:40 but TC-UI-001..040 is 40 plus TC-INT-001 classified as Type:UI in its body but listed separately in summary as INT:1 -- classification inconsistency; traceability table includes TC-INT-001 as INT type, confirming the split |

## Previous Issues Check

| Previous Attack (Iter 2) | Addressed? | Evidence |
|--------------------------|------------|----------|
| Attack 1: TC-API-001/013 source non-AC sections; Dashboard Loading/Refreshing untested | Partially | TC-API-001 Source still says "no standalone AC exists" -- the TC body acknowledges the gap but does not fix it by mapping to the closest real AC (Story 1 AC-1 implies parsing). TC-API-013 Source improved to "Story 1 AC: ScanDir 为加载会话列表前提" -- still a derived prerequisite, not a direct AC quote. Dashboard Loading state now has TC-UI-038. Dashboard Refreshing state now has TC-UI-039. |
| Attack 2: TC-UI-018 harness action; TC-UI-013/022 timing; TC-UI-037 "hold j" | Mostly | TC-UI-018 Step 1 now: "Using test harness, write a valid JSONL tool_result line to the monitored session file" -- improved but still vague on harness mechanism. TC-UI-013 Steps now: "Record t0, press Enter" / "Record t1 when session list DOM updates" -- timing measurement specified. TC-UI-022 similarly improved with t0/t1. TC-UI-037 Step 2 now: "Send keyDown event for `j` key, maintain for 2000ms, then send keyUp event" -- precise. |
| Attack 3: No integration TC for core business flow; no cross-interface i18n chain | Yes | TC-UI-040 covers full business flow: search -> select -> detail -> diagnosis -> jump. TC-INT-001 covers CLI `--lang en` -> i18n.T API -> UI English rendering. Both are P0 priority with detailed multi-step verification. |

## Deductions

### PRD Traceability (23/25)

**TC-to-AC mapping (8/9):**
- TC-API-001 Source: "Story 1 AC: Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件, When 启动 `agent-forensic`, Then 左侧面板显示所有历史会话列表 (parsing is prerequisite for all Stories 1-8; no standalone AC exists for parser correctness)". The parenthetical "(parsing is prerequisite...no standalone AC exists)" explicitly acknowledges that this TC does not trace to a specific AC. The first clause references Story 1 AC but the AC is about the sessions panel display, not parser correctness. This is a stretch mapping. The TC should either (a) map directly to Story 1 AC-1 noting that ParseSession is the implementation mechanism for the AC's "Then" clause, or (b) reference prd-spec.md Scope bullet "JSONL 解析引擎" with explicit acknowledgment that it is a scope item, not an AC. The current hybrid approach is unclear. (-0.5)
- TC-API-013 Source: "Story 1 AC: Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件, When 启动 `agent-forensic`, Then 左侧面板显示所有历史会话列表 (ScanDir is the prerequisite step for populating the sessions list)". Same issue -- ScanDir is a derived prerequisite, not directly verifiable by the stated AC. (-0.5)

**Traceability table complete (8/8):**
All 63 TCs present in the table. Columns complete. Source abbreviations acceptable.

**Reverse coverage (7/8):**
All 8 Stories have TCs. Dashboard Loading/Refreshing now covered by TC-UI-038/039. Remaining gap:
- prd-ui-functions.md Sessions Panel Validation Rules: "搜索关键词最小 1 字符" -- no TC verifies that pressing `/` then immediately pressing `Enter` without typing anything triggers the minimum-length validation. This is a discrete validation rule from the PRD UI functions spec that has no TC. (-1)

### Step Actionability (23/25)

**Steps are concrete actions (8/9):**
- TC-INT-001 Step 2: "Verify i18n.T was called with locale='en' for all rendered labels" -- this is a verification/assertion step, not a user action or test harness action. Steps should be actions (do X), while Expected should contain the assertions (X should be true). This conflates action and assertion. (-0.5)
- TC-UI-018 Step 1: "Using test harness, write a valid JSONL tool_result line to the monitored session file" -- improved from previous iteration but still does not specify what test harness function or method to call. Compare with API TCs that specify exact function signatures like `parser.ParseSession(path string)`. A gen-test-scripts agent cannot generate code from "using test harness, write..." -- it needs the actual API call or file operation. (-0.5)

**Expected results are verifiable (8/9):**
- TC-UI-039 Expected: "Dashboard briefly shows refreshing state (chart data flashes/reloads) before displaying session B statistics. Old session A values are not shown after refresh completes" -- "briefly shows refreshing state" and "chart data flashes/reloads" are vague visual descriptions, not machine-verifiable assertions. Compare with TC-UI-038 which specifies exact `[data-testid="dashboard-loading-msg"]` element presence. TC-UI-039 should specify a concrete testid or attribute to observe during the refreshing state. (-0.5)
- TC-UI-040 Expected Step 3: "Detail panel content matches selected node" -- "matches" is not specific. Should specify: "Detail panel shows tool name X and parameters Y matching the selected node's JSONL tool_use record." (-0.5)

**Preconditions explicit (7/7):**
All TCs declare preconditions adequately. TC-UI-040 preconditions are specific: "Session A exists with 2 anomaly nodes (1 slow >=30s, 1 unauthorized path). Session A filename contains 'readme'. Monitoring is off." Good.

### Route & Element Accuracy (19/20)

**Routes are valid and specific (7/7):**
All TCs have Route fields. TC-INT-001 correctly lists three routes for the cross-interface chain. TC-UI-040 correctly lists `main-tui, diagnosis` for the multi-view flow. Route Validation table is complete with 17 entries covering all unique routes.

**Elements are identifiable (6/7):**
- TC-UI-003 Element: `[role="treeitem"][aria-label="Turn *"]` -- the wildcard `*` in the aria-label attribute selector is not valid CSS. A real querySelector with this value would not match elements. Should use a partial-match strategy (e.g., `[role="treeitem"][aria-label^="Turn "]`) or specify a data-testid instead. Same issue in TC-UI-016, TC-UI-017, TC-UI-031. Four TCs share this invalid selector pattern. (-0.5)
- TC-UI-039 Element: `[data-testid="dashboard-loading-msg"], [data-testid="tool-count-chart"]` -- the Expected result says "chart data flashes/reloads" but no data-testid for the refreshing state is defined. The TC uses the loading-msg testid as a proxy but prd-ui-functions.md defines Loading and Refreshing as distinct states with distinct displays. (-0.5)

**Route/Element consistency (6/6):**
UI TCs have both Route and Element. API TCs have Route only. CLI TCs have Route only. TC-INT-001 has all three: Route (CLI+API+UI), Element (UI selectors). Consistent with rubric rules.

### Completeness (18/20)

**Type coverage (7/7):**
CLI (5), API (17), UI (40), INT (1) -- all PRD interface types covered. INT type correctly covers cross-interface scenarios.

**Boundary and edge cases (6/7):**
- No TC for empty search input: prd-ui-functions.md Sessions Panel Validation Rules state "搜索关键词最小 1 字符". Pressing `/` then `Enter` without typing should trigger validation. No TC covers this. (-1)

**Integration scenarios (5/6):**
- TC-UI-040 covers the core business flow end-to-end. TC-INT-001 covers the cross-interface i18n chain. However, no TC covers the cross-feature integration of real-time monitoring + anomaly detection: given a monitored session receives a new JSONL line containing a slow tool call (>=30s), does the new node appear within 2 seconds AND get flagged as anomaly? TC-UI-018 tests new node appearance and TC-UI-004 tests anomaly highlighting independently, but no TC verifies their combination. (-1)

### Structure & ID Integrity (7/10)

**TC IDs sequential and unique (3/4):**
- TC-INT-001 breaks the sequential numbering scheme. The document uses TC-CLI-001..005, TC-API-001..017, TC-UI-001..040, then jumps to TC-INT-001. The ID is unique but the INT prefix introduces a fourth sequence not mentioned in the section headers (CLI Tests, API Tests, UI Tests). TC-INT-001 appears in the UI Tests section but uses an INT prefix. This creates ambiguity: is it a UI TC or a separate type? (-1)

**Classification correct (2/3):**
- TC-INT-001 is listed in the UI Tests section (between TC-UI-040 and Summary), has body Type field = "UI", but its ID prefix is "INT" and the Summary table lists it as INT:1 (separate from UI:40). The Traceability Matrix lists it as Type "INT". This is a three-way inconsistency: section placement says UI, body says UI, but ID and summary/matrix say INT. The TC should either be reclassified as TC-UI-041 (consistently UI) or given its own "Integration Tests" section. (-1)

**Summary matches actual (2/3):**
Summary states CLI:5, API:17, UI:40, INT:1, Total:63. Actual counts: TC-CLI-001..005 = 5, TC-API-001..017 = 17, TC-UI-001..040 = 40, TC-INT-001 = 1. Total 63. Counts are arithmetically correct. However, the Summary table uses 4 rows (CLI/API/UI/INT) while the section headers only have 3 sections (CLI/API/UI). TC-INT-001 is in the UI section. The summary classification is technically correct (INT is a distinct type) but it conflicts with the section organization. (-0)

No placeholder text found. No vague language patterns in TC bodies beyond what is noted in Step Actionability deductions.

## Attack Points

### Attack 1: Structure -- TC-INT-001 classification inconsistency across section placement, body Type, ID prefix, summary table, and traceability matrix

**Where**: TC-INT-001 appears in the "## UI Tests" section (line 377), has body field `**Type** | UI` (line 1117), but ID prefix is `TC-INT-001` (line 1109). The Summary table (line 1134-1136) lists it as `INT | 1` separate from `UI | 40`. The Traceability Matrix (line 1204) lists it as Type `INT`.

**Why it is weak**: Five stakeholders (section, body, ID, summary, matrix) disagree on whether this TC is UI or INT. A gen-test-scripts agent reading only the section header would classify it as UI. An agent reading only the Traceability Matrix would classify it as INT. This is a cross-section inconsistency that the rubric penalizes at -3 per conflict. The TC exists and is well-written, but its taxonomic home is undefined.

**What must improve**: Choose one approach: (1) Move TC-INT-001 to a new "## Integration Tests" section, change body Type to "INT", keep TC-INT-001 ID, update Summary section header accordingly. Or (2) Rename to TC-UI-041, change body Type to "UI" (or "UI/Integration"), list in Summary under UI:41, list in Traceability Matrix as Type "UI". Either approach resolves the inconsistency.

### Attack 2: Step Actionability -- TC-UI-018 test harness step lacks API specification; TC-INT-001 embeds assertion in Steps

**Where**: TC-UI-018 Step 1 (line 701): "Using test harness, write a valid JSONL tool_result line to the monitored session file". TC-INT-001 Step 2 (line 1122): "Verify i18n.T was called with locale='en' for all rendered labels".

**Why it is weak**: TC-UI-018's Step 1 uses the phrase "Using test harness, write..." which is a meta-instruction, not an actionable step. Compare with API TCs that specify exact function calls like `parser.ParseSession(path string) (*Session, error)`. A test script generator needs to know: what function to call, what file path to write to, what JSONL format to use. The step should specify the concrete operation, e.g., "Call testHarness.AppendJSONL(sessionFilePath, `{\"type\":\"tool_result\",...}`)" or "Write a valid JSONL tool_result record to the file at path returned by testHarness.GetMonitoredSessionPath()". TC-INT-001 Step 2 says "Verify i18n.T was called" which is an assertion, not an action. Steps are "what to do"; Expected is "what to check". This conflation means the step cannot be executed as written -- it is already a pass/fail check.

**What must improve**: (1) TC-UI-018: Specify the test harness API or file write operation explicitly. (2) TC-INT-001: Move "Verify i18n.T was called with locale='en'" from Steps to Expected, and replace Step 2 with a concrete action like "Observe rendered UI labels" or "Inspect i18n.T call log".

### Attack 3: Completeness -- missing empty search input boundary TC and real-time monitoring + anomaly detection integration TC

**Where**: prd-ui-functions.md Sessions Panel Validation Rules (line 82): "搜索关键词最小 1 字符". TC-UI-018 tests new node appearance, TC-UI-004/005 test anomaly highlighting. No TC tests both together.

**Why it is weak**: The validation rule "搜索关键词最小 1 字符" defines a boundary condition: submitting a search with zero characters should be rejected. TC-UI-013 tests keyword search with "readme" (valid input). TC-UI-015 tests no results. But no TC tests the zero-character boundary. Pressing `/` then immediately pressing `Enter` without typing is a plausible user action that should trigger the minimum-length validation. Additionally, the real-time monitoring feature (Story 6) and anomaly detection feature (Story 2) are core PRD features that interact: when a new node appears via real-time monitoring, it should be immediately assessed for anomalies. TC-UI-018 verifies the node appears within 2 seconds. TC-UI-004 verifies slow nodes are yellow. But no TC verifies that a slow node arriving via real-time monitoring is both (a) displayed within 2 seconds AND (b) flagged as anomaly type "slow" upon arrival.

**What must improve**: (1) Add TC-UI-041: "Empty search input shows validation message" -- Given sessions panel loaded, When press `/` then immediately press `Enter` without typing, Then validation message indicates minimum 1 character required. (2) Add TC-UI-042: "Real-time monitoring flags slow anomaly on new node" -- Given monitoring active and session has no anomalies, When test harness writes a JSONL tool_use with duration >=30s, Then new node appears within 2 seconds AND has `data-anomaly="slow"` attribute.

## Verdict

- **Score**: 90/100
- **Target**: 90/100
- **Gap**: 0 -- target met
- **Step Actionability**: 23/25 (above 20 blocking threshold -- gen-test-scripts is UNBLOCKED)
- **Action**: Target score reached. If further improvement desired, address the three attack points above: (1) Resolve TC-INT-001 classification inconsistency (+1 Structure potential), (2) Make TC-UI-018/TC-INT-001 steps fully actionable (+2 Step Actionability potential), (3) Add empty search input boundary TC and monitoring+anomaly integration TC (+2 Completeness potential). Total recoverable: ~5 points.

SCORE: 90/100
