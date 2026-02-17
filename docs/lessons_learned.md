# Lessons Learned: Unknown command "mcp" for "ralph2"

## Summary
The system failed to recognize the `mcp` subcommand despite the code being present. This was caused by a failed build process that silently left an outdated binary in place, coupled with incorrect third-party library implementation and unused imports.

## Incident Details

### Symptoms
*   `Error: unknown command "mcp" for "ralph2"` when running `./ralph2.exe mcp`.
*   The help output (`./ralph2.exe --help`) did not list the `mcp` command.
*   Compilation logs showed the build failed, but subsequent execution commands were still triggered using the existing (old) binary.

### Root Causes
*   **Failed Build Process**: The `go build` command failed due to syntax and logic errors, preventing the creation of a new binary.
*   **Incorrect API Usage**: The `github.com/metoro-io/mcp-golang` library was initialized without a required transport argument (`NewServer()` vs `NewServer(transport)`).
*   **Code Rot/Lint Errors**: Unused imports (e.g., `"os"`) in `internal/service/orchestrator.go` caused the Go compiler to abort the build (Go treats unused imports as errors).
*   **Lack of Build Verification**: The workflow did not check if the binary was successfully replaced before attempting to run it.
*   **Module Resolution Ping-Pong**: Multiple attempts to `go get` a library failed or resolved to different versions because of conflicting command syntax (`/server` subpackage vs root package) and inconsistent module cache states.
*   **Environment Path Fragility**: Relying on `$env:PATH` modifications in a single command chain can lead to "command not found" errors if a subsequent step expects the tool (like `go`) to be globally available without the local prefix.

### Containment
*   **Manual Clean**: Deleted the old `ralph2.exe` to ensure that any failed build would result in a "file not found" error rather than running an outdated version.
*   **Dependency Audit**: Verified the correct function signatures in the `mcp-golang` SDK using `go doc`.

### Countermeasures (AI Guidelines)
*   **Verify Build Success**: ALWAYS check for a successful exit code from the compiler before running the resulting binary.
*   **Immediate Linting**: Proactively remove unused imports and variables. Go is strict; one unused import breaks the entire build.
*   **Alias Conflicting Namespaces**: When a library name (e.g., `mcp`) conflicts with a local variable or package name, use an alias (e.g., `mcp_sdk "github.com/.../mcp"`) to prevent "variable is not a type" errors.
*   **Read the SDK Patterns**: Check function definitions (using `go doc` or `view_file`) before assuming a parameter-less constructor exists for third-party libraries.
*   **Clean Build**: Add a `rm -f binary_name` (or Windows equivalent) before building to guarantee that you are testing the latest code.
*   **Standardize Path Access**: Use absolute paths or a consistent `PATH` setup for build tools (`go`, `npm`) to avoid "not recognized" errors across different shell environments.
*   **Audit Module Imports**: Before implementing a third-party library, use `go list -m all` or `go doc` to confirm the exact package path and available constructors.
*   **Check Before Prepending PATH**: Run `where.exe <tool>` or `<tool> --version` FIRST to verify whether a tool is already available. Don't blindly prepend `$env:PATH` on every command — it wastes tokens, clutters output, and masks real issues.
