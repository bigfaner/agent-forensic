# Freeform Review: Session Experience Enhancement

**Reviewer Persona:** TUI Forensic Tooling & Data Integration Specialist
**Date:** 2026-06-04

---

## Section 1: Background Assessment

This proposal addresses eight concrete deficiencies in the agent-forensic TUI that collectively undermine the tool's core mission: letting analysts reliably inspect Claude Code sessions. The problems range from trivial interaction bugs (key bindings only respond to uppercase letters) to fundamental data completeness failures (some sessions simply do not appear in the session list). The author has clearly spent time with the actual product -- the evidence section cites specific issue numbers and the "Assumptions Challenged" table shows they did genuine homework, overturning the assumption that a `sessionName` field existed in Claude metadata and discovering that sessions-index.json's `summary` field is the closest available alternative.

The core technical approach is a three-phase rollout: quick interaction fixes first, then data layer improvements, then UI enhancements. This sequencing makes sense -- fix the buttons before rewiring the data, then extend the display. The proposal leans on Claude Code's existing `sessions-index.json` file as a data source for session titles, which is a pragmatic choice since that file already contains a `summary` field per session. For the watcher integration, the proposal plans to connect an already-implemented `watcher.go` to the Bubble Tea message loop. For the `--session` CLI flag, it relies on the existing Cobra framework.

The proposal rests on several assumptions that warrant scrutiny. First, it assumes `sessions-index.json` has a stable schema across Claude Code versions, which the risk table acknowledges as "Likelihood: L, Impact: M" but which my filesystem investigation suggests is more nuanced -- the file does not exist in every project directory. Second, it assumes the `ScanProjectsDir` bug has a discoverable root cause that can be fixed within the project's scope, but the current code at `/Users/fanhuifeng/Projects/ai/agent-forensic/internal/parser/jsonl.go:205-240` is a straightforward `filepath.WalkDir` that only skips `subagents/` subdirectories, so the bug may be more subtle than the proposal anticipates. Third, it treats the watcher integration as a straightforward wiring task, but the existing watcher watches a single directory while sessions span dozens of project subdirectories.

---

## Section 2: Key Risks

**Risk: sessions-index.json does not exist in all project directories.**

The proposal states: "从 sessions-index.json 读取 summary 作为标题，fallback 到当前行为" and "sessions-index.json 不是所有项目都有，需要 fallback". It also says "Edge case: sessions-index.json 不存在或 summary 为空 → 回退到当前行为（首条用户消息）".

I verified the actual filesystem: there are approximately 102 project directories under `~/.claude/projects/`, but only about 10 of them contain a `sessions-index.json` file. The agent-forensic project itself -- the very tool we are enhancing -- does not have one. This means the fallback path is not an edge case; it is the common case. The proposal underweights how often the fallback triggers and does not specify how the system discovers which projects have index files. If the code opens `sessions-index.json` per-project and parses it during session loading, that is 102 file-open attempts (most returning FileNotFound) on every startup. The proposal says "sessions-index.json 解析不阻塞 UI（异步加载）" but does not address the I/O pattern of probing for a file that exists in fewer than 10% of directories.

This matters because session loading already happens in a background `tea.Cmd`, and adding 102 speculative file-open syscalls plus JSON parsing for the 10 that exist could measurably affect the "< 2s" non-functional requirement for scanning 100+ directories. The proposal's success criterion "会话标题使用 summary 字段；不存在时回退到首条用户消息" will silently degrade to the fallback for over 90% of sessions.

**Risk: ScanProjectsDir bug root cause may be outside the code.**

The proposal states: "排查 ScanProjectsDir 逻辑，可能涉及递归深度或过滤条件" and "ScanProjectsDir 的 bug 根因复杂（如权限、符号链接）".

Having read the actual implementation, `ScanProjectsDir` is a 35-line function that does `filepath.WalkDir` on the projects directory, skipping only `subagents/` subdirectories and collecting `.jsonl` files. There is no recursion depth limit, no symlink handling, and no permission-checking logic. If sessions are genuinely missing, the root cause is likely not in this function -- it could be that sessions exist outside the `~/.claude/projects/` hierarchy (e.g., in worktree directories), or that the missing sessions are in directories that `os.ReadDir` silently skips due to permission errors, or that the file extension differs. The proposal's risk table rates this as "M" likelihood and "H" impact, but the mitigation "先诊断具体原因再决定修复方案" is vague. Without a concrete reproduction case (which sessions are missing and why), this item could balloon into an open-ended investigation that blocks Phase 2.

This matters because the proposal's success criterion is "会话列表显示的项目数 >= `find ~/.claude/projects -name '*.jsonl' | wc -l`" -- on the current machine, that is 2,345 files. If `ScanProjectsDir` returns all 2,345, the bug may already be fixed, or the bug may be in the subsequent `maxRecentSessions = 20` truncation or in the parsing step, not in discovery.

**Risk: Watcher monitors a single directory, not the entire projects tree.**

