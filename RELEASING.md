# Releasing Ralph2

Ralph2 uses [GoReleaser](https://goreleaser.com/) and GitHub Actions for automated releases.

## How to Release

1.  **Ensure all changes are on `main` and passing CI.**
2.  **Tag the commit:**
    ```bash
    git tag -a v1.0.0 -m "Release v1.0.0"
    ```
3.  **Push the tag:**
    ```bash
    git push origin v1.0.0
    ```
4.  **GitHub Actions will detect the tag and trigger the release pipeline.**
5.  Check the "Releases" section on GitHub for the binaries and changelog.

## Manual Dry-Run

You can run GoReleaser locally to verify the build:
```bash
goreleaser release --snapshot --clean
```
(Requires `goreleaser` to be installed locally).
