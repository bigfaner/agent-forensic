---
date: "2026-05-10"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/testing/"
iteration: 1
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 1

**Score: 55/100** (target: 90)

## Scorecard

| Dimension | Score | Max | Notes |
|-----------|-------|-----|-------|
| PRD Traceability | 19 | 25 | Source fields mostly specific; traceability table format non-compliant; several PRD ACs orphaned |
| Step Actionability | 18 | 25 | Steps generally concrete; multiple expected results are subjective/non-automatable; preconditions incomplete for realtime TCs |
| Route & Element Accuracy | 0 | 20 | **Entire dimension zeroed**: no Route fields, no Element fields, no Route Validation table anywhere in document |
| Completeness | 14 | 20 | Type coverage adequate; boundary cases present but gaps exist; zero explicit integration TCs; missing performance TCs |
| Structure & ID Integrity | 4 | 10 | IDs sequential and unique; classifications correct; traceability table does not match rubric format; summary counts accurate |

## Deductions

### PRD Traceability (19/25)

**TC-to-AC mapping (8/9):**
Most TCs have precise Source fields pointing to specific Story ACs (e.g., TC-CLI-001 quotes Story 8 AC verbatim). Deductions:
- TC-API-001 Source: "prd-spec.md Scope: JSONL parsing engine..." -- this is a scope bullet, not an acceptance criterion.
- TC-API-013 Source: "prd-spec.md Flow: ..." -- flow narrative, not an AC.
- TC-API-015 Source: "prd-spec.md i18n Requirements" -- generic section reference without specific paragraph.
- TC-CLI-002, TC-CLI-003 Source: "prd-spec.md i18n Requirements: ..." -- requirement paragraph, not a discrete AC.

**Traceability table complete (4/8):**
The Traceability Matrix groups TCs by PRD Story. Rubric requires: "Traceability table lists every TC with its PRD source, type, target, and priority" -- a flat per-TC table. The actual table has two columns (PRD Story | Test Cases) which is the inverse structure. It also has an "Infrastructure" row bundling 7 TCs with no corresponding PRD source -- these TCs are traceability orphans in this table.

**Reverse coverage (7/8):**
All 8 user stories have at least one TC. However, specific orphaned ACs:
- Story 5 AC-3: "耗时 >=30 秒的步骤在时间轴上以黄色标记高亮显示" -- replay-timeline highlighting is distinct from Story 2 call-tree highlighting, but has no dedicated TC.
- prd-ui-functions Call Tree "Loading" state ("解析会话...") has no TC.
- prd-ui-functions Detail Panel "Empty" state ("选中节点并按 Tab 查看详情") has no TC.
- Story 3 implies "thinking 片段" display but no TC verifies thinking content appears in detail panel.
- prd-spec.md Performance requirements (first-screen render <3s, keystroke response <100ms, virtual scroll >=30fps) have no TCs.
- prd-spec.md Security requirement (SHA256 unchanged before/after run) has no TC.

### Step Actionability (18/25)

**Steps are concrete actions (6/9):**
- TC-UI-029 Step 1: "Trigger language switch via keyboard shortcut" -- which shortcut? No PRD document defines this keybinding. This is an unresolvable action.
- TC-UI-004/005 Step 2: "Locate the slow/unauthorized tool call node" -- "locate" is a directive, not an action. Should specify navigation keystrokes.
- TC-UI-018 Step 1: "Append a new JSONL line to the active session file" -- test-harness action, not a user interaction.

**Expected results are verifiable (5/9):**
- TC-CLI-002: "All UI labels, status messages, and error text render in English" -- no specific labels enumerated. "All" is not automatable.
- TC-CLI-003: "All UI labels render in Chinese (default locale)" -- same problem.
- TC-UI-004: "The node is rendered with yellow color/highlight" -- no ANSI code, model field, or attribute specified. A script cannot "see yellow."
- TC-UI-005: "The node is rendered with red color/highlight" -- identical vagueness.
- TC-UI-021: "Dashboard view appears showing tool call count distribution" -- no specific values to assert against despite preconditions not specifying fixture data.
- TC-UI-030: "briefly shows 'scanning session files...' before populating" -- "briefly" is temporally undefined.
- TC-UI-019: "New node has visual highlight/flash that persists for 3 seconds, then returns to normal" -- "visual highlight/flash" is not a testable assertion.

