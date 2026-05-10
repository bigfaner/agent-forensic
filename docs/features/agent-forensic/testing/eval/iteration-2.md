---
date: "2026-05-10"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/testing/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 84/100** (target: 90)

## Scorecard

| Dimension | Score | Max | Notes |
|-----------|-------|-----|-------|
| PRD Traceability | 21 | 25 | Traceability table now flat per-TC format; several AC source fields still imprecise; TC-API-013 sources a flow narrative not an AC; 2 PRD ACs still lack dedicated TCs |
| Step Actionability | 21 | 25 | Most expected results now concrete and machine-verifiable; TC-UI-018 Step 1 still uses test-harness action; TC-UI-019 preconditions lack process specification; TC-UI-037 "hold j key for 2 seconds" is ambiguous as automation instruction |
| Route & Element Accuracy | 17 | 20 | Route and Element fields added to all TCs; Route Validation table present; API Routes use function signatures (good); some Route values are generic (e.g., `main-tui`); CLI Routes are command patterns (correct); minor consistency issues |
| Completeness | 16 | 20 | All interface types covered; new boundary TCs (TC-API-016, TC-API-017, TC-CLI-005); performance TCs added (TC-UI-035/036/037); SHA256 TC added; still missing explicit integration flow TCs |
| Structure & ID Integrity | 9 | 10 | IDs sequential and unique; classification correct; summary counts match (5+17+37=59); frontmatter uses `sources` (plural, correct); traceability table reformatted to flat per-TC |

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Route & Element Accuracy — complete absence (0/20) | Partially | Route and Element fields now present on all TCs. Route Validation table added as required section. API Routes use function signatures like `parser.ParseSession(path string) (*Session, error)`. UI Routes use view identifiers like `main-tui/sessions-panel`. Elements use `data-testid` and `aria-label` selectors. However, several UI Routes are the same generic `main-tui` without panel specificity, and API Routes are function signatures (acceptable but unconventional). |
| Attack 2: Step Actionability — visual/i18n expected results not machine-verifiable | Mostly | TC-UI-004 Expected now: "node.AnomalyType == 'slow'. Rendered output contains ANSI escape sequence `\033[33m`". TC-UI-005 similarly specifies ANSI code `\033[31m`. TC-UI-019 now specifies: "node.HighlightUntil == t0 + 3s" and "node.HasNewHighlight == true". TC-CLI-002 now enumerates specific strings: "Status bar contains 'j/k:nav Enter:expand...'". These are now machine-verifiable. Remaining issues: TC-UI-018 Step 1 "Append a new JSONL line" is a harness action; TC-UI-037 "hold j key for 2 seconds" is ambiguous. |
| Attack 3: PRD Traceability — traceability table wrong format, orphaned ACs | Mostly | Traceability table reformatted to flat per-TC with columns: TC ID | Source | Type | Target | Priority. New TCs added for previously orphaned ACs: TC-UI-031 (Story 5 AC-3 replay timeline yellow), TC-UI-032 (Call Tree Loading state), TC-UI-033 (Detail Panel Empty state), TC-UI-034 (thinking fragment), TC-UI-035/036/037 (performance), TC-CLI-004 (SHA256). However TC-API-001 Source still says "prd-spec.md Scope: JSONL 解析引擎..." which is a scope bullet, not an AC. TC-API-013 sources "prd-spec.md Flow: 扫描目录" which is a flow narrative. |

## Deductions

### PRD Traceability (21/25)

**TC-to-AC mapping (7/9):**
Most TCs now have precise Source fields. Remaining issues:
- TC-API-001 Source: "prd-spec.md Scope: JSONL 解析引擎：解析 `~/.claude/` 下的 Claude Code 会话 JSONL 文件，提取结构化数据" -- this is a scope checklist bullet, not an acceptance criterion. (-1)
- TC-API-013 Source: "prd-spec.md Flow: 用户启动 `agent-forensic`，工具扫描 `~/.claude/` 目录查找 JSONL 会话文件" -- this is a flow description narrative step, not an AC. (-1)
- TC-API-015 Source: "prd-spec.md i18n Requirements" -- section reference without specific paragraph or quote. Vague compared to the specific Story AC references used elsewhere. (-0.5, rounding to 0)

**Traceability table complete (7/8):**
The flat per-TC table is now present with all 59 TCs listed. All columns (TC ID, Source, Type, Target, Priority) are filled. No TC is missing from the table. Deduction:
- The table abbreviates some source references compared to the TC body. E.g., TC-CLI-001 table entry says "Story 8 AC: ~/.claude/ 目录不存在" while the TC body has the full AC text. This is acceptable but means a reader cannot reconstruct the full AC from the table alone. (-1)

