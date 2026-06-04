# Proposal Evaluation Report — Iteration 3 (Final)

**Proposal**: Session Experience Enhancement
**Date**: 2026-06-04
**Evaluator**: Adversary (CTO persona)
**Previous Iteration**: Iteration 2 (Score: 836/1000)

---

## Iteration-2 Issue Tracking

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | Comparison table uses only self-invented alternatives | **Addressed** | New row: "外部工具组合（lnav + jq）" with honest pros/cons and rejection rationale about closed-loop workflow |
| 2 | Tool descriptions thin — one sentence each | **Addressed** | Each tool now gets a full paragraph: what is adopted, what is rejected, and why. lnav (pattern adoption + standalone rejection), jq (simplicity argument), VisiData (offset inspiration + architecture mismatch) |
| 3 | Cross-domain analogies from adjacent domains only | **Addressed** | Third analogy added: "数字取证时间线（Autopsy/Sleuth Kit）" — Turn→ToolUse→Result hierarchy inspired by Autopsy's browsable timeline |
| 4 | No calendar anchor | **Partially Addressed** | "16-24 总工时，含 20% 诊断缓冲" — total hours with buffer, but still no calendar projection (e.g., "2-3 weeks") or phase-level timeboxes |
| 5 | Phase ordering vs P0-P3 mismatch | **Addressed** | Phase 1 now has explicit rationale: "按键 bug 阻塞所有后续测试和验证操作——如果按键不响应，无法在 TUI 内导航到特定会话或 Turn 来验证 Phase 2/3 的数据修复效果" |
| 6 | Missing watcher race condition | **Addressed** | New risk row: "Watcher 会话切换期间的竞态条件 | M | M" with detailed mitigation including 2s polling fallback and immediate ParseIncremental on switch |
| 7 | TaskOutput SC uses subjective "可读展示" | **Addressed** | SC rewritten: "解析成功的 JSON/XML 内容按缩进和换行格式化显示；解析失败的原始内容按终端宽度自动换行显示" — concrete and testable |

**Resolution rate**: 6 fully addressed, 1 partially addressed, 0 unaddressed.

---

## Phase 1 — Reasoning Audit

### Problem → Solution Trace

All 8 defects map to in-scope items. The mapping is complete and unchanged across all iterations. No orphan problems or phantom solutions.

### Solution → Evidence Trace

The urgency section provides P0-P3 ranking with per-item time-cost estimates. The Phase 1 rationale now explicitly links the solution ordering to the evidence chain: key bug must be fixed first because it blocks verification of all subsequent work. This closes a reasoning gap identified in iteration 2.

The evidence remains author-estimated rather than user-observed — no telemetry, surveys, or support tickets are cited. This is a structural limitation of the proposal that was not addressed, though it was flagged in both previous iterations.

### Evidence → Success Criteria Trace

Each in-scope item has a corresponding SC. The iteration-3 revisions resolved the final SC tension:

- **TaskOutput SC** (Cluster F from iteration 2): Now specifies two concrete display behaviors — formatted for parseable content, auto-wrapped for raw content. **Satisfiable and measurable**.

### Self-Contradiction Check — SC Consistency Deep-Dive

**Cluster A: Watcher / Refresh**
- SC: "编辑 .jsonl 文件后，写入暂停 500ms 内或累计最多 2s 后，TUI 自动刷新"
- NFR: "收到 WatcherEventMsg 启动 500ms tick，同文件后续事件重置；设 2s 最大延迟上限"
- Risk: "切换会话时 remove-watch 与 add-watch 之间存在窗口...切换期间以 2s 间隔轮询目标目录直到 watch 建立"
- **Satisfiable**: The 2s polling during watch transition aligns with the 2s max delay in the SC. The risk mitigation is compatible with the debounce design.

**Cluster B: Session Discovery**
- SC: "会话列表项目数 ≥ find ... | wc -l；若诊断确认权限/符号链接问题，则可访问目录计数匹配即可并记录 warning"
- In Scope #4: "先诊断（对比 ScanProjectsDir 输出 vs find 结果），确认缺失会话和根因后修复"
- **Satisfiable**: Conditional SC provides clear fallback path.

**Cluster C: Key Bug**
- SC: "按键 bug 根因确认并修复（非输入模式下所有按键正常响应，覆盖 80+ i18n 键）"
- In Scope #1: "验证根因（key normalization vs handler routing）后修复"
- **Satisfiable**: No change needed.