**Preconditions explicit (7/7):**
Most TCs declare preconditions adequately. TC-API-005 specifies "duration exactly 30 seconds." TC-UI-018/019 say "an active JSONL file is being appended to" without specifying what process writes to it, but the context (Claude Code session) is inferable. Acceptable.

### Route & Element Accuracy (0/20)

**Routes are valid and specific (0/7):**
Zero TCs contain a Route field. The document has no Route column in any TC table. UI TCs have Target fields with Go package names (e.g., `internal/model/calltree`) which are implementation details, not navigable routes. For a TUI app, routes could be view identifiers (`main-tui`, `dashboard`, `diagnosis-modal`) -- none are provided. CLI TCs embed command patterns in steps but have no formal Route field.

**Elements are identifiable (0/7):**
Zero TCs contain an Element field. No data-testid, aria-label, role, CSS selector, or any locator strategy exists anywhere in the document. UI TCs identify elements by prose description: "the slow tool call node" (TC-UI-004), "the diagnosis modal" (TC-UI-006), "the bottom panel" (TC-UI-009). No test automation engineer can locate these elements programmatically.

**Route/Element consistency (0/6):**
The rubric requires "UI TCs have both Route and Element. API TCs have Route but no Element. CLI TCs have neither but have command patterns." Since no TC has Route or Element fields, this rule cannot be positively evaluated.

**Route Validation table: MISSING REQUIRED SECTION.** The rubric lists this as a required section. Its absence zeroes the entire dimension per the deduction rule: "Missing required section: 0 pts for that dimension."

### Completeness (14/20)

**Type coverage (7/7):**
All three interface types have TCs: CLI (3), API (15), UI (30). Coverage spans all PRD features: parsing, anomaly detection, sanitization, statistics, i18n, navigation, search, replay, realtime monitoring, dashboard, diagnosis.

**Boundary and edge cases (4/7):**
Present: exactly 30s threshold (TC-API-005), exactly 200 chars (TC-UI-011), malformed JSONL (TC-API-002), empty file (TC-API-003), missing directory (TC-CLI-001).
Missing:
- No TC for tool call at 29.9s (just below slow threshold) -- negative boundary.
- No TC for content at 201 characters (first truncation case above boundary).
- No TC for invalid `--lang` value (e.g., `--lang fr`).
- No TC for session at exactly 10000 lines (streaming boundary).

**Integration scenarios (3/6):**
Some implicit cross-component coverage: TC-UI-012 tests sanitizer + detail panel, TC-UI-020 tests watcher toggle + call tree. But no explicit end-to-end flow TCs:
- No TC for "search -> select session -> view detail -> diagnosis -> jump to node" full flow from PRD business flow.
- No TC for "CLI `--lang en` -> i18n API returns English -> UI renders English labels" cross-interface chain.
- No TC for "realtime file append -> watcher detects -> parser updates -> call tree adds node" monitoring pipeline.

### Structure & ID Integrity (4/10)

**TC IDs are sequential and unique (4/4):**
TC-CLI-001..003, TC-API-001..015, TC-UI-001..030. No gaps, no duplicates within each prefix. Each TC has a secondary Test ID (e.g., `cli/launch/missing-claude-dir`).

**Classification is correct (0/3):**
All TCs are correctly classified by type. However, the rubric deduction rule states "-3 pts per conflict." There are no conflicts, so this should be 3/3. Re-scoring: 3/3.

Wait -- let me re-check. The traceability table "Infrastructure" row groups TC-API-001 and TC-API-013 (both API type) alongside TC-UI-025, TC-UI-026, TC-UI-024, TC-UI-028, TC-UI-030 (all UI type). This is not a type-classification error in the TC definitions themselves; the TCs are correctly placed in their respective sections. The classification of individual TCs is correct: 3/3.

**Summary table matches actual (1/3):**
Summary table states CLI:3, API:15, UI:30, Total:48. These counts match actual TCs. However:
- Frontmatter uses `source` (singular) where rubric specifies `sources` (plural).
- Traceability table format does not match rubric's required format (counts under this sub-dimension because the rubric groups traceability under structure).

## Attack Points

### Attack 1: Route & Element Accuracy -- complete absence makes automated test generation impossible