The proposal states: "Watcher 依赖 fsnotify，macOS 已验证可用" and "接入文件监视器实现自动刷新 + 手动刷新键". The success criterion says: "编辑 .jsonl 文件后 1s 内 TUI 自动刷新对应会话数据".

The existing `watcher.go` constructor takes a single `dir string` and adds only that directory to fsnotify: `w.fsw.Add(w.dir)`. It does not recursively watch subdirectories. Sessions live at `~/.claude/projects/{encoded-project-path}/{uuid}.jsonl` -- which means they are nested two or more levels deep. To watch all sessions, the watcher would need to monitor each project subdirectory individually, or monitor the top-level `projects/` directory and handle the fact that fsnotify on macOS (using kqueue) may not propagate events from deeply nested paths. The proposal does not address this architectural mismatch at all. It says "watcher.go 已实现，需连接 Bubble Tea 消息循环" but the real gap is not message-loop wiring -- it is that the watcher watches the wrong granularity of the filesystem.

This matters because the proposal's watcher integration might work for a single in-progress session (where you know the exact directory) but will fail for the broader use case of detecting changes across all sessions.

**Risk: --session UUID cross-project search requires scanning all 2,345 JSONL files or probing 102 directories.**

The proposal states: "CLI `--session <UUID>` 参数（item 17）：通过 UUID 搜索并直接打开指定会话" and "Edge case: UUID 在多个项目中存在 → 搜索所有项目，取最新匹配" and "会话发现扫描 100+ 项目目录应在 < 2s 内完成".

The challenge is that the UUID is embedded in the filename (`{uuid}.jsonl`), not in any index. To find a UUID, the system must either scan all 2,345 filenames (which is an O(n) directory walk) or parse each `sessions-index.json` to check `sessionId` fields. The former is fast (filesystem metadata only) but still requires walking all subdirectories. The latter only works for the ~10% of projects that have an index file. The proposal does not specify which strategy it will use, nor does it account for the fact that `--session` implies a synchronous CLI startup path -- the user types the command and expects the TUI to open on the target session immediately, without the normal async loading flow. The non-functional requirement of "< 2s" is for "scanning 100+ project directories," but finding a specific UUID is a different operation than listing all sessions.

This matters because the user experience depends on this being fast. If the implementation naively calls `ScanProjectsDir` (which walks the entire tree) and then filters by UUID, it will work but may be slow on systems with many projects. If it tries to use `sessions-index.json`, it will miss sessions in unindexed projects.

**Risk: Watcher integration lacks debounce despite the non-functional requirement specifying it.**

The proposal states: "Watcher 事件去重：同一文件 500ms 内多次写入合并为一次刷新". However, searching the codebase for "debounce" returns zero results. The existing `watcher.go` has no debounce logic -- it emits a `WatchEvent` for every `fsnotify.Write` event immediately. The `handleWatcherEvent` function in `app.go` at line 803-831 calls `parser.ParseIncremental` on every watcher event, which re-parses from offset 0 each time (the offset argument is hardcoded to 0). There is no batching, no timer, and no deduplication.

This matters because during active Claude Code sessions, JSONL files receive rapid sequential writes (tool calls, results, streaming text). Without debounce, the TUI would receive and process dozens of events per second, each triggering a full incremental parse. This could cause rendering jank and wasted CPU. The proposal acknowledges the risk in its risk table ("Watcher 在高频写入场景下产生过多事件, M/M") but the mitigation "500ms debounce 合并" is a design aspiration, not something that exists in the code.

**Problem: ParseIncremental is called with offset 0 in handleWatcherEvent.**

The proposal does not address this at all, but examining the existing watcher integration code reveals a bug: `handleWatcherEvent` at line 815 calls `parser.ParseIncremental(msg.FilePath, 0)` with a hardcoded offset of 0, meaning every watcher event triggers a re-parse of the entire file from the beginning. The watcher's `WatchEvent` struct includes an `Offset` field that tracks where new content begins, but this offset is never passed through the `WatcherEventMsg` to the handler. The `WatcherEventMsg` struct only contains `FilePath` and `Lines`, discarding the offset information. This means the "incremental" refresh is not incremental at all -- it re-reads the whole file on every event.

This matters because for large session files, full re-parsing on every filesystem event will be a performance disaster, especially combined with the lack of debounce.

**Problem: Key binding issue may not be a simple case-sensitivity problem.**

The proposal states: "按键仅支持大写，小写不响应（item 20）" and the solution is "所有 key binding 同时匹配大写和小写字母". However, examining the actual code at `internal/model/calltree.go:372-399`, the key matching uses `msg.String()` which returns lowercase strings like `"down"`, `"up"`, `"n"`, `"p"`, `"s"`, `"m"`. The switch cases already use lowercase. If uppercase keys are not responding, the problem might be in how Bubble Tea normalizes key events, or in a parent handler intercepting uppercase keys before they reach the panel-specific handlers. The proposal's solution "修改 key matching 逻辑" assumes the fix is to add uppercase cases, but the existing code is already lowercase. The real bug may be elsewhere.