**Reverse coverage (7/8):**
All 8 user stories have at least one TC. New TCs added for previously orphaned items: TC-UI-031 (Story 5 AC-3), TC-UI-032 (loading state), TC-UI-033 (empty state), TC-UI-034 (thinking), TC-UI-035/036/037 (performance), TC-CLI-004 (SHA256). Remaining gaps:
- Story 7 AC: Dashboard Loading state "计算统计数据..." (from prd-ui-functions.md Dashboard States) has no dedicated TC. TC-UI-021 only verifies populated dashboard, not the loading-to-populated transition. (-1)
- prd-ui-functions.md Dashboard "Refreshing" state ("数据闪烁刷新") when switching sessions has no dedicated TC. TC-UI-022 verifies the data refreshes but does not assert on the refreshing visual state. (-0.5)

### Step Actionability (21/25)

**Steps are concrete actions (7/9):**
Most steps are now specific keystroke-level actions. Issues:
- TC-UI-018 Step 1: "Append a new JSONL line to the active session file" -- this describes a test harness setup action, not a user interaction with the TUI. The step does not specify what tool or method performs the append. A gen-test-scripts agent cannot generate a Playwright action for this. (-1)
- TC-UI-037 Step 2: "Hold `j` key for 2 seconds to trigger rapid scrolling" -- "hold for 2 seconds" is ambiguous for automation. Should specify: "Press and hold `j` key (keyDown event) for 2000ms then release (keyUp event)" or "Send 30 consecutive `j` keypresses within 2 seconds". (-1)

**Expected results are verifiable (7/9):**
Significant improvement. TC-UI-004/005 now specify ANSI codes and model fields. TC-UI-019 specifies `node.HighlightUntil` and `HasNewHighlight`. TC-CLI-002 enumerates specific strings. Remaining issues:
- TC-UI-013 Expected: "Results appear within 500ms" -- this is a timing assertion with no specification of HOW to measure. Should specify: "t(results_rendered) - t('readme' submitted) < 500ms" or similar. (-1)
- TC-UI-022 Expected: "Dashboard data refreshes to show session B statistics within 500ms" -- same timing measurement ambiguity. No specification of what observable change marks the start/end of the 500ms window. (-1)
- TC-UI-035 Expected: "t1 - t0 < 3000ms" -- here the document DOES specify timestamps (Step 1 records t0, Step 2 records t1). This inconsistency between TC-UI-035 (which defines timing measurement) and TC-UI-013/022 (which do not) suggests the author knows how to write timing assertions but did not apply it uniformly.

**Preconditions explicit (7/7):**
Most TCs declare preconditions adequately. TC-API-005 specifies "duration exactly 30 seconds." TC-UI-019 now specifies "node.HighlightUntil" as the check. TC-UI-037 specifies "A JSONL session file with >10000 lines." Acceptable.

### Route & Element Accuracy (17/20)

**Routes are valid and specific (6/7):**
Routes are now present on all TCs. UI Routes use view identifiers, API Routes use function signatures, CLI Routes use command patterns. Issues:
- TC-UI-025 Route: `main-tui` -- too generic. This TC tests Tab cycling across all three panels. The Route should be `main-tui` with additional note about cycling, but it matches no single panel. Marginally acceptable. (-0.5)
- TC-UI-029 Route: `main-tui` -- same generic route. TC-UI-029 tests language switching which affects the entire app. The generic route is semantically correct but less useful for navigation targeting. (-0.5)

**Elements are identifiable (6/7):**
All UI TCs now have Element fields using `data-testid` or `aria-label` selectors. This is a major improvement. Issues:
- TC-UI-031 Element: `[role="treeitem"][aria-label="Turn *"]` -- the wildcard `*` in aria-label is not a valid CSS selector pattern. Should use a partial match strategy like `[role="treeitem"]` with a note about Turn-level selection, or specify the exact aria-label format. (-0.5)
- TC-UI-016 and TC-UI-017 use the same Element: `[role="treeitem"][aria-label="Turn *"]` -- same wildcard issue. (-0.5)
- TC-UI-026 has two Routes (`main-tui`, `diagnosis`) which is correct for a cross-view TC, but the Element field is `[data-testid="diagnosis-modal"]` which only covers the overlay, not the main view quit behavior. (-0.5)

**Route/Element consistency (5/6):**
UI TCs have both Route and Element. API TCs have Route but no Element. CLI TCs have Route (command pattern) but no Element. This matches the rubric rule. Issues:
- TC-CLI-004 Route: `agent-forensic` -- but the TC involves recording SHA256 before and after. The Route field should arguably show the full command sequence including the quit action, but the current format is acceptable as a launch route. (-0)

### Completeness (16/20)

