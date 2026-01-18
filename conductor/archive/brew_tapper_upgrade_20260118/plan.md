# Implementation Plan - Upgrade brew-tapper for Multi-Platform Support & Modernization

## Phase 1: Project Modernization (Go Upgrade)
- [ ] Task: Upgrade Go version to 1.25.6
    - [ ] Update `go.mod` to use Go 1.25.6
    - [ ] Update `Dockerfile` to use Go 1.25.6 base image (or appropriate builder image)
    - [ ] Update `Makefile` or CI configs if they hardcode Go versions
- [ ] Task: Upgrade Dependencies
    - [ ] Run `go get -u ./...` to upgrade all dependencies
    - [ ] Run `go mod tidy` to clean up
    - [ ] Verify project compiles with `go build ./...`
    - [ ] Run existing tests `go test ./...` to ensure no regressions from upgrades
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Project Modernization (Go Upgrade)' (Protocol in workflow.md)

## Phase 2: Implement Multi-Platform Support Logic
- [ ] Task: TDD - Prepare Test Cases for Multi-Platform Formula
    - [ ] Create/Update `pkg/brew/formula_upgrade_test.go`
    - [ ] Add test case: `TestUpgrade_MultiPlatform` containing a mock Formula content with `arm64` blocks
    - [ ] Define expected output structure (ensure it matches Homebrew DSL for multiple platforms)
    - [ ] Run tests and confirm failure (Red Phase)
- [ ] Task: Update Data Structures
    - [ ] Modify `pkg/brew/formula.go`: Add `DarwinArm64Sha256`, `LinuxArm64Sha256` (if needed) fields to `Formula` struct
- [ ] Task: Implement New Regex Logic
    - [ ] Modify `pkg/brew/formula_upgrade.go`
    - [ ] Add/Update Regex patterns to capture and replace `sha256` for specific OS/Arch conditions
    - [ ] Update `format()` function to apply these new replacements
- [ ] Task: Verify Implementation
    - [ ] Run tests and confirm pass (Green Phase)
    - [ ] Refactor regex logic if complex or duplicated
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Implement Multi-Platform Support Logic' (Protocol in workflow.md)

## Phase 3: Integration & Final Verification
- [ ] Task: CLI Analysis & CLI Verification (Dry Run)
    - [ ] Build the tool: `make build` (or `go build cmd/tapper`)
    - [ ] Run `tapper` locally against a dummy file or mock server to verify behavior
- [ ] Task: Documentation & Cleanup
    - [ ] Update `README.md` to reflect new capabilities and Go version requirement
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration & Final Verification' (Protocol in workflow.md)