This matters because the proposal may be fixing the wrong thing. If the issue is in event routing rather than case matching, adding uppercase aliases will not solve it.

---

## Section 3: Improvement Suggestions

**建议： Specify the sessions-index.json discovery strategy explicitly.**

Instead of the vague "从 sessions-index.json 读取 summary", the proposal should specify: during `ScanProjectsDir`, for each project directory encountered, check if `sessions-index.json` exists; if so, parse it and build a map of `sessionId -> summary`. This avoids N+1 file-open syscalls during session title rendering. The map can be populated once during the initial scan and reused. This addresses the risk that 90% of projects lack the index file by making the probe a natural side-effect of the directory walk rather than a separate operation. After adopting this, the proposal would describe a concrete data structure (e.g., `map[string]string` from session UUID to summary) and a loading path that integrates with the existing `SessionsLoadedMsg`.

**建议： Require a concrete reproduction for ScanProjectsDir before Phase 2.**

Before committing to "修复会话列表完整性" as a Phase 2 deliverable, the proposal should include a diagnostic step: run `ScanProjectsDir` on the actual `~/.claude/projects/` directory, compare the count against `find ~/.claude/projects -name '*.jsonl' | wc -l`, and identify which specific sessions are missing. If the counts match, the bug is elsewhere (perhaps in parsing, or in the `maxRecentSessions = 20` truncation, or in the user's perception of what "all sessions" means). This addresses the risk of an open-ended investigation by establishing a clear go/no-go gate. After adopting this, Phase 2 would begin with a diagnostic task that either confirms the bug location or reclassifies the issue.

**建议： Redesign the watcher integration to handle multi-directory monitoring.**

The proposal should acknowledge that the current `watcher.go` monitors a single directory and specify how it will be extended. The most practical approach is to watch the specific directory of the currently-selected session (which is known from `currentSession.FilePath`). This is sufficient for the primary use case (monitoring an active session) without requiring recursive directory watching. The watcher's `dir` field would be updated whenever the user selects a different session. This addresses the risk of trying to watch the entire projects tree and the architectural mismatch with fsnotify's semantics. After adopting this, the watcher integration would be scoped to "watch the current session's directory" rather than "watch all sessions."

**建议： Pass the watcher offset through WatcherEventMsg and fix ParseIncremental usage.**

The `WatcherEventMsg` struct should include the `Offset` field from the watcher's `WatchEvent`, and `handleWatcherEvent` should pass this offset to `ParseIncremental` instead of hardcoding 0. This is a prerequisite for any watcher integration to work efficiently. The proposal should call out this existing bug explicitly. This addresses the risk of full-file re-parsing on every watcher event. After adopting this, the incremental refresh pipeline would be: watcher detects new bytes at offset N -> WatcherEventMsg carries offset N -> ParseIncremental reads only from offset N -> new entries appended to call tree.

**建议： Add debounce to the watcher event pipeline.**

The proposal's non-functional requirement specifies "500ms debounce" but the implementation plan does not mention where this debounce logic will live. It should be specified: either in `watcher.go` (coalescing events in a time window before emitting to the channel) or in `app.go` (using a Bubble Tea tick to batch incoming watcher events). The Bubble Tea approach is more idiomatic: when a `WatcherEventMsg` arrives, start a 500ms tick; if another event for the same file arrives before the tick fires, reset the timer; when the tick fires, process the accumulated lines. This addresses the risk of event flooding during active sessions. After adopting this, the proposal would describe the debounce mechanism and its interaction with the message loop.

**建议： Specify the --session UUID lookup strategy.**

The proposal should state whether UUID lookup will use filename matching (fast, O(n) directory walk but covers all sessions) or sessions-index.json lookup (fast for indexed sessions but misses 90% of them), or a hybrid approach (check index files first, fall back to directory walk). Given that filenames embed the UUID as `{uuid}.jsonl`, the most reliable approach is `filepath.WalkDir` with a filename prefix match, which requires walking 102 directories but avoids reading file contents. This addresses the risk of the UUID search being incomplete or slow. After adopting this, the proposal would include a concrete algorithm: "Walk ~/.claude/projects/ looking for files matching `{uuid}.jsonl`; if multiple matches found, select by newest mtime."

**建议： Investigate the key binding bug before specifying the fix.**

Before committing to "所有 key binding 同时匹配大写和小写字母", the proposal should include a diagnostic step to determine why lowercase keys work in some panels but not others. The existing code already uses lowercase case labels (`"n"`, `"p"`, `"s"`, `"m"`), so if uppercase keys are the only ones responding, the issue may be in how terminal key events are normalized or in a parent handler's key routing logic. This addresses the problem that the proposed fix may target the wrong layer. After adopting this, Phase 1's key binding fix would begin with a diagnostic task that identifies the actual root cause.