**Type coverage (7/7):**
All three interface types have TCs: CLI (5), API (17), UI (37). Coverage spans all PRD features including new performance and security TCs.

**Boundary and edge cases (5/7):**
Significant improvement. New boundary TCs:
- TC-API-016: 29.9s (below slow threshold) -- addresses previous gap
- TC-API-017: 201 characters (above truncation boundary) -- addresses previous gap
- TC-CLI-005: invalid `--lang fr` -- addresses previous gap

Remaining gaps:
- No TC for session at exactly 10000 lines (streaming boundary: is it streamed or full-parsed?). TC-API-004 tests >10000 lines but not the boundary value. (-1)
- No TC for content at exactly 200 characters + 1 byte (201 tested, but 200 boundary already covered by TC-UI-011). Actually TC-API-017 covers 201 and TC-UI-011 covers exactly 200. The gap at 200+1 is covered. Withdrawn.
- No TC for empty search input (search keyword minimum 1 character rule from prd-ui-functions.md). Pressing `/` then `Enter` without typing anything. (-1)

**Integration scenarios (4/6):**
Some improvement. TC-UI-012 tests sanitizer + detail panel integration. TC-UI-020 tests watcher toggle + call tree. TC-UI-025 tests cross-panel Tab cycling. However, still no explicit end-to-end flow TCs:
- No TC for the full business flow from prd-spec.md: "search -> select session -> view detail -> diagnosis -> jump to node". This is the core user journey and should have at least one integration TC. (-1)
- No TC for cross-interface chain: "CLI `--lang en` -> i18n API returns English -> UI renders English labels". TC-CLI-002 tests the CLI flag and TC-API-014 tests the API lookup independently but no TC verifies the full chain. (-1)

### Structure & ID Integrity (9/10)

**TC IDs are sequential and unique (4/4):**
TC-CLI-001..005, TC-API-001..017, TC-UI-001..037. No gaps, no duplicates. Each TC also has a secondary Test ID (e.g., `cli/launch/missing-claude-dir`).

**Classification is correct (3/3):**
All TCs are correctly classified: CLI TCs in CLI section, API TCs in API section, UI TCs in UI section. No cross-section misplacements.

**Summary table matches actual (2/3):**
Summary table states CLI:5, API:17, UI:37, Total:59. Actual counts: CLI TC-CLI-001..005 = 5, API TC-API-001..017 = 17, UI TC-UI-001..037 = 37. Total 59. Counts match.
- However, TC-UI-035/036/037 are performance TCs classified as UI type. They test UI rendering performance, so the classification is defensible, but they arguably cross into performance/non-functional testing which is a separate category in many test frameworks. The rubric does not require a "Performance" type, so this is acceptable. (-0)
- Traceability table includes all 59 entries, matching the summary. (-0)
- Deduction: The frontmatter is missing the `status` field that some rubric implementations expect. The document has `feature`, `generated`, `sources` -- all required fields are present. No deduction here.
- Minor: The Traceability Matrix table header row (line 1062-1063) is followed by 59 data rows, all correctly formatted. (-0)

Final -1 deduction: TC-API-004 Route field says `parser.ParseIncremental(path string, batchSize int) (<-chan PartialResult, error)` but the TC step says "requesting first batch" without specifying batchSize value. The Route signature includes the parameter but the step does not provide a concrete value. Minor inconsistency between Route specification and step actionability.

## Attack Points

### Attack 1: PRD Traceability — TC-API-001 and TC-API-013 source non-AC sections, and Dashboard Loading/Refreshing states remain untested

**Where**: TC-API-001 Source: "prd-spec.md Scope: JSONL 解析引擎：解析 `~/.claude/` 下的 Claude Code 会话 JSONL 文件，提取结构化数据". TC-API-013 Source: "prd-spec.md Flow: 用户启动 `agent-forensic`，工具扫描 `~/.claude/` 目录查找 JSONL 会话文件". Dashboard Loading state ("计算统计数据...") and Refreshing state from prd-ui-functions.md have no TC.

**Why it's weak**: Scope bullets and flow narratives are not acceptance criteria -- they are design context. A scope bullet cannot be verified as true/false. The PRD has specific acceptance criteria in Stories and in prd-ui-functions.md Validation Rules that should be the Source. For TC-API-001, the actual traceable source would be the JSONL parsing engine's functional requirements (which are implied across multiple Story ACs but not explicitly stated as a standalone AC -- this itself is a PRD gap, but the TC should acknowledge this). TC-API-013 should reference Story 1 AC (which implicitly requires scanning) or prd-spec.md Flow step 1 with explicit acknowledgment that it is a flow step. Additionally, prd-ui-functions.md Dashboard States table defines Loading and Refreshing states that have no TCs -- the dashboard appears in 4 Story 7 ACs but only the Populated state is tested.

