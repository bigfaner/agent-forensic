---
id: "3"
title: "Boundary & Layout Tests"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["2"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 3: Boundary & Layout Tests

## Description

Test the TUI's resilience to edge cases: terminal resize, empty/error states, minimum terminal size, and i18n rendering across locales. These tests ensure the UI degrades gracefully and renders correctly in all conditions.

## Reference Files
- `docs/proposals/tui-e2e-testing/proposal.md` — Source proposal
- `tests/e2e_go/helpers.go` — Infrastructure (Task 1)
- `internal/model/app.go` — Layout calculation, size warning
- `internal/model/statusbar.go` — Responsive hints (60/80/100 col breakpoints)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `tests/e2e_go/boundary_test.go` | Boundary & layout tests |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] **Minimum size test**: resize to 80x24 → view renders without crash, shows main layout
- [ ] **Below minimum test**: resize to 60x15 → view shows yellow size warning message
- [ ] **Resize adaptation test**: resize from 120x40 to 80x24 → panels recalculate widths, status bar truncates hints
- [ ] **Wide terminal test**: resize to 200x50 → layout uses full width, all panels visible
- [ ] **Empty session list test**: load model with no sessions → sessions panel shows empty state message
- [ ] **Error state test**: load model with invalid session data → error state displayed
- [ ] **No-anomaly diagnosis test**: open diagnosis on entry without anomalies → shows "no anomalies" message
- [ ] **Status bar responsive test**: at 60 cols shows basic hints, at 80 cols adds more, at 100 cols shows full hints + monitoring indicator
- [ ] **i18n layout test**: same resize scenarios run in both `zh` and `en`, verifying both locales render without overflow or truncation
- [ ] Total: 5+ test functions covering all scenarios

## Implementation Notes

1. **Resize testing**: Send `tea.WindowSizeMsg{Width: w, Height: h}` to the model. The AppModel recalculates panel widths: sessions = width/4 (min 25), callTree+detail share remaining width, callTree gets 67% of content height.

2. **Size warning**: AppModel shows a yellow warning when terminal is below 80x24. Check for the warning text in the view output.

3. **Status bar breakpoints**: StatusBarModel has 3 hint priority levels:
   - >=60 cols: basic navigation hints
   - >=80 cols: adds diagnosis and replay hints
   - >=100 cols: adds session/call shortcuts + monitoring indicator

4. **Empty states**: Create AppModel, don't call SetSessions(). Sessions panel should show empty state (localized message).

5. **No-anomaly diagnosis**: Navigate to a normal (non-anomaly) entry, press `d`. The diagnosis modal should show "no anomalies found" state.

6. **i18n consistency**: Run the same boundary test with both locales to catch layout issues specific to Chinese character widths (CJK characters are typically 2 columns wide).
