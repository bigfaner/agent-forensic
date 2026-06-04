iteration: 0
title: "Pre-Revision (Freeform Findings)"
ATTACK_POINTS:
- "**[high]** sessions-index.json 仅存在于约 10% 的项目目录中，fallback 是常态而非边缘情况 | quote: \"there are approximately 102 project directories under `~/.claude/projects/`, but only about 10 of them contain a `sessions-index.json` file\" | improvement: 明确 sessions-index.json 的发现策略——在 ScanProjectsDir 遍历时顺带检查每个项目目录是否存在 sessions-index.json，构建 sessionId→summary 映射，避免 N+1 文件打开"
- "**[high]** 每次启动 102 次文件探测可能影响 <2s 性能目标 | quote: \"that is 102 file-open attempts (most returning FileNotFound) on every startup\" | improvement: 将 sessions-index.json 的探测合并到已有的 ScanProjectsDir 目录遍历中，作为同一次 I/O 的附带操作"
- "**[high]** 现有 watcher 仅监控单目录，无法覆盖嵌套的会话文件 | quote: \"the real gap is not message-loop wiring -- it is that the watcher watches the wrong granularity of the filesystem\" | improvement: 明确 watcher 集成策略——仅监控当前选中会话所在目录，在切换会话时更新 watcher 的监控路径"
- "**[high]** watcher 无 debounce 逻辑，高频写入将导致性能灾难 | quote: \"Without debounce, the TUI would receive and process dozens of events per second, each triggering a full incremental parse\" | improvement: 在 proposal 中明确 debounce 机制——在 Bubble Tea 消息循环中使用 tick 实现 500ms 合并"
- "**[high]** ParseIncremental offset 硬编码为 0，每次重解析整个文件 | quote: \"handleWatcherEvent calls parser.ParseIncremental(msg.FilePath, 0) with a hardcoded offset of 0\" | improvement: 在 proposal 中明确修复 WatcherEventMsg 结构体以传递 offset，并将此作为 watcher 集成的前置条件"
- "**[medium]** ScanProjectsDir bug 根因未知，调查可能无限膨胀 | quote: \"Without a concrete reproduction case (which sessions are missing and why), this item could balloon into an open-ended investigation\" | improvement: 在 Phase 2 前增加诊断步骤——运行 ScanProjectsDir 并与 find 结果对比，确认具体缺失的会话"
- "**[medium]** 按键 bug 可能不是大小写问题，现有代码已使用小写匹配 | quote: \"If uppercase keys are not responding, the problem might be in how Bubble Tea normalizes key events\" | improvement: 在 Phase 1 增加诊断步骤，确认按键 bug 的根因再指定修复方案"
BORDERLINE_FINDINGS:
- "**[medium]** --session UUID 搜索未指定具体策略 | 在文件名前缀匹配（可靠但需遍历）和 sessions-index.json 查询（快但覆盖率低）之间未做选择 | 建议：明确使用 filepath.WalkDir 文件名前缀匹配策略"
SKIPPED_FINDINGS:
- "**[low]** 建议：在 ScanProjectsDir 中顺带构建 sessionId→summary 映射 (subjective preference, covered by attack point 1)"
- "**[low]** 建议：watcher 改为仅监控当前会话目录 (subjective preference, covered by attack point 3)"
- "**[low]** 建议：明确 debounce 实现位置 (subjective preference, covered by attack point 4)"
- "**[low]** 建议：明确 UUID 查找策略 (subjective preference, covered by borderline finding)"
- "**[low]** 建议：先诊断按键 bug 根因 (subjective preference, covered by attack point 7)"
- "**[low]** 建议：传递 watcher offset (subjective preference, covered by attack point 5)"
rubric:
  all_dimensions: N/A