**Cluster D: sessions-index.json**
- SC: "sessions-index.json 存在时标题用 summary 字段；不存在时回退到首条用户消息"
- Constraint: "sessions-index.json 仅约 10% 项目有，fallback 是常态"
- In Scope #3: "ScanProjectsDir 遍历时附带检查 sessions-index.json 构建 sessionId→summary 映射，fallback 到当前行为（~90% 需 fallback）"
- **Satisfiable**: Fallback is explicit and quantified.

**Cluster E: CLI --session UUID**
- SC: "`--session <valid-uuid>` 启动后直接展示目标会话，无需翻页"
- Edge case: "UUID 在多个项目中存在 → 搜索所有项目，取最新匹配（以 JSONL 文件 mtime 为准）"
- In Scope #5: "filepath.WalkDir 文件名前缀匹配搜索 UUID 并直接打开"
- **Satisfiable**: WalkDir supports multi-project search; mtime criterion is defined.

**Cluster F: TaskOutput Display (Revised)**
- SC: "TaskOutput 调用结果在详情面板中展示：解析成功的 JSON/XML 内容按缩进和换行格式化显示；解析失败的原始内容按终端宽度自动换行显示"
- In Scope #6: "解析 TaskOutput 工具调用内容并展示"
- Risk: "TaskOutput 格式多样解析不完整 | M | L | 先覆盖常见格式，异常显示原始内容"
- **Satisfiable**: The SC now defines two concrete display modes. The risk mitigation aligns — "异常显示原始内容" maps to "解析失败的原始内容按终端宽度自动换行显示."

**No new contradictions detected.** All SC clusters are internally satisfiable.

---

## Phase 2 — Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

**Problem stated clearly (37/40)**: Eight specific defects with item references and P0-P3 prioritization. The three-category framing ("数据完整性、交互细节和内容展示") provides a clear mental model. The Phase 1 rationale now links urgency to execution order, reinforcing the problem framing. Deduction: Minor — the three categories do not perfectly align with the solution phases, creating potential confusion about whether "交互细节" problems are more or less important than "数据完整性" problems.

**Evidence provided (35/40)**: Item references provide traceability. P0-P3 ranking with per-item impact quantification ("手动翻页浪费 30-60s/次", "约 40% Turn 详情需外部工具", "按键每分钟多次触发") provides concrete estimates. Deduction: Evidence is author-estimated, not user-observed. The "40%" figure for Turn details needing external tools has no source. No usage telemetry, user complaints, or support tickets are cited across three iterations. This is a structural limitation unlikely to change further.

**Urgency justified (27/30)**: The Phase 1 rationale strengthens the urgency argument by establishing a dependency chain: key bug blocks all verification. "延迟修复 P0/P1 意味核心分析流程持续不可靠" is a valid cost-of-delay statement. Deduction: No data on user base size or incident frequency. The urgency is reasonable but not externally validated.

**Score: 99/110**

### 2. Solution Clarity (120 pts)

**Approach is concrete (39/40)**: Three-phase breakdown with specific technical approaches per item. The Phase 1 rationale ("按键 bug 阻塞所有后续测试和验证操作") resolves the ordering question from iteration 2. Each item names its implementation strategy. A reader can explain back what will be built. Deduction: Minor — the key bug item says "验证根因（key normalization vs handler routing）" identifying two hypotheses but does not discuss the fix path for each. If root cause is "handler routing," what changes?

**User-facing behavior described (42/45)**: Key Scenarios section covers 9 scenarios. The TaskOutput SC now specifies concrete display behavior. The watcher SC specifies measurable thresholds. Deduction: Item 25 (session discovery) — the user sees more sessions, but is there a count indicator or status message confirming completeness? Item 18 display — the SC now covers it well, but the Key Scenarios section still says only "展示解析后的任务输出" without describing the dual display modes.

**Technical direction clear (34/35)**: Specific technical choices named throughout. The watcher constraint now includes the switch strategy ("仅监控当前会话目录，切换时更新 watch target") and the race condition mitigation ("切换期间以 2s 间隔轮询目标目录直到 watch 建立"). The offset prerequisite for WatcherEventMsg is flagged as a risk item. Deduction: The "watch target" switching mechanism's implementation detail (remove-watch then add-watch vs. atomic swap) is not specified. fsnotify may not support atomic watch replacement.

