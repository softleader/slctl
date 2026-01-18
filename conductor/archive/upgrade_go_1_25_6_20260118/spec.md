# Track Specification: Upgrade Go Version to 1.25.6

## 1. Overview
本 Track 的目標是將 `slctl` 專案的 Go 語言版本從目前的 1.24+ 升級到 **1.25.6**。這包含更新開發環境設定、相依性宣告、CI/CD 自動化流程以及相關說明文件，並確保在升級後所有測試仍能正常通過。

## 2. Functional Requirements
- **Go 版本升級**：將專案中所有參考到 Go 版本的地方更新為 `1.25.6`。
- **環境一致性**：確保本機開發與 GitHub Actions CI/CD 環境使用相同的 Go 版本。
- **文件同步**：更新專案內部文件（如 `docs/`）中提及的 Go 版本資訊。

## 3. Technical Changes
- **`go.mod`**：更新 `go` 行之版本宣告。
- **GitHub Actions (`.github/workflows/*.yml`)**：更新 `actions/setup-go` 中指定的 `go-version`。
- **GoReleaser (`.goreleaser.yml`)**：如有指定 Go 版本，需同步更新。
- **Tech Stack Document (`conductor/tech-stack.md`)**：更新技術棧文件中的 Go 版本說明。
- **`docs/` 目錄**：檢查並更新 `docs/Compiling-Source.md` 等文件中提及的 Go 版本要求。
- **Makefile**：檢查是否有硬編碼的 Go 版本檢查。

## 4. Acceptance Criteria
- [ ] `go.mod` 中的 Go 版本更新為 `1.25.6`。
- [ ] 所有 GitHub Actions 工作流成功完成且無報錯。
- [ ] 執行 `go test ./...` 通過所有單元測試與整合測試。
- [ ] 執行 `make build`（或相關編譯指令）成功產出執行檔。
- [ ] `conductor/tech-stack.md` 已反映最新版本。
- [ ] `docs/` 下的所有相關文件皆已更新為正確的 Go 版本資訊。

## 5. Out of Scope
- 更新 Homebrew/Chocolatey 的外部發佈說明文件。
- 大規模重構程式碼以使用 Go 1.25 的新特性（僅專注於升級與相容性）。
