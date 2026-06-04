# Proposal Evaluation Report — Iteration 2

**Proposal**: Session Experience Enhancement
**Date**: 2026-06-04
**Evaluator**: Adversary (CTO persona)
**Previous Iteration**: Iteration 1 (Score: 697/1000)

---

## Iteration-1 Issue Tracking

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | Industry benchmarking: no real tools cited | **Partially Addressed** | Three tools now referenced (lnav, jq, VisiData), but comparison table still uses only self-invented alternatives |
| 2 | Solution creativity: self-admitted zero innovation | **Addressed** | Cross-domain inspirations added (IDE incremental indexing, WAL replay) |
| 3 | Watcher SC contradicts debounce design | **Addressed** | SC now "500ms pause or 2s cumulative"; debounce NFR adds 2s cap |
| 4 | Session count SC assumes fixable root cause | **Addressed** | Conditional SC added: "若诊断确认权限/符号链接问题，则可访问目录计数匹配即可并记录 warning" |
| 5 | Timeline unanchored | **Partially Addressed** | Diagnostic time-boxes (30min) and escalation path added, but no total calendar estimate |
| 6 | Missing regression risk | **Addressed** | Risk row added: "8 项并行修改回归面大" with mitigation via per-phase full test runs |
| 7 | "取最新匹配" undefined | **Addressed** | Edge case now specifies "以 JSONL 文件 mtime 为准" |
| 8 | Generic urgency without quantification | **Addressed** | P0-P3 priority ranking with per-item time loss estimates |

**Resolution rate**: 5 fully addressed, 3 partially addressed, 0 unaddressed.

---

## Phase 1 — Reasoning Audit

### Problem → Solution Trace

All 8 defects map to in-scope items. The mapping is complete and unchanged from iteration 1. No orphan problems or phantom solutions.

### Solution → Evidence Trace

The urgency section now provides priority ranking (P0-P3) with concrete time-cost estimates ("手动翻页浪费 30-60s/次", "约 40% Turn 详情需外部工具"). This substantively improves the evidence chain for why these items matter. However, the evidence is still author-estimated rather than user-observed — there is no usage telemetry, user survey, or support ticket data.

### Evidence → Success Criteria Trace

Each in-scope item has a corresponding SC. The iteration-2 revisions resolved the two SC tensions identified in iteration 1:

1. **Watcher SC** (Cluster A): Now "写入暂停 500ms 内或累计最多 2s 后" — the 2s cap resolves the contradiction with the debounce mechanism. **Satisfiable**.

2. **Session count SC** (Cluster B): Now conditional — "若诊断确认权限/符号链接问题，则可访问目录计数匹配即可并记录 warning." **Satisfiable with documented fallback**.

### Self-Contradiction Check — SC Consistency Deep-Dive

**Cluster A: Watcher / Refresh**
- SC: "编辑 .jsonl 文件后，写入暂停 500ms 内或累计最多 2s 后，TUI 自动刷新"
- NFR: "收到 WatcherEventMsg 启动 500ms tick，同文件后续事件重置；设 2s 最大延迟上限"
- **Satisfiable**: The 2s absolute cap prevents indefinite delay under sustained writes. The SC language "累计最多 2s" matches the NFR's "2s 最大延迟上限." No contradiction.

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
- **Satisfiable**: Fallback is explicit and acknowledged as the common case.

**Cluster E: CLI --session UUID (New)**
- SC: "`--session <valid-uuid>` 启动后直接展示目标会话，无需翻页"
- Edge case: "UUID 在多个项目中存在 → 搜索所有项目，取最新匹配（以 JSONL 文件 mtime 为准）"
- In Scope #5: "filepath.WalkDir 文件名前缀匹配搜索 UUID 并直接打开"
- **Satisfiable**: The implementation approach (WalkDir with filename prefix match) supports the multi-project search behavior. The "最新" criterion is now defined (mtime).

**New Tension Detected — Cluster F: TaskOutput Display**

