# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

## 5. Key Documents

```
docs/
├── ARCHITECTURE.md              # 分层架构总览（按需加载）
├── conventions/                 # 项目级规范（按需加载，不自动消耗 token）
│   └── INDEX.md                 # ← 规范索引，按 scope/keyword 查找文件
├── lessons/                     # 经验教训（按需加载）
│   └── INDEX.md                 # ← 按 category/task-trigger 查找
├── proposals/<slug>/            # 早期提案（brainstorm 输出）
└── features/<slug>/             # Feature 工作区
    ├── manifest.md              # Feature 索引 & 可追溯性映射
    ├── prd/                     # 需求文档
    ├── design/                  # 技术设计 & API 文档
    ├── ui/                      # UI 设计规格
    └── tasks/                   # 任务定义 & 执行记录
```

### Convention Loading

执行任务前，**必须**按文件模式加载对应规范（不是"按需"，是强制）：

**TUI / 前端任务**（`internal/model/*.go`，或任何导入 `bubbletea`/`lipgloss` 的文件）：
1. 读取 `docs/conventions/` 下全部 5 个文件（总计 ~800 行）
2. 读取 `docs/lessons/lesson-tui-tech-design-mockup.md` 和 `lesson-tui-visual-verify.md`
3. 编码前，用一句话陈述哪些 convention 章节适用于当前任务

**Parser / 数据任务**（`internal/parser/`、`internal/stats/`）：
1. 读取 `docs/conventions/tui-data-contracts.md`

**TUI Mockup 要求**：编写新的 `View()`/`Render()` 函数前，先输出 ASCII mockup 供用户确认。Mockup 必须包含：
- 每个面板的 box-drawing 字符图示
- 精确尺寸数值（禁止"大约"/"适当"）
- 5 个边界场景：窄终端(80×24)、宽终端(140+)、混合数字宽度、长路径(>50 chars)、无数据
- 字符调色板（每个视觉元素指定 Unicode 字符 + code point）
- 颜色从 `tui-layout-ui.md` 色表选取

**Scope Guard**：Feature tasks 完成后的 fix 迭代阶段：
- 提交前运行 `git diff --stat HEAD`，如任何 .go 文件 +50 行非 test 代码 → 停下，创建 task
- 同一文件累积 3 个 fix/style 提交 → 停下，运行 `/learn-lesson` 提取 convention
- 连续 vibe coding 提交达 5 个 → 暂停，检查是否存在结构性问题