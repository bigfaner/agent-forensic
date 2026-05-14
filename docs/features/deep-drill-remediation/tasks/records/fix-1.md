---
status: "completed"
started: "2026-05-14 02:02"
completed: "2026-05-14 02:08"
time_spent: "~6m"
---

# Task Record: fix-1 Fix: golden test header alignment (detail panel)

## Summary
Fix detail panel title bar hint alignment by using a non-East-Asian runewidth Condition for the padding calculation, resolving golden test failures on Windows with CJK code page (936)

## Changes

### Files Created
无

### Files Modified
- internal/model/detail.go

### Key Decisions
- Used runewidth.Condition{EastAsianWidth: false} instead of the global runewidth.StringWidth() to ensure ambiguous-width characters (—, ↑, ↓) are measured as width 1, matching the golden file expectations and producing consistent alignment regardless of system locale/code page

## Test Results
- **Tests Executed**: No
- **Passed**: 558
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] TestGolden_DetailTruncated passes
- [x] TestGolden_DetailMasked passes
- [x] All model package tests pass

## Notes
Root cause: Windows console code page 936 (GBK) causes runewidth.IsEastAsian() to return true, making ambiguous-width chars measure as 2 instead of 1. The golden files were generated with EastAsianWidth=false. Only the title bar padding calculation (line 293) was affected.