- In Scope #6: "解析 TaskOutput 工具调用内容并展示"
- Risk: "TaskOutput 格式多样解析不完整 | M | L | 先覆盖常见格式，异常显示原始内容"
- SC: "TaskOutput 调用结果在详情面板中可读展示"
- **Tension (minor)**: The SC says "可读展示" but the risk mitigation says "异常显示原始内容." Raw content may not be "readable." This is a minor ambiguity — the SC should clarify whether "可读" means "formatted for readability" or "visible" (even if raw). Tagged as **ambiguous — requires author clarification**.

---

## Phase 2 — Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

**Problem stated clearly (35/40)**: The problem is 8 specific defects with item references. The P0-P3 prioritization adds clarity about relative importance. Deduction: The three-category framing ("数据完整性、交互细节和内容展示") still does not perfectly map to the solution phases (Quick Fixes, Data Layer, UI Enhancement). The categorization is organizational rather than analytical — a reader might ask "is item 24 (title quality) a data integrity issue or a content display issue?" This is a minor framing concern, not a blocking issue.

**Evidence provided (33/40)**: Item references provide traceability. The urgency section now includes concrete time-cost estimates ("手动翻页浪费 30-60s/次", "约 40% Turn 详情需外部工具"). The P0-P3 ranking with per-item impact quantification is a significant improvement. Deduction: Evidence is still author-estimated rather than observed. No usage data, user complaints, or support tickets are cited. The "40%" figure for Turn details needing external tools is stated without source.

**Urgency justified (25/30)**: The P0-P3 framework with specific cost-per-incident estimates is a major improvement over the vague "频繁" claim in iteration 1. "按键每分钟多次触发" and "手动翻页浪费 30-60s/次" are concrete. "延迟修复 P0/P1 意味核心分析流程持续不可靠" is a reasonable cost-of-delay argument. Deduction: The urgency still lacks data on how many users are affected, how often forensic sessions occur, or whether there are workarounds that mitigate the impact.

**Score: 93/110**

### 2. Solution Clarity (120 pts)

**Approach is concrete (38/40)**: Three-phase breakdown with specific technical approaches per item. Each item names the implementation strategy (e.g., "filepath.WalkDir 文件名前缀匹配", "ScanProjectsDir 附带操作"). A reader could explain back what will be built. Deduction: Phase ordering rationale is still unstated — why is key bug Phase 1 rather than Phase 2? The P0-P3 urgency ranking would suggest data integrity (P0) should come first, but Phase 1 is "Quick Fixes" (interaction).

**User-facing behavior described (38/45)**: Key Scenarios section covers happy paths and edge cases. Improvements since iteration 1: the edge cases now specify "以 JSONL 文件 mtime 为准" for multi-project UUID resolution. Deduction: Some items still lack explicit user-facing behavior description:
- Item 18 (TaskOutput): "解析 TaskOutput 工具调用内容并展示" — how is it displayed? Inline in the detail panel? Collapsible? Truncated for long outputs?
- Item 25 (session discovery): "修复会话发现完整性" — the user sees more sessions, but is there a count indicator? A status message?

**Technical direction clear (33/35)**: Specific technical choices named (filepath.WalkDir, Cobra, fsnotify, tick-based debounce with 2s cap). The offset precondition for WatcherEventMsg is well-flagged as a prerequisite. Deduction: The key bug says "验证根因（key normalization vs handler routing）" — this identifies two hypotheses but does not discuss what the fix looks like for each scenario. If the root cause is "handler routing," what is the fix path?

**Score: 109/120**

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (28/40)**: Significant improvement. Three real tools are now cited:
- **lnav**: "实时 tail + 结构化日志" — the proposal adopts the "实时跟踪+结构化展示" pattern.
- **jq**: "JSON 流处理" — the `--session UUID` search is framed as a simplified jq filter, but file prefix matching is chosen as simpler.
- **VisiData**: "按需加载策略" — offset-based incremental parsing is inspired by VisiData's approach.

Each reference includes what is adopted and what is rejected. Deduction: The descriptions are thin — one sentence each. No version numbers, no links, no architectural comparison. "lnav" is described as "独立应用" but no detail on how lnav's architecture differs from agent-forensic's embedded TUI. The benchmarking is present but shallow.

