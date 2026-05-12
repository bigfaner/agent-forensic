---
created: 2026-05-12
source: prd/prd-ui-functions.md
status: Draft
---

# UI Design: Deep Drill Analytics

## Design System

Inherits the established TUI design system from the base agent-forensic feature (Vercel-inspired, monochrome dark theme). All color tokens, typography rules, layout grid conventions, and keyboard interaction patterns are consistent with the existing design documented at `docs/features/agent-forensic/ui/ui-design.md`.

### Extended Color Tokens

| Role | ANSI Color | Usage |
|------|-----------|-------|
| File Read | Bright Green (#55FF55) | Read operation bars and labels |
| File Edit | Bright Red (#FF5555) | Edit/Write operation bars and labels |
| Hook Timeline | Bright Magenta (#FF55FF) | Hook timeline connectors and markers |
| Strategy Change | Bright Yellow (#FFFF55) | Thinking chain strategy change markers |

### Extended Key Bindings

| Key | Context | Action |
|-----|---------|--------|
| `a` | Call Tree, SubAgent node selected | Open SubAgent full-screen overlay |
| `a` | Call Tree, non-SubAgent node | No-op |

### Extended Status Bar

Dashboard Active state key hints updated to:

```
s:back  1:session  j/k:nav  Tab:panel  Esc:back  a:agent  m:mon 监听:{状态}  q:quit
```

`a:agent` shown only when in Dashboard and current session has subagents (priority 3, >= 100 cols).

### Inherited Components

The following components from the base feature design (`docs/features/agent-forensic/ui/ui-design.md`) are inherited without modification:

- **UF-5 (Diagnosis Summary)**: PRD Navigation Architecture entry #5 (`d` key) opens the existing Diagnosis overlay defined in the base feature. This component is not redefined here because Phase 1 (Deep Drill Analytics) does not extend it. Phase 2 (Repeat Operation Detection) will add anomaly subtypes and extend the evidence block layout.
- **UF-6 (Status Bar)**: The base Status Bar component is inherited; only the key hints are extended via the Extended Status Bar table above.

### Terminal Resize Behavior

All overlays and panels re-render on SIGWINCH using the same responsive strategy as the base design:

- **Overlays (UF-2)**: Overlay dimensions recalculate as 80% x 90% of the new terminal size. Content reflows within the new dimensions using the section height allocation rules below. If terminal shrinks below 80x24, the base design's full-screen warning is shown.
- **Panels (UF-1, UF-3, UF-4, UF-5, UF-6)**: Text truncation and bar chart widths recalculate based on the new panel width. File paths re-truncate; bar chart `█` lengths rescale proportionally.
- **Minimum terminal size**: 80x24 (inherited from base design). Below this, no component content is rendered.

### Emoji and ASCII Fallback

The base feature design uses emoji indicators (`🟡`, `🔴`, `📦`, `⚠`) as the primary rendering. When the terminal does not support Unicode emoji (detected via `TERM` environment or explicit `--ascii` flag), the following ASCII replacements apply:

| Emoji | ASCII Fallback | Context |
|-------|---------------|---------|
| `📦` | `[A]` | SubAgent node marker (UF-1) |
| `⏳` | `...` | Loading state suffix (UF-1) |
| `⚠` | `!` | Error state suffix (UF-1) |
| `🟡` | `[S]` | Slow anomaly marker (base UF-5) |
| `🔴` | `[U]` | Unauthorized anomaly marker (base UF-5) |

Detection: check `TERM` for known limited terminals (e.g., `linux`, `dumb`) or honor `--ascii` CLI flag. Default is emoji rendering.

---

## Component: SubAgent Inline Expand (UF-1)

### Placement

- **Mode**: existing-page
- **Target**: Call Tree Panel
- **Position**: Below SubAgent parent node, indented 2 additional spaces (depth 2)

### Layout Structure

```
│ ● Turn 2 (5.1s)
│   ├─ Read src/main.go (0.3s)
│   ├─ SubAgent ×3 (4.8s) 📦
│   │  ├─ Read internal/model/app.go (0.2s)
│   │  ├─ Edit internal/model/app.go (1.5s)
│   │  └─ Bash go test ./... (2.8s)
│   └─ Write config.yaml (0.2s)
```

- SubAgent children indented 2 spaces deeper than parent (depth 2 vs depth 1)
- Tree connectors follow same `├─ └─` pattern at depth 2
- Children sorted by JSONL appearance order
- Loading state: SubAgent line shows `📦 ⏳` suffix while parsing (ASCII: `[A] ...`)
- Error state: SubAgent line shows `📦 ⚠` suffix, children hidden (ASCII: `[A] !`)
- Overflow: >50 children → last visible line shows `│  ... +N more` in text-secondary

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Collapsed | `├─ SubAgent ×3 (4.8s) 📦` (ASCII: `[A]`) | Default; Enter to expand |
| Loading | `├─ SubAgent ×3 (4.8s) 📦 ⏳` (ASCII: `[A] ...`) | Parsing subagents/ JSONL; async, non-blocking |
| Expanded | Children visible at depth 2 | Enter to collapse; children navigable |
| Error | `├─ SubAgent ×3 (4.8s) 📦 ⚠` (ASCII: `[A] !`) | JSONL missing or corrupt; stays collapsed |
| Overflow | `│  ... +5 more` as last child | >50 children; text-secondary color |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `Enter` on collapsed SubAgent | Expand (load JSONL) | Loading indicator → children appear or error |
| `Enter` on expanded SubAgent | Collapse | Children hidden |
| `j`/`k` over children | Navigate child nodes | Same highlight as depth-1 nodes |
| `Enter` on child node | Show detail in Detail Panel | Detail panel updates |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Child tool name | SubAgent session tool_use.name | String |
| Child duration | SubAgent session tool_use duration | "(Xs)" |
| Child count | SubAgent session tool_use count | "×N" |
| Loading indicator | Parse state | "⏳" (ASCII: "...") |
| Error indicator | Parse error | "⚠" (ASCII: "!") |

---

## Component: SubAgent Full-Screen Overlay (UF-2)

### Placement

- **Mode**: new-page
- **Target**: SubAgent Analysis Overlay
- **Position**: Centered overlay, 80% width × 90% height (screen dimensions), with 1-cell border; rendered above Dashboard/Call Tree content which dims to text-secondary

### Layout Structure

```
     ┌─ SubAgent: internal/model/app.go refactor ─── 4 tools, 12.3s ────────┐
     │                                                                       │
     │  Tool Statistics                                                      │
     │  ────────────────                                                     │
     │  Read    ████████████  12      Edit    ██████  5                      │
     │  Bash    ██████████   10      Write   ███  3                          │
     │                                                                       │
     │  File Operations (top 20)                                             │
     │  ─────────────────────                                                │
     │  internal/model/app.go    Read ×5  Edit ×3  ████████████  8          │
     │  cmd/root.go              Read ×3  Edit ×1  ██████  4                │
     │  internal/i18n/i18n.go    Read ×2            ████  2                  │
     │                                                                       │
     │  Duration Distribution                                                │
     │  ───────────────────                                                  │
     │  Bash     ████████████████████  8.2s  (67%)                           │
     │  Edit     ████████             3.1s  (25%)                            │
     │  Read     ████                 1.0s   (8%)                            │
     │                                                                       │
     │  Esc:close  j/k:scroll  Tab:sections                                 │
     └───────────────────────────────────────────────────────────────────────┘
```

- Title: "SubAgent: {title} — {N} tools, {duration}" in bold white
- Three sections separated by `────` dividers in text-secondary
- **Tool Statistics**: Horizontal bar chart (same pattern as Dashboard tool calls), sorted by count descending
- **File Operations**: Per-file rows with Read (green) and Edit (red) counts, total operations bar, sorted by total descending, max 20 rows
- **Duration Distribution**: Horizontal bar chart with time and percentage, sorted by duration descending

**Section height allocation:**

Available content height = overlay height (90% of terminal) - title (1 line) - footer (1 line) - 3 dividers (3 lines). The content area is split among the three sections:

| Section | Allocation | Rounding |
|---------|-----------|----------|
| Tool Statistics | 25% of content lines | Rounded up |
| File Operations | 50% of content lines | Rounded down |
| Duration Distribution | 25% of content lines | Takes remainder |

Example: 36-row terminal, overlay = 32 rows, content = 27 lines. Tool Statistics = 7 lines, File Operations = 13 lines, Duration Distribution = 7 lines. Each section scrolls independently via j/k when content exceeds its allocation.

**Tab interaction:** Tab cycles cursor between the three section headers (Tool Statistics -> File Operations -> Duration Distribution). The focused section header renders in cyan; unfocused headers remain in bold white. j/k scrolls content only within the focused section. On overlay open, focus defaults to Tool Statistics.

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Loading | "Loading subagent data..." centered | Async JSONL parse |
| Populated | Three-section layout | Scroll if content exceeds height |
| Empty | "No data" centered in text-secondary | SubAgent JSONL has 0 tool calls |
| Error | "Failed to load: {message}" in bright red | JSONL parse failure |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `Esc` | Close overlay, return to Call Tree | Cursor returns to SubAgent parent node |
| `j`/`k` | Scroll content | Virtual scroll within overlay |
| `Tab` | Cycle section focus | Focused section header turns cyan; j/k scrolls within focused section only |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Title | SubAgent session title + stats | "{title} — {N} tools, {duration}" |
| Tool bars | Tool call counts | `█` × proportion, count right-aligned |
| File rows | File operations | Path (40 chars) + Read ×N + Edit ×M + bar + total |
| Duration bars | Tool durations | `█` × proportion + time + percentage |

---

## Component: Turn Overview File Operations (UF-3)

### Placement

- **Mode**: existing-page
- **Target**: Detail Panel (Turn Overview mode)
- **Position**: Below existing "tools: N calls" block, before anomaly summary

### Layout Structure

```
│  tools: 5 calls, 12.3s
│    Read           ×3  2.1s
│    Edit           ×1  3.5s
│    Bash           ×1  6.7s
│  files:
│    internal/model/app.go  R×2  E×1
│    cmd/root.go            R×1
```

- Section label "files:" in bright cyan (same as "tools:" label)
- File rows: path (truncated to panel width, `...filename` format) + `R×N` (green) + `E×N` (red)
- Sorted by operation count descending
- Max 20 rows; overflow shows `+N more` in text-secondary
- Hidden when turn has no Read/Write/Edit calls

### States

| State | Visual | Behavior |
|-------|--------|----------|
| No file ops | Section hidden | "files:" label not rendered |
| Has file ops | File list | Scroll with parent Detail panel |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| None specific | Display-only | Inherits Detail panel scroll |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| File path | Read/Write/Edit input.file_path | Truncated to panel width |
| Read count | Read calls for this file | "R×N" in bright green |
| Edit count | Write/Edit calls for this file | "E×N" in bright red |

---

## Component: SubAgent Statistics in Detail (UF-4)

### Placement

- **Mode**: existing-page
- **Target**: Detail Panel
- **Position**: Replaces tool detail content when SubAgent child node selected

### Layout Structure

```
│  subagent stats:
│    tools: 8 calls, 15.2s
│      Read           ×3  2.1s
│      Edit           ×2  6.5s
│      Bash           ×2  5.8s
│      Write          ×1  0.8s
│    files:
│      internal/model/app.go  R×2  E×2
│      cmd/root.go            R×1  E×1
│    duration: avg 1.9s, peak Bash go test (5.2s)
```

- Section label "subagent stats:" in bright cyan
- "tools:" sub-block: same format as Turn Overview tool stats
- "files:" sub-block: same format as Turn Overview file ops
- "duration:" sub-block: average + peak call
- Tab toggles between this stats view and individual tool detail view

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Stats view | SubAgent statistics | Default when selecting SubAgent child |
| Tool detail | Individual tool input/output | After Tab toggle |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `Tab` | Toggle stats ↔ tool detail | View switches; title updates |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Tool stats | SubAgent session tool aggregation | Same as Turn Overview |
| File list | SubAgent file operations | Same as Turn Overview |
| Duration stats | Average and peak | "avg Xs, peak {tool} ({duration})" |

---

## Component: Dashboard File Operations Panel (UF-5)

### Placement

- **Mode**: existing-page
- **Target**: Dashboard overlay
- **Position**: Below existing Custom Tools block, new section

### Layout Structure

```
│  File Operations (top 20)
│  ────────────────────────────────────────────────
│  internal/model/app.go  ██████████████████  R×5  E×3  8
│  cmd/root.go            ██████████        R×3  E×1  4
│  internal/i18n/i18n.go  ██████            R×2        2
│  config/production.yml  ████              E×1        1
```

- Section header: "File Operations (top 20)" in bold white
- Divider: `────` in text-secondary
- Each row: path (40 chars, truncated with `...` prefix) + horizontal bar + `R×N` in bright green + `E×N` in bright red + total count
- Bar: `█` characters proportional to total ops, max 20 chars
- Sorted by total operations descending
- Hidden when session has no Read/Write/Edit calls

### States

| State | Visual | Behavior |
|-------|--------|----------|
| No file ops | Section hidden | Entire block not rendered |
| Has file ops | Bar chart with file paths | Static, updates on session change |
| >20 files | Top 20 shown | "+N more" in text-secondary at bottom |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `j`/`k` | Scroll Dashboard content vertically | Panel scrolls; new rows appear at viewport edge. Uses same virtual-scroll mechanism as base Dashboard (see `docs/features/agent-forensic/ui/ui-design.md` UF-4 Dashboard View) |
| `Tab` | Cycle focus to next Dashboard section | File Operations section header highlighted in cyan when focused |
| `s` / `Esc` | Return to Call Tree view | Dashboard closes (same as base Dashboard interaction) |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| File path | Aggregated file_path | Truncated 40 chars |
| Read bar | Read count | `█` × proportion, bright green |
| Edit bar | Edit count | `█` × proportion, bright red |
| Total count | Read + Edit count | Integer, right-aligned |

---

## Component: Dashboard Hook Analysis Panel (UF-6)

### Placement

- **Mode**: existing-page
- **Target**: Dashboard overlay
- **Position**: Replaces existing Hook column in Custom Tools block + new timeline section below

### Layout Structure

```
│  Hook Statistics
│  ────────────────────────────────────────────────
│  PreToolUse::Bash       ×12
│  PreToolUse::Edit       ×5
│  PostToolUse::Bash      ×8
│  PostToolUse::Edit      ×3
│  Stop                   ×2
│
│  Hook Timeline (by Turn)
│  ────────────────────────────────────────────────
│  Legend: ●PreToolUse(green) ●PostToolUse(cyan) ●Stop(yellow) ●user-prompt(magenta)
│  T1  ●PreToolUse::Bash ●PreToolUse::Bash ●PostToolUse::Bash ●PreToolUse::Edit ●PostToolUse::Edit
│  T2  ●PreToolUse::Bash ●PreToolUse::Bash ●PostToolUse::Bash ●PostToolUse::Bash
│  T3  ●PreToolUse::Edit ●PostToolUse::Edit ●Stop
```

**Hook Statistics section:**
- Section header: "Hook Statistics" in bold white
- Each row: `HookType::Target` in white + `×N` right-aligned
- Target extraction failure: show only `HookType` (no `::` suffix)
- Sorted by count descending

**Hook Timeline section:**
- Section header: "Hook Timeline (by Turn)" in bold white
- **Legend row**: displayed immediately below the section header, before the first Turn row. Shows one entry per marker type with its color and a text label: `● PreToolUse  ● PostToolUse  ● Stop  ● user-prompt-submit` (markers rendered in their respective colors: bright green, bright cyan, bright yellow, bright magenta; labels in text-secondary)
- Each row: Turn label (`T{N}`) in text-secondary + `●` markers for each hook trigger, followed by the full `HookType::Target` label (matching the Statistics section naming)
- Marker format: `●HookType::Target` using the same full names as Hook Statistics (e.g., `●PreToolUse::Bash`, `●PostToolUse::Bash`), not abbreviated
- Marker color by hook type: PreToolUse = bright green, PostToolUse = bright cyan, Stop = bright yellow, user-prompt-submit = bright magenta
- Turn label: 3 chars wide, right-aligned
- Markers separated by spaces, max 30 markers per line; overflow wraps to next line with continuation indent

### States

| State | Visual | Behavior |
|-------|--------|----------|
| No hooks | Both sections hidden | Entire block not rendered |
| Has hooks | Statistics + Timeline | Static, updates on session change |
| Target extraction failed | HookType without `::` | Fallback display |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `j`/`k` | Scroll Dashboard content vertically | Panel scrolls; new rows appear at viewport edge. Uses same virtual-scroll mechanism as base Dashboard (see `docs/features/agent-forensic/ui/ui-design.md` UF-4 Dashboard View) |
| `Tab` | Cycle focus to next Dashboard section | Hook Analysis section header highlighted in cyan when focused |
| `s` / `Esc` | Return to Call Tree view | Dashboard closes (same as base Dashboard interaction) |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Hook label | HookType + TargetCommand | "HookType::Target" or "HookType" |
| Hook count | Trigger count | "×N" |
| Timeline Turn label | Turn index | "T{N}" |
| Timeline markers | Per-Turn hook triggers | `●HookType::Target` with type-specific color |