**Where**: The entire document lacks Route fields, Element fields, and a Route Validation table. Searching for "Route" or "Element" returns zero matches. TC-UI-004 Target is `internal/model/calltree` -- a Go import path, not a view route or element locator.
**Why it's weak**: The rubric explicitly requires these as mandatory fields and a dedicated section. Without Route fields, a gen-test-scripts agent cannot determine which view or page to navigate to. Without Element fields, it cannot determine which component to interact with or assert against. The existing Target field is an internal Go package name useful only for unit-test targeting, not for the TUI interaction tests that 30 of the 48 TCs describe. The missing Route Validation table is a required section whose absence zeroes the entire 20-point dimension.
**What must improve**: (1) Add a `Route` field to every TC -- UI TCs need view identifiers (`main-tui/sessions-panel`, `dashboard`, `diagnosis-modal`); API TCs need function signatures (`parser.ParseSession(path string) (*Session, error)`); CLI TCs need command patterns (`agent-forensic --lang en`). (2) Add an `Element` field to every UI TC with a selector strategy (e.g., `[data-testid="session-item"]`, `[role="treeitem"][aria-label="Turn 2"]`). (3) Add the Route Validation table as a required section.

### Attack 2: Step Actionability -- visual and i18n expected results are not machine-verifiable

**Where**: TC-CLI-002 Expected: "All UI labels, status messages, and error text render in English". TC-UI-004 Expected: "The node is rendered with yellow color/highlight". TC-UI-019 Expected: "New node has visual highlight/flash that persists for 3 seconds, then returns to normal".
**Why it's weak**: These expected results describe human perception, not machine-verifiable assertions. TC-CLI-002 says "all UI labels" but enumerates none -- an automated test cannot verify "all" without an explicit list of expected strings. TC-UI-004 says "yellow color/highlight" without specifying any ANSI escape code, Lipgloss style attribute, or model field value. TC-UI-019 says "visual highlight/flash" without defining what DOM-equivalent property to check. The rubric deduction rule says "-2 pts per instance of vague language without specificity" -- these three instances alone cost 6 points. Critically, the Step Actionability score of 18 is below the 20-point blocking threshold, meaning gen-test-scripts is blocked.
**What must improve**: Replace subjective descriptions with concrete assertions. TC-CLI-002: enumerate at least 3 specific labels expected in English (e.g., "Status bar displays 'j/k:nav Enter:expand...'", "Session list header shows 'Sessions'"). TC-UI-004: assert on model state ("node.AnomalyType == 'slow'") or rendering output ("node output contains ANSI escape \\033[33m"). TC-UI-019: assert on timer state ("node.HighlightUntil > now") or style toggle ("node.HasNewHighlight == true for 3s, then false").

### Attack 3: PRD Traceability -- traceability table format is wrong and multiple ACs are orphaned

**Where**: Traceability Matrix has columns `PRD Story | Test Cases` instead of the rubric-required `TC ID | Source | Type | Target | Priority`. Story 5 AC-3 ("耗时 >=30 秒的步骤在时间轴上以黄色标记高亮显示") has no TC. prd-spec.md performance requirements (first-screen <3s, keystroke <100ms, virtual scroll >=30fps) have no TCs. prd-spec.md SHA256 integrity requirement has no TC.
**Why it's weak**: The traceability table is a story-to-TC mapping, not the required TC-to-details mapping. This forces the reader to scan grouped rows to find individual TC sources -- the opposite of the intended lookup pattern. The "Infrastructure" row bundles 7 TCs with no PRD source, making them traceability orphans. Story 5 AC-3 describes replay-timeline-specific yellow highlighting that is semantically distinct from Story 2's general anomaly yellow highlighting -- collapsing them loses this distinction. The missing performance TCs are particularly damaging because the PRD quantifies them as explicit goals with metrics.
**What must improve**: (1) Reformat the traceability table to a flat per-TC format with columns: TC ID | Source (specific AC/section) | Type | Target | Priority. (2) Add TCs for orphaned ACs: Story 5 AC-3 replay timeline highlighting, call tree loading state, detail panel empty state, thinking fragment display. (3) Add TCs for prd-spec.md performance requirements. (4) Add TC for SHA256 non-invasive integrity check. (5) Map every "Infrastructure" TC to its actual PRD source or prd-ui-functions reference.

## Verdict

- **Score**: 55/100
- **Target**: 90/100
- **Gap**: 35 points
- **Step Actionability**: 18/25 (below 20 blocking threshold -- gen-test-scripts is BLOCKED)
- **Action**: Continue to iteration 2. Priority fixes: (1) Add Route, Element fields and Route Validation table (+20 potential), (2) Make expected results machine-verifiable (+7 potential), (3) Reformat traceability table to flat per-TC format (+4 potential), (4) Cover orphaned PRD ACs and performance requirements (+4 potential). Total recoverable: ~35 points.