**At least 3 meaningful alternatives (18/30)**: The comparison table still lists only self-invented alternatives: do nothing, bug-fix only, and full enhancement. "Bug-fix only" remains a straw man — it is positioned as "其余 6 项被推迟" with no argument for why this might be a valid strategy (e.g., ship P0/P1 fixes now, defer P2/P3 to next iteration). The three industry tools are referenced in a separate section but are not evaluated as alternatives. An alternative like "use lnav + jq externally instead of building into TUI" was not considered.

**Honest trade-off comparison (18/25)**: The table evaluates scope breadth only. "改动量较大" is still vague. The comparison does not address risk accumulation, timeline implications, or the trade-off between batch delivery vs. incremental delivery.

**Chosen approach justified against benchmarks (15/25)**: The justification is "各项独立，风险可控" — an assertion. The industry section explains what each tool inspired but does not justify why building all 8 features into the TUI is better than composing with external tools. The "参考三个工具" section is descriptive, not evaluative.

**Score: 79/120**

### 4. Requirements Completeness (110 pts)

**Scenario coverage (37/40)**: 9 scenarios covering happy paths, edge cases, and error scenarios. Improvements since iteration 1: the multi-project UUID edge case now specifies the resolution strategy (mtime). The sessions-index.json fallback is explicit. Missing scenario: What happens when the user switches sessions while a watcher event is being debounced? This concurrency interaction is still not covered. Also missing: corrupted JSONL files with partial or malformed lines — how does incremental parsing handle this?

**Non-functional requirements (33/40)**: Four NFRs with concrete numbers. The debounce NFR now includes the 2s cap, resolving the iteration-1 contradiction. Improvements: "tick-based debounce——收到 WatcherEventMsg 启动 500ms tick，同文件后续事件重置；设 2s 最大延迟上限" is specific and implementable. Deduction: Still missing: memory usage impact of loading sessions-index.json data for 100+ projects. Startup time impact of the enhanced ScanProjectsDir traversal. Performance characteristics of the --session UUID filesystem search across large project counts.

**Constraints & dependencies (28/30)**: Well-specified. Improvements: the watcher constraint now specifies "仅监控当前会话目录，切换时更新 watch target" and the WatcherEventMsg offset prerequisite is clearly stated. The sessions-index.json 10% coverage rate is honestly stated. Deduction: Still missing: maximum JSONL file size concern, corrupted JSONL handling.

**Score: 98/110**

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (22/40)**: Significant improvement. The proposal now articulates two cross-domain ideas and identifies the sessions-index.json sidecar optimization as the core design insight. The "discovery 合并到 ScanProjectsDir 遍历，避免 N+1 探测" is a legitimate optimization. Deduction: The innovation is modest — this is still primarily a catch-up feature batch. The cross-domain analogies (IDE incremental indexing, WAL replay) are reasonable but standard in log-processing applications.

**Cross-domain inspiration (25/35)**: Two cross-domain ideas are now cited:
1. IDE incremental indexing: "Watcher+ParseIncremental 复用'监视→增量解析→更新 UI'模式"
2. WAL replay: "JSONL 类似 WAL，offset 即 position，增量读行即 replay"

These are relevant and well-connected to the implementation. Deduction: Both analogies are from adjacent domains (developer tools, databases), not from truly different fields. Consider whether ideas from domains like version control (blame/annotate), network monitoring (packet replay), or digital forensics tools (Autopsy timeline) could further enrich the approach.

**Simplicity of insight (22/25)**: The sessions-index.json piggyback on ScanProjectsDir traversal is elegant — "利用 sessions-index.json 获取会话摘要（仅约 10% 项目有），discovery 合并到 ScanProjectsDir 遍历，避免 N+1 探测." The offset-based incremental parsing leveraging JSONL's append-only nature is also clean. These are "why didn't I think of that" insights.

**Score: 69/100**

### 6. Feasibility (100 pts)

**Technical feasibility (37/40)**: All components are within the current tech stack. The watcher integration strategy is now clearer: "仅监控当前会话目录，切换时更新 watch target." The offset prerequisite for WatcherEventMsg is flagged as H/H risk with a concrete fix path. Deduction: The watcher's add-watch/remove-watch cycle during session switching introduces a race condition window that is still not addressed. If the user rapidly switches sessions, the watcher could miss events during the transition.

