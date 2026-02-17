# Loki Mode Continuity Log

## Current State
**Phase:** Core Development Loop
**Status:** Implementing Story 4.1: Windows Job Object Wrapper
**Active Agent:** Orchestrator (Loki)
**Branch:** story/4-1-windows-sandbox

## Immediate Plan (RARV)
1.  [x] **Foundation:** Story 2.3 complete (Autonomous Plan->Build loop).
2.  [x] **Branching:** `story/4-1-windows-sandbox` created.
3.  [ ] **Reason:** Need physical constraints on child processes for security.
4.  [ ] **Act:** Implement `SetInformationJobObject` to enforce memory limits and kill-on-close.
5.  [ ] **Act:** Create test to verify sandbox isolation.
6.  [ ] **Verify:** Test creates a process, closes the job, and process tree is wiped.

## Mistakes & Learnings
*   **Merge Overhead:** Resolved.
*   **Git Dirty Checks:** Need to ensure binaries like `ralph2.exe` are consistently ignored to pass the "dirty check" in `run.go`.
*   **Windows Constants:** `syscall` package in Go lacks some JobObject constants. Will manually define them using `unsafe` and `uintptr`.
