# Implementation Plan - Upgrade Go Version to 1.25.6

## Phase 1: Core Configuration
- [x] Task: 更新 `go.mod` 中的 Go 版本至 1.25.6 2deaea3
- [ ] Task: 更新 `conductor/tech-stack.md` 中的 Go 版本資訊
- [ ] Task: 執行 `go mod tidy` 驗證相依性並同步環境
- [ ] Task: Conductor - User Manual Verification 'Core Configuration' (Protocol in workflow.md)

## Phase 2: Infrastructure & CI/CD
- [ ] Task: 更新 GitHub Actions 工作流 (`.github/workflows/*.yml`) 以使用 Go 1.25.6
- [ ] Task: 檢查並更新 `.goreleaser.yml` 中的 Go 版本設定
- [ ] Task: 檢查 `Makefile` 是否有硬編碼的 Go 版本檢查或限制
- [ ] Task: Conductor - User Manual Verification 'Infrastructure & CI/CD' (Protocol in workflow.md)

## Phase 3: Documentation & Final Verification
- [ ] Task: 更新 `docs/` 目錄下提及 Go 版本的文件（如 `docs/Compiling-Source.md`）
- [ ] Task: 執行完整測試套件 (`go test ./...`) 確保功能正常
- [ ] Task: 執行專案編譯 (`make build`) 驗證產出物
- [ ] Task: Conductor - User Manual Verification 'Documentation & Final Verification' (Protocol in workflow.md)
