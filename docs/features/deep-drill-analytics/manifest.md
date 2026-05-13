---
feature: "deep-drill-analytics"
status: tasks
---

# Feature: deep-drill-analytics

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | SubAgent 下钻、文件追踪、Hook 增强、Turn 效率、重复检测、思考链、成功率 7 个分析维度，分 Phase 1/2 交付 |
| User Stories | prd/prd-user-stories.md | 9 个用户故事覆盖 Phase 1（SubAgent 下钻、文件统计、Hook 分析）和 Phase 2（效率、重复、思考链、成功率） |
| UI Functions | prd/prd-ui-functions.md | 6 个 UI Function：SubAgent 内联展开、全屏视图、Turn 文件列表、SubAgent 统计、文件排行面板、Hook 分析面板 |
| UI Design | ui/ui-design.md | 6 个组件设计，继承现有 TUI 设计系统，新增 SubAgent 展开/overlay、文件操作面板、Hook 分析面板 |
| Tech Design | design/tech-design.md | parser/stats/model 三层扩展，4 个新接口，5 个新数据模型，6 个集成点，懒加载 SubAgent |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| P1-1 SubAgent Drill-down | Interface 1 (SubAgent Loading), Integration 1/3/6 | UF-1 SubAgent Inline Expand | existing-page:Call Tree 面板 | 1.1, 2.1, 3.1 |
| P1-1 SubAgent Drill-down | Interface 4 (Extended Stats), Integration 6 | UF-2 SubAgent Full-Screen Overlay | new-page:SubAgent Analysis Overlay | 1.4, 2.2, 3.2 |
| P1-1 SubAgent Drill-down | Interface 4 (Extended Stats), Integration 3 | UF-4 SubAgent Statistics View | existing-page:Detail 面板 | 1.4, 2.4, 3.3 |
| P1-2 File Read/Write Tracking | Interface 2 (File Path Extraction), Integration 2 | UF-3 Turn File Operations | existing-page:Detail 面板 | 1.2, 2.3, 3.3 |
| P1-2 File Read/Write Tracking | Interface 2 (File Path Extraction), Integration 4 | UF-5 Dashboard File Operations | existing-page:Dashboard overlay | 1.2, 2.5, 3.4 |
| P1-3 Hook Analysis Enhancement | Interface 3 (Hook Target), Integration 5 | UF-6 Dashboard Hook Analysis | existing-page:Dashboard overlay | 1.3, 2.6, 3.5 |
| P2-1 Turn Efficiency Analysis | — | — | — | — |
| P2-2 Repeat Operation Detection | — | — | — | — |
| P2-3 Thinking Chain Visualization | — | — | — | — |
| P2-4 Cost & Success Rate | — | — | — | — |