**Score: 115/120**

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (35/40)**: Significant improvement over iteration 2. Three real tools cited with expanded descriptions:
- **lnav**: Pattern adoption (real-time tracking + structured display) and rejection (standalone mode, context-switching cost). The rejection rationale is specific: "取证分析需在 TUI 内完成会话导航与内容查看的闭环，频繁切换外部工具会打断分析心流."
- **jq**: UUID search framed as simplified jq filter; rejection based on simplicity (WalkDir sufficient) and avoidance of external process dependency plus learning curve.
- **VisiData**: Offset-based loading inspiration; rejection based on architecture mismatch (Bubble Tea vs. VisiData's table model, JSONL tree structure vs. tabular data).

Deduction: No version numbers, no links, no benchmarks. "lnav" is described functionally but no specific version or feature set is referenced. The "独立应用" characterization of lnav is accurate but shallow — lnav does support piped input and could theoretically be integrated as a subprocess.

**At least 3 meaningful alternatives (24/30)**: The comparison table now includes "外部工具组合（lnav + jq）" as a genuinely different approach. This is no longer just scope-breadth comparison. The verdict includes a specific rationale: "取证分析需要会话列表浏览→Turn 详情→工具结果的闭环工作流，外部工具组合无法提供." Deduction: "仅修 bug" still reads as a straw man — positioned as "其余 6 项被推迟" with no argument for why phased delivery might be valid (e.g., ship P0/P1 now, defer P2/P3).

**Honest trade-off comparison (21/25)**: The external tools alternative now has real pros ("无需修改代码；lnav 已有成熟过滤和格式化能力") and cons ("无法集成会话导航；需手动拼路径、切换工具"). This is a meaningful improvement. Deduction: The "改动量较大" assessment for the full enhancement is still qualitative. No estimation of lines of code, files touched, or integration complexity.

**Chosen approach justified against benchmarks (20/25)**: The justification is now multi-part: (1) each item is independent, (2) the closed-loop workflow cannot be achieved with external tools, (3) the phase structure contains risk. This is a more complete argument than iteration 2's "各项独立，风险可控." Deduction: The assertion that external tools "cannot provide" the workflow is strong but untested. Has the author tried using lnav + jq for the forensic workflow? What specific friction points were encountered?

**Score: 100/120**

### 4. Requirements Completeness (110 pts)

**Scenario coverage (38/40)**: 9 scenarios covering happy paths, edge cases, and errors. The watcher race condition is now addressed in the risk table with a mitigation (2s polling during transition). Deduction: Two scenarios still missing: (1) corrupted JSONL files with partial/malformed lines during incremental parsing — what does ParseIncremental do with a half-written line? (2) concurrent write + read during active forensic session — the watcher fires while the user is reading; does the UI scroll or hold position?

**Non-functional requirements (35/40)**: Four NFRs with concrete numbers. The debounce NFR is specific and implementable. Deduction: Still missing: memory usage impact of sessions-index.json loading for 100+ projects; startup time impact of enhanced ScanProjectsDir traversal; performance of --session UUID search across large project counts. These were flagged in iteration 2 and remain unaddressed.

**Constraints & dependencies (28/30)**: Well-specified. The watcher constraint is now explicit about single-directory limitation and the switch strategy. The sessions-index.json 10% coverage is honestly quantified. The WatcherEventMsg offset prerequisite is clearly stated. Deduction: Still missing: maximum JSONL file size, corrupted JSONL handling.

**Score: 101/110**

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (28/40)**: The proposal now articulates three cross-domain ideas plus the sessions-index.json sidecar optimization. The Autopsy timeline analogy adds a genuine forensics-domain reference. The "discovery 合并到 ScanProjectsDir 遍历，避免 N+1 探测" remains the core design insight. Deduction: This is still primarily a catch-up feature batch. The innovation is in the implementation details (sidecar optimization, offset-based incremental parsing, tick-based debounce with max delay cap) rather than in the problem space or user interaction model.

**Cross-domain inspiration (30/35)**: Three cross-domain ideas:
1. IDE incremental indexing — "Watcher+ParseIncremental 复用'监视→增量解析→更新 UI'模式"
2. WAL replay — "JSONL 类似 WAL，offset 即 position，增量读行即 replay"
3. Digital forensics timeline (Autopsy/Sleuth Kit) — "Turn→ToolUse→Result 的层级展示思路"

The addition of Autopsy moves beyond adjacent domains into the specific problem domain (digital forensics). The Autopsy analogy is apt: organizing events as a browsable timeline mirrors the Turn→ToolUse→Result hierarchy. Deduction: The Autopsy analogy is present but could be deeper. Autopsy's timeline also supports filtering, tagging, and bookmarking — none of which are adopted or discussed.

**Simplicity of insight (23/25)**: The sessions-index.json piggyback on ScanProjectsDir traversal remains elegant. The offset-based incremental parsing leveraging JSONL's append-only nature is clean. The tick-based debounce with 2s absolute cap is a practical design that avoids both excessive refreshes and indefinite delays. These are genuine "why didn't I think of that" insights.

**Score: 81/100**

### 6. Feasibility (100 pts)

**Technical feasibility (38/40)**: All components within current tech stack. The watcher integration now addresses the race condition with a concrete mitigation (2s polling fallback + immediate ParseIncremental). The Phase 1 dependency rationale (key bug blocks all verification) adds a practical sequencing argument. Deduction: The atomicity of fsnotify watch switching is not discussed. The remove-watch/add-watch cycle may have platform-specific behavior (particularly on macOS, which is the validated platform).

**Resource & timeline feasibility (26/30)**: "16-24 总工时，含 20% 诊断缓冲" with 30min diagnostic timeboxes and escalation paths. The 20% buffer is explicit. Deduction: Still no calendar estimate. "16-24 hours" at different developer paces could mean 2 days or 2 weeks. The phase-level timeboxes are missing — how much time for Phase 1 vs. Phase 2 vs. Phase 3?

**Dependency readiness (28/30)**: All dependencies verified: Cobra, bubbletea, fsnotify, sessions-index.json format. The macOS fsnotify validation is noted. Deduction: The WatcherEventMsg offset field is listed as a prerequisite but its current implementation state is not specified — is the field already defined with value 0, or does it need to be added to the struct?

**Score: 92/100**

### 7. Scope Definition (80 pts)

**In-scope items are concrete (28/30)**: 8 items with item numbers, technical approaches, and prerequisite specifications. The watcher item now includes the monitoring strategy, switch handling, and offset prerequisite. The session discovery item includes the diagnostic gate. Deduction: Item 25 depends on a diagnostic phase that may reveal an out-of-scope root cause. The conditional fallback is in the SC but the scope item does not mention the fallback condition.

**Out-of-scope explicitly listed (23/25)**: 6 out-of-scope items, now including "手动/自动刷新切换" which is explicitly deferred. Deduction: Missing: backward compatibility of --session CLI flag with existing scripts that may pass unknown flags.

**Scope is bounded (21/25)**: Diagnostic timeboxing (30min per item, 20% buffer total) and escalation path ("创建 follow-up issue") prevent scope creep. The phase structure provides sequential bounding. Deduction: Still no calendar anchor. The phases are sequenced but not time-bounded in calendar terms. A sprint boundary or delivery date would strengthen commitment.

**Score: 72/80**

### 8. Risk Assessment (90 pts)

**Risks identified (28/30)**: 7 risks (up from 6 in iteration 2). New addition:
- "Watcher 会话切换期间的竞态条件 | M | M | 切换期间以 2s 间隔轮询目标目录直到 watch 建立，且切换时立即触发一次完整 ParseIncremental"

This addresses the blindspot from iteration 2. The mitigation is concrete and implementable. Deduction: The --session UUID search performance with 500+ project directories is not flagged. The WalkDir traversal with prefix matching is O(n) but the constant factor depends on filesystem speed.

**Likelihood + impact rated (27/30)**: Ratings are honest and consistent. The ScanProjectsDir risk (M/H) is realistic. The watcher race condition (M/M) is appropriately rated — it is a real concern but with an acceptable mitigation. The ParseIncremental offset risk (H/H) correctly signals a blocking prerequisite. Deduction: The "8 项并行修改回归面大" risk (M/H) mitigation says "每 Phase 完成后运行全量测试" — but this is standard practice, not a specific mitigation for 8-item regression risk. A more targeted mitigation would be per-item test coverage requirements.

**Mitigations are actionable (28/30)**: Most mitigations are specific and implementable. The watcher race condition mitigation is particularly well-specified: "切换期间以 2s 间隔轮询目标目录直到 watch 建立，且切换时立即触发一次完整 ParseIncremental." The diagnostic gate for ScanProjectsDir is concrete. The sessions-index.json mitigation includes version checking. Deduction: "先覆盖常见格式，异常显示原始内容" for TaskOutput is still vague — "常见格式" is not enumerated. How many TaskOutput formats exist in the wild?

**Score: 83/90**

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (28/30)**: Most SC are concrete and verifiable. The TaskOutput SC is now specific: "解析成功的 JSON/XML 内容按缩进和换行格式化显示；解析失败的原始内容按终端宽度自动换行显示." The watcher SC specifies two thresholds ("500ms 内或累计最多 2s"). The key bug SC specifies "80+ i18n 键." Deduction: The watcher SC's "写入暂停" concept — how is "pause" detected from the implementation? The debounce tick resets on each event, so "pause" means "no event for 500ms." This is implicit and should be stated explicitly.

**Coverage is complete (23/25)**: All 8 in-scope items have corresponding SC. The WatcherEventMsg offset prerequisite (listed in Constraints and Scope #8) still has no dedicated SC — if this prerequisite is not met, item 8 is blocked. Deduction: This was flagged in iteration 2 and remains unaddressed. A simple SC like "WatcherEventMsg 携带正确的 offset 值，handleWatcherEvent 传递给 ParseIncremental" would close the gap.

**SC internal consistency (24/25)**: All SC clusters are internally satisfiable. The TaskOutput SC tension from iteration 2 is fully resolved. The watcher SC is consistent with the debounce design and the race condition mitigation. No intra-SC contradictions detected. Deduction: The watcher SC says "编辑 .jsonl 文件后" — but the watcher monitors the session directory, not a specific file. If the directory contains multiple .jsonl files, does the SC apply to all of them or only the active session file? This is a minor ambiguity.

**Score: 75/80**

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (34/35)**: Full coverage of 8 defects. The Phase 1 rationale now explicitly links the execution order to the problem: key bug must be fixed first because it blocks verification of all other fixes. This closes the reasoning gap between P0-P3 urgency ranking and Phase 1-3 execution sequence. Deduction: Minor — the urgency section lists P0 as "数据完整性（item 25）" but this is Phase 2, not Phase 1. The Phase 1 rationale (blocking verification) is a valid override, but the urgency section could be updated to reflect this dependency explicitly.

**Scope ↔ Solution ↔ Success Criteria aligned (27/30)**: Improved alignment. The TaskOutput SC now matches the dual display mode described in the solution. The watcher risk mitigation aligns with the watcher SC thresholds. Deduction: The WatcherEventMsg offset prerequisite is listed in Scope #8 and Constraints but has no dedicated scope item or SC. This is a tracking gap — the prerequisite is a deliverable (code change) that should be tracked. If it proves difficult, item 8 is blocked without explicit visibility.

**Requirements ↔ Solution coherent (23/25)**: Clean mapping. The edge cases align with the implementation approaches. The NFRs map to the technical design. Deduction: The NFR "会话发现扫描 100+ 项目目录 < 2s" does not have a corresponding risk entry for performance degradation. If the WalkDir traversal with sessions-index.json lookups takes longer than 2s, no mitigation path is documented.

**Score: 84/90**

---

## Phase 3 — Blindspot Hunt

**[blindspot-1]** The sessions-index.json cost-benefit ratio is questionable. The proposal honestly states "仅约 10% 项目有，fallback 是常态" and "discovery 合并到 ScanProjectsDir 遍历." The implementation cost includes: JSON parsing during ScanProjectsDir, mapping construction, fallback logic, version checking, and testing for both paths. For 10% coverage, the benefit is "nicer session titles in the session list." The proposal should acknowledge this trade-off explicitly — is this feature worth the implementation and maintenance cost for a 10% hit rate? If the goal is to establish the infrastructure for future sessions-index.json adoption (as Claude Code adoption increases), this should be stated.

**[blindspot-2]** The "外部工具组合" alternative in the comparison table is rejected with: "取证分析需要会话列表浏览→Turn 详情→工具结果的闭环工作流，外部工具组合无法提供." However, this rejection conflates two different capabilities: (1) real-time file monitoring (watcher), which is genuinely harder externally, and (2) session discovery and UUID search, which could be composed with jq. The proposal does not consider a hybrid approach: build the watcher and UI enhancements in TUI, but use jq-compatible UUID search as an external filter. The all-or-nothing framing may overstate the case for in-house implementation of all 8 items.

**[blindspot-3]** The proposal states "Phase 2 含 2 个诊断项，每项时间上限 30min；超时未定位则创建 follow-up issue 不阻塞." This is a pragmatic approach, but it introduces an untracked deliverable: the "follow-up issue." If the ScanProjectsDir root cause is complex, the follow-up issue represents deferred scope. The proposal should state what happens to the session discovery SC if the follow-up issue is not resolved — does item 25 ship in a degraded state, or is it held?

**[blindspot-4]** The watcher monitors "当前会话目录" but the proposal does not specify how the "current session directory" is determined. Is it the directory of the currently selected session? If the user navigates away from a session in the TUI, does the watcher switch? If the user is viewing Turn details (not the session list), is the "current session" still being watched? The watcher's lifecycle coupling to UI navigation state is not described.

**[blindspot-5]** The Autopsy/Sleuth Kit analogy is present but shallowly explored. The proposal says "Autopsy 将文件系统事件组织为可浏览的时间线视图，启发了 Turn→ToolUse→Result 的层级展示思路." However, the proposal does not describe what the Turn→ToolUse→Result hierarchy looks like in the UI. Is it a tree view? An indented list? A collapsible accordion? The analogy is cited but the user-facing design is not specified beyond "层级展示." This is a design gap masquerading as an inspiration.

---

## Bias Detection Report

- Annotated regions: 7 attack points / 12 paragraphs = density 0.58
- Unannotated regions: 14 attack points / 24 paragraphs = density 0.58
- Ratio (annotated/unannotated): 1.00

Interpretation: Attack density is equal for annotated and unannotated regions (ratio 1.00), indicating no attention bias toward pre-revised content. The evaluation applied equal scrutiny to both revised and unrevised sections.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 2 |
|-----------|-------|-----|--------------------|
| Problem Definition | 99 | 110 | +6 |
| Solution Clarity | 115 | 120 | +6 |
| Industry Benchmarking | 100 | 120 | +21 |
| Requirements Completeness | 101 | 110 | +3 |
| Solution Creativity | 81 | 100 | +12 |
| Feasibility | 92 | 100 | +2 |
| Scope Definition | 72 | 80 | +3 |
| Risk Assessment | 83 | 90 | +4 |
| Success Criteria | 75 | 80 | +5 |
| Logical Consistency | 84 | 90 | +4 |
| **Total** | **902** | **1000** | **+66** |

---

## Top Attack Points (For Reference)

1. **[Industry Benchmarking]** The "仅修 bug" alternative remains a straw man — positioned only as "其余 6 项被推迟" with no argument for why phased delivery might be valid. Quote: *"仅修 bug (20, 25) | 最小改动 | 其余 6 项被推迟 | Rejected: 机会成本低应一并解决"* — Present a genuine phased delivery argument (e.g., ship P0/P1 in v1, defer P2/P3 to next iteration) rather than dismissing bug-fix-only as "opportunity cost is low."

2. **[Scope Definition]** No calendar anchor — "16-24 总工时" is effort, not schedule. Quote: *"单人开发，预计 16-24 总工时，含 20% 诊断缓冲"* — Add a calendar projection (e.g., "at 4h/day, 4-6 working days") to make the scope commitment concrete.

3. **[Requirements Completeness]** Missing NFRs for memory and startup performance. The enhanced ScanProjectsDir traversal with sessions-index.json loading has no memory budget. Quote: *"会话发现扫描 100+ 项目目录 < 2s"* — Add NFRs for memory impact and startup time of the enhanced discovery.

4. **[Success Criteria]** WatcherEventMsg offset prerequisite has no dedicated SC. If this prerequisite fails, item 8 is blocked without visibility. Quote: *"前置条件：WatcherEventMsg 需携带 offset 字段"* — Add SC: "WatcherEventMsg 携带正确 offset 值，handleWatcherEvent 传递给 ParseIncremental."

5. **[Solution Creativity]** The Autopsy analogy is cited but the UI design it inspired is not described. Quote: *"Autopsy 将文件系统事件组织为可浏览的时间线视图，启发了 Turn→ToolUse→Result 的层级展示思路"* — Specify what "层级展示" looks like in the TUI: tree view, indented list, or collapsible sections.
