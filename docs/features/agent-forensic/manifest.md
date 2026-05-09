---
feature: "agent-forensic"
status: tasks
---

# Feature: agent-forensic

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | Agent 行为诊断 TUI：调用树、回放、实时监听、仪表盘、异常标记、证据提取、i18n |
| User Stories | prd/prd-user-stories.md | 8 个用户故事覆盖浏览、异常定位、详情、搜索、回放、实时监听、仪表盘、边界条件 |
| UI Functions | prd/prd-ui-functions.md | 6 个 UI 功能：Sessions、Call Tree、Detail、Dashboard、Diagnosis、Status Bar |
| UI Design | ui/ui-design.md | Vercel 风格 TUI 设计：6 个组件的布局、状态、交互和数据绑定规格 |
| Tech Design | design/tech-design.md | Go + Bubble Tea 架构：6 个 Bubble Tea model、JSONL parser、file watcher、anomaly detector、sanitizer、i18n |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| UF-1 Sessions Panel | model/sessions.go, parser/jsonl.go | Sessions Panel (ui-design §UF-1) | new-page:main-tui | 1.1, 2.1, 3.1, 4.1, 4.2 |
| UF-2 Call Tree Panel | model/calltree.go, detector/anomaly.go | Call Tree Panel (ui-design §UF-2) | new-page:main-tui | 1.1, 2.1, 2.2, 3.2, 4.1 |
| UF-3 Detail Panel | model/detail.go, sanitizer/sanitizer.go | Detail Panel (ui-design §UF-3) | new-page:main-tui | 1.1, 2.3, 3.3, 4.1 |
| UF-4 Dashboard View | model/dashboard.go, stats/stats.go | Dashboard View (ui-design §UF-4) | new-page:dashboard | 1.1, 2.5, 3.4, 4.1 |
| UF-5 Diagnosis Summary | model/diagnosis.go, detector/anomaly.go | Diagnosis Summary (ui-design §UF-5) | new-page:diagnosis | 1.1, 2.2, 3.5, 4.1 |
| UF-6 Status Bar | model/statusbar.go | Status Bar (ui-design §UF-6) | new-page:main-tui | 2.4, 3.6, 4.1 |
| i18n | i18n/i18n.go, i18n/locales/*.yaml | All components | — | 2.4, 4.2 |
| Real-time monitoring | watcher/watcher.go | Call Tree (new nodes) | — | 2.6, 3.2, 4.1 |
