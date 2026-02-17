# Loop Mitigation & Self-Healing Protocols

When the agent detects a "Stuck" state, it must execute the following strategies in order.

## Phase 1: Logic Shift (Low Cost)
**Trigger:** 3x Duplicate Outputs or 5x Tool Stutter.
1. **Linear Temp Escalation:** Increase LLM temperature by **+0.3** to force non-deterministic pathing [Source: Ralph v2.4].
2. **Semantic Rephrasing:** Rewrite the previous "Thought" block using entirely different vocabulary.
3. **Tool Cooling:** If the loop involves a specific tool (e.g., `grep`), ban that tool for the next 3 turns. Force usage of an alternative (e.g., `find` or python script) [Source: Agent Loop].

## Phase 2: Structural Reset (Medium Cost)
**Trigger:** Phase 1 failed.
1. **Context Pruning:** Delete the last 5 user/assistant message pairs from context to remove the "bias of failure." Retain only the System Prompt and `plan.md`.
2. **The "Fresh Shell" Pattern:** Restart the shell subprocess to clear any hanging env vars or zombie processes [Source: Loki Mode].

## Phase 3: The "Nuclear" Option (High Cost)
**Trigger:** Phase 2 failed.
1. **Git Branch Reset:** Hard reset the current working branch to `main` or the last known good commit.
   - Command: `git reset --hard origin/main`
   - Rationale: Codebase state is likely polluted/broken beyond repair [Source: Kinetic Conductor].
2. **Mark Task Failed:** Mark the current task in `plan.md` as `[FAILED]` and skip to the next independent task (if applicable).

## Phase 4: Termination
**Trigger:** Phase 3 failed.
1. **Alert User:** Output `CRITICAL_FAILURE: UNABLE_TO_RECOVER`.
2. **Halt:** Stop execution.
