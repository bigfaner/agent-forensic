---
id: "4"
title: "Real-time Monitoring Flow Test"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["3"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 4: Real-time Monitoring Flow Test

## Description

Test the real-time monitoring pipeline: watcher events trigger incremental parsing, which adds entries to the call tree with flash indicators that expire after 3 seconds. This is the only flow that requires `tea.TestProgram` for its `tea.Cmd` support.

## Reference Files
- `docs/proposals/tui-e2e-testing/proposal.md` — Source proposal
- `tests/e2e_go/helpers.go` — Infrastructure (Task 1)
- `internal/model/calltree.go` — AddEntry, flash mechanism, flashTickMsg
- `internal/model/app.go` — WatcherEventMsg handling

## Affected Files

### Create
| File | Description |
|------|-------------|
| `tests/e2e_go/monitoring_test.go` | Real-time monitoring flow tests |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] **AddEntry flash test**: send WatcherEventMsg → call tree shows new entry with `[NEW]` flash indicator
- [ ] **Flash expiry test**: after flash, advance time past 3s → `[NEW]` indicator removed from view
- [ ] **Sequential events test**: send multiple WatcherEventMsgs → all entries appear, each with flash indicator
- [ ] **Auto-expand test**: new entry in collapsed turn → turn auto-expands to show new entry
- [ ] **Monitoring toggle test**: press `m` → status bar shows monitoring enabled; press again → disabled
- [ ] **Integration journey test**: enable monitoring → receive event → view shows flash → navigate to entry → view detail → wait for flash expiry → flash gone
- [ ] Total: 5+ test functions

## Implementation Notes

1. **WatcherEventMsg**: This custom message carries `{FilePath, Lines []string}`. The AppModel.Update() handler parses lines and calls CallTreeModel.AddEntry(). Since this is a custom message (not a tea.Cmd), it can be sent via direct `Update()` call.

2. **Flash mechanism**: CallTreeModel.AddEntry() sets a flash expiry time (3s from now). A `flashTickMsg` cleans up expired flashes. To test expiry, manually send a `flashTickMsg` after constructing a model with expired flash times.

3. **tea.TestProgram usage**: For the integration journey test, use `tea.TestProgram` to handle the `flashTickMsg` Cmd chain. The program sends the initial message, the model returns a tick Cmd, and the program executes it.

4. **Flash timing**: Don't rely on real time in tests. Construct entries with explicit flash expiry times:
   - New entry: set flash expiry to `time.Now().Add(3 * time.Second)`
   - To test expiry: send update after 4+ seconds (or mock time)

5. **Auto-expand**: When a new entry is added to a collapsed turn, the turn should auto-expand. Verify by checking that the turn icon changes from `●` (collapsed) to `▼` (expanded) and child entries are visible.

6. **Monitoring indicator**: StatusBarModel shows `监听:开`/`Watch:ON` (green) when monitoring is active. Verify this appears in the view after pressing `m`.