**What must improve**: (1) Map TC-API-001 to the closest actual AC (e.g., "prd-spec.md Scope In Scope: JSONL 解析引擎 (implied by Stories 1-8 all requiring parsed data)"). (2) Map TC-API-013 to Story 1 AC-1 which requires sessions to be loaded (scanning is a prerequisite). (3) Add TC for Dashboard Loading state: press `s`, observe "计算统计数据..." message, then observe populated dashboard. (4) Add TC for Dashboard Refreshing state visual indicator.

### Attack 2: Step Actionability — TC-UI-018 harness action and TC-UI-013/022 timing assertions lack measurement specification

**Where**: TC-UI-018 Step 1: "Append a new JSONL line to the active session file". TC-UI-013 Expected: "Results appear within 500ms". TC-UI-022 Expected: "Dashboard data refreshes to show session B statistics within 500ms". TC-UI-037 Step 2: "Hold `j` key for 2 seconds to trigger rapid scrolling".

**Why it's weak**: TC-UI-018's first step is not a user action -- it is a test infrastructure operation. A gen-test-scripts agent generating Playwright code cannot produce a "append JSONL line" step. The step should specify the harness mechanism: e.g., "Using test harness, append a valid JSONL tool_use record to the session file at [path]". TC-UI-013 and TC-UI-022 assert timing constraints ("within 500ms") without defining the measurement endpoints. Compare with TC-UI-035 which correctly defines t0 and t1 with specific capture points. The inconsistency suggests these TCs cannot be automated consistently. TC-UI-037's "hold j key for 2 seconds" does not specify whether this is a continuous keyDown event or repeated keypresses, which produces very different scroll behavior.

**What must improve**: (1) TC-UI-018: Replace Step 1 with a test-harness-specific instruction: "Using test harness, write a valid JSONL tool_result line to the monitored session file". (2) TC-UI-013: Add timing measurement steps: "Record t0 when Enter is pressed after typing 'readme'. Record t1 when session list DOM updates with filtered results. Assert t1 - t0 < 500ms." (3) TC-UI-022: Add similar timing measurement. (4) TC-UI-037: Change "Hold `j` key" to "Send keyDown event for `j` key, maintain for 2000ms, then send keyUp event".

### Attack 3: Completeness — no integration TC for the core business flow (search -> select -> detail -> diagnosis -> jump)

**Where**: prd-spec.md Business Flow Description defines a 9-step user journey: "scan sessions -> browse -> select -> search -> browse call tree -> anomaly visible -> Tab detail -> d diagnosis -> Enter jump to node". No single TC or TC chain covers this end-to-end flow.

**Why it's weak**: The PRD's Business Flow Description is the primary definition of how users interact with the tool. Individual TCs verify each step in isolation (TC-UI-013 for search, TC-UI-002 for select, TC-UI-009 for detail, TC-UI-006 for diagnosis, TC-UI-007 for jump). But no TC verifies the sequential flow works: that selecting a session after searching still loads the correct tree, that pressing Tab after diagnosis closes the modal and moves to detail, that jumping from evidence correctly repositions after a session switch. These cross-step state transitions are where integration bugs most commonly appear. The previous iteration report called this out explicitly and it remains unaddressed. Similarly, no TC verifies the CLI -> API -> UI cross-interface chain for i18n: TC-CLI-002 tests English launch, TC-API-014 tests lookup, TC-UI-029 tests keyboard switch, but no TC proves `--lang en` reaches the i18n API and the result reaches the UI.

**What must improve**: (1) Add TC-UI-038: End-to-end business flow -- "Given sessions exist with anomalies, When user searches for a keyword, selects a session, browses call tree, presses Tab for detail, presses d for diagnosis, presses Enter to jump, Then each transition works correctly and state is consistent across steps." (2) Add TC-INT-001: Cross-interface i18n chain -- "Given `--lang en` flag is set, When app launches, Then i18n.T is called with locale 'en' and all rendered UI labels match the English locale strings."

## Verdict

- **Score**: 84/100
- **Target**: 90/100
- **Gap**: 6 points
- **Step Actionability**: 21/25 (above 20 blocking threshold -- gen-test-scripts is UNBLOCKED)
- **Action**: Continue to iteration 3. Priority fixes: (1) Add integration TCs for core business flow and i18n cross-interface chain (+4 Completeness potential), (2) Make timing assertions measurable in TC-UI-013/022 (+2 Step Actionability potential), (3) Fix TC-API-001/013 source fields to reference actual ACs (+2 Traceability potential). Total recoverable: ~8 points, only 6 needed to reach target.