**Resource & timeline feasibility (25/30)**: Improvement: "Phase 2 含 2 个诊断项，每项时间上限 30min；超时未定位则创建 follow-up issue 不阻塞" adds diagnostic time-boxing and an escalation path. This is a pragmatic approach. Deduction: Still no total calendar estimate. "8-12 个任务，每任务 1-2 小时" remains a rough range (8-24 hours). With diagnostic time-boxing, the range narrows but is still wide. No buffer is included for integration testing between phases.

**Dependency readiness (28/30)**: All dependencies verified. The fsnotify macOS validation is noted. The sessions-index.json format is verified. The Cobra framework is already integrated. Good.

**Score: 90/100**

### 7. Scope Definition (80 pts)

**In-scope items are concrete (27/30)**: 8 items with item numbers and technical approaches. Improvements: the watcher item now specifies the monitoring strategy ("仅监控当前会话目录，切换时更新 watch target") and its prerequisite (fix WatcherEventMsg offset). The session discovery item specifies the diagnostic gate. Deduction: Item 25 still depends on a diagnostic phase that may reveal an out-of-scope root cause. The conditional fallback is now in the SC but not in the scope item itself.

**Out-of-scope explicitly listed (22/25)**: 6 out-of-scope items. Improvements: "手动/自动刷新切换（todo 23 的 UI 开关部分）" is now explicitly deferred. Deduction: Still missing: backward compatibility of the --session CLI flag with existing scripts. Memory usage impact from sessions-index.json loading.

**Scope is bounded (20/25)**: The diagnostic time-boxing (30min per diagnostic item) adds boundary. The "超时未定位则创建 follow-up issue" escalation path prevents scope creep. Deduction: Still no calendar anchor. The phases are sequenced but not time-bounded in calendar terms. Is this a sprint? A week? Open-ended?

**Score: 69/80**

### 8. Risk Assessment (90 pts)

**Risks identified (27/30)**: 6 risks identified (up from 4 in iteration 1). New additions:
- "8 项并行修改回归面大 | M | H" — addresses the missing regression risk from iteration 1.
- "ParseIncremental offset 硬编码为 0 | H | H" — elevates a known defect to risk status.

Deduction: Still missing: watcher race condition during rapid session switching (add-watch/remove-watch window). The --session UUID search performance with large project counts is not flagged as a risk.

**Likelihood + impact rated (25/30)**: Ratings are honest. The ScanProjectsDir risk (M/H) is realistic. The sessions-index.json risk (L/M) is appropriately low given the fallback is the norm. Deduction: The ParseIncremental offset risk (H/H) is correctly identified as a blocking prerequisite. The mitigation ("前置：WatcherEventMsg 携带 offset，handleWatcherEvent 传递而非硬编码") is concrete. However, H/H items in a risk table typically signal a blocking dependency, not a risk — this is a known defect that must be fixed before item 8 can proceed. It would be better tracked as a prerequisite than a risk.

**Mitigations are actionable (27/30)**: Improvements: the regression mitigation now specifies "每 Phase 完成后运行全量测试（just test），通过后才进入下一 Phase." The diagnostic gate for ScanProjectsDir is concrete ("对比 ScanProjectsDir 输出与 find 结果"). The sessions-index.json mitigation specifies "版本号检查 + graceful fallback." Deduction: "先覆盖常见格式，异常显示原始内容" for TaskOutput is still vague — what is a "常见格式"? How many formats exist? Is there an enumeration?

**Score: 79/90**

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (26/30)**: Most SC are testable with concrete verification methods. Improvements:
- Watcher SC: "写入暂停 500ms 内或累计最多 2s 后，TUI 自动刷新" — now has two measurable thresholds.
- Session count SC: conditional with documented fallback.
- Key bug SC: "覆盖 80+ i18n 键" is quantifiable.

Deduction: "TaskOutput 调用结果在详情面板中可读展示" — "可读" is subjective. What makes it readable? Formatted? Syntax highlighted? Truncated? The watcher SC has two conditions ("写入暂停 500ms 内" OR "累计最多 2s") — is this an OR condition? How is "写入暂停" detected from the debounce tick?

**Coverage is complete (22/25)**: All 8 in-scope items have corresponding SC. Improvement: the conditional session count SC provides a fallback path. Deduction: The WatcherEventMsg offset prerequisite (listed in Scope #8 and Constraints) has no dedicated SC. If this prerequisite is not met, item 8 is blocked, but there is no SC to verify it.

**SC internal consistency (22/25)**: The two iteration-1 contradictions are resolved. The SC set is now internally consistent. The minor tension on TaskOutput "可读展示" vs. "异常显示原始内容" is a definitional ambiguity, not a logical contradiction. No intra-SC contradictions detected.

**Score: 70/80**

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (32/35)**: Full coverage of 8 defects. The P0-P3 priority ranking in the urgency section aligns roughly with the phased delivery — Phase 1 addresses P1 (key bug), Phase 2 addresses P0 (session discovery) and P1 (CLI param), Phase 3 addresses P2 (full conversation display, watcher). Deduction: The phase ordering does not strictly follow priority — P0 (session discovery) is in Phase 2, not Phase 1. The proposal should explain why Quick Fixes are Phase 1 despite being P1/P3 items. Is there a dependency? Or is it strictly "easiest first"?

**Scope ↔ Solution ↔ Success Criteria aligned (26/30)**: Improved alignment. The conditional SC for session discovery resolves the scope-SC mismatch. The watcher SC now matches the debounce design. Deduction: The WatcherEventMsg offset prerequisite is listed in Scope #8's description and in Constraints, but has no dedicated scope item or SC. If this prerequisite proves difficult, item 8 is blocked without explicit tracking. This is a tracking gap, not a logical inconsistency.

**Requirements ↔ Solution coherent (22/25)**: Clean mapping. The edge case "UUID 在多个项目中存在 → 搜索所有项目，取最新匹配（以 JSONL 文件 mtime 为准）" now aligns with the implementation approach. Deduction: The NFR "会话发现扫描 100+ 项目目录 < 2s" does not have a corresponding risk entry for performance degradation. If the WalkDir traversal with sessions-index.json lookups takes longer than 2s, there is no mitigation path documented.

**Score: 80/90**

---

## Phase 3 — Blindspot Hunt

**[blindspot-1]** The comparison table evaluates only scope breadth (do nothing vs. bug-fix only vs. full enhancement) but never considers approach alternatives. The proposal references lnav, jq, and VisiData as inspirations but never evaluates "use these tools externally" as an alternative to building all 8 features into the TUI. For a forensic analyst, composing agent-forensic's session discovery with lnav's real-time viewing might be a valid strategy that avoids the complexity of watcher integration. The proposal assumes all features must be built in-house without justifying this assumption.

**[blindspot-2]** The watcher integration introduces a subtle operational complexity that is underexplored. The proposal specifies "仅监控当前会话目录，切换时更新 watch target" but does not discuss what happens during the transition window. When the user switches from session A to session B:
1. Remove watch on session A directory
2. Add watch on session B directory
3. Between steps 1 and 2, writes to session B are missed

This is a windowed data loss scenario. Under rapid session switching (e.g., comparing two sessions), the user could miss updates. The proposal should discuss whether this is acceptable or whether a watch pool is needed.

**[blindspot-3]** The "Assumptions Challenged" section is a positive addition, but it reveals a concerning pattern: two of three assumptions were "Overturned" (sessionName field does not exist; session list does not load all sessions). This suggests the proposal was written with incomplete understanding of the codebase. The proposal should acknowledge this and state whether additional assumptions have been validated before proceeding to implementation.

**[blindspot-4]** The proposal states "sessions-index.json 仅约 10% 项目有，fallback 是常态" — this means 90% of sessions will not benefit from the sessions-index.json optimization. The effort to implement this optimization (JSON parsing, mapping construction, fallback logic) may not justify the 10% improvement. The proposal does not include a cost-benefit analysis for this specific item. Quote: *"sessions-index.json 仅约 10% 项目有，fallback 是常态；discovery 必须合并到 ScanProjectsDir 遍历."*

**[blindspot-5]** The NFR "会话发现扫描 100+ 项目目录 < 2s" sets a performance target, but the proposal does not discuss what happens if this target is not met. Is there a fallback strategy? Loading indicators? Background scanning? The 2s target is stated as a requirement but treated as an assumption that the current approach will meet it.

---

## Bias Detection Report

- Annotated regions: 10 attack points / 12 paragraphs = density 0.83
- Unannotated regions: 18 attack points / 26 paragraphs = density 0.69
- Ratio (annotated/unannotated): 1.20

Interpretation: Attack density is slightly higher for annotated regions (ratio 1.20), suggesting mild attention bias toward pre-revised content. However, the difference is small and does not indicate significant unfairness. Several attacks on annotated regions are confirmations that revisions resolved issues, which inflates the count.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 1 |
|-----------|-------|-----|--------------------|
| Problem Definition | 93 | 110 | +13 |
| Solution Clarity | 109 | 120 | +9 |
| Industry Benchmarking | 79 | 120 | +34 |
| Requirements Completeness | 98 | 110 | +8 |
| Solution Creativity | 69 | 100 | +26 |
| Feasibility | 90 | 100 | +7 |
| Scope Definition | 69 | 80 | +6 |
| Risk Assessment | 79 | 90 | +10 |
| Success Criteria | 70 | 80 | +16 |
| Logical Consistency | 80 | 90 | +10 |
| **Total** | **836** | **1000** | **+139** |

---

## Top Attack Points (Prioritized for Revision)

1. **[Industry Benchmarking]** Comparison table still uses only self-invented alternatives — the three industry tools (lnav, jq, VisiData) are referenced as inspirations but never evaluated as alternatives. Quote: *"参考三个工具：(1) lnav... (2) jq... (3) VisiData..."* followed by a comparison table that lists only "Do nothing", "仅修 bug", and "完整增强". — Include an industry-validated alternative in the comparison table (e.g., "Use lnav/jq externally for advanced features, build only core fixes into TUI") and justify why the in-house approach is preferred.

2. **[Industry Benchmarking]** Tool descriptions are thin — one sentence each with no architectural comparison. Quote: *"lnav——实时 tail + 结构化日志，采纳'实时跟踪+结构化展示'理念但它是独立应用，仅借鉴模式"* — Expand each reference to 2-3 sentences: what specific feature/pattern is adopted, what is rejected, and why.

3. **[Solution Creativity]** Cross-domain analogies are from adjacent domains only (developer tools, databases). Quote: *"借鉴两个跨领域思路：(1) IDE 增量索引...(2) WAL replay..."* — Consider whether ideas from digital forensics tools (e.g., Autopsy's timeline view, Sleuth Kit's metadata layering) or version control (git blame's incremental annotation) could further enrich the approach.

4. **[Scope Definition]** No calendar anchor for the timeline. Quote: *"单人开发，预计 8-12 个任务，每任务 1-2 小时"* — Add a total time estimate with buffer (e.g., "2-3 weeks with 20% buffer for diagnostics") and phase-level timeboxes.

5. **[Solution Clarity]** Phase ordering does not follow P0-P3 priority. P0 (session discovery) is Phase 2, P1 (key bug) is Phase 1. Quote: *"Phase 1 — Quick Fixes：按键 bug...诊断面板加会话标题"* vs. urgency section *"P0 数据完整性（item 25）"* — Either reorder phases to match priority, or explicitly explain why Quick Fixes come first (e.g., "Phase 1 items are low-risk and unblock development velocity for Phase 2").

6. **[Risk Assessment]** Missing watcher race condition during rapid session switching. The proposal specifies "切换时更新 watch target" but does not discuss the transition window where events can be missed. — Add a risk entry for this race condition with a mitigation (e.g., "during transition, poll session directory on a 2s interval until watch is established").

7. **[Success Criteria]** TaskOutput SC uses subjective "可读展示." Quote: *"TaskOutput 调用结果在详情面板中可读展示"* vs. risk mitigation *"先覆盖常见格式，异常显示原始内容"* — Define "可读": either change to "展示解析后内容或原始内容" or specify that "可读" means "formatted with line breaks and indentation."
