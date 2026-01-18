# Track Specification: Upgrade brew-tapper for Multi-Platform Support & Modernization

## 1. Overview
本 Track 旨在升級 `homebrew-tap` 儲存庫中的 `brew-tapper` 工具，使其能夠處理 Apple Silicon (arm64) 架構的 Formula 更新。同時，將專案的 Go 版本升級至 1.25.6 並更新相關依賴套件，以確保安全性與效能。這將使 `slctl` 的 Release Workflow 能夠自動化地同步更新 Homebrew Formula，並正確支援現代 macOS 硬體。

## 2. Goals
1.  **擴充 Regular Expression 邏輯**：修改 `brew-tapper` 核心邏輯，使其能辨識並替換 Formula 中 `arm64` 架構對應的 `url` 和 `sha256`。
2.  **專案現代化**：升級 Go Runtime 至 1.25.6，並更新 `go.mod` 中所有依賴至最新穩定版本。
3.  **保持向下相容**：確保現有的 `x86_64` (Intel) 和 `Linux` 更新邏輯不受影響。
4.  **驗證與發布**：透過單元測試與整合測試驗證修改後的邏輯，並確保新版 `brew-tapper` 能被 `slctl` 的 CI 流程正確呼叫。

## 3. Scope
### In Scope
- **功能實作**：
    - 修改 `brew-tapper/pkg/brew/formula.go` 以擴充 `Formula` struct 支援 arm64 欄位。
    - 修改 `brew-tapper/pkg/brew/formula_upgrade.go` 中的 Regex 和 `format` 函式，實作多平台替換邏輯。
- **維護工程**：
    - 更新 `go.mod` 將 Go 版本設為 `1.25.6`。
    - 執行 `go get -u ./...` 與 `go mod tidy` 升級依賴。
- **測試與驗證**：
    - 增加完整的單元測試 (`_test.go`) 覆蓋各種 Formula 格式情境。
    - 在本地進行 Dry Run 測試。
    - 協助使用者將編譯後的 binary 或呼叫邏輯整合進 `slctl` 的 workflow (提供指引或範例)。

### Out of Scope
- 重寫整個 `homebrew-tap` 的維護方式 (僅專注於升級現有的 `brew-tapper` 工具)。
- 修改 `slctl` 本身的 Go 程式碼 (僅涉及 Release Workflow 的設定調整引用)。

## 4. Technical Details
- **Target Repository**: `softleader/homebrew-tap` (利用 multi-workspace access)
- **Language**: Go 1.25.6
- **Key Files**:
    - `pkg/brew/formula.go`: 新增 `DarwinArm64Sha256` 等欄位。
    - `pkg/brew/formula_upgrade.go`: 更新 Regex 如 `darwinArm64Sha256Regexp`。
    - `go.mod`, `go.sum`: 依賴管理檔案。
    - `.github/workflows`: (若有) 確保 CI 環境使用 Go 1.25.6。
- **Integration**: `slctl` 的 GitHub Actions 將執行 `brew-tapper` binary，傳參數進行更新。

## 5. Acceptance Criteria
1.  **單元測試通過**：包含 arm64 屬性的 Formula 字串能被正確解析並替換 SHA256，且不破壞原有結構。
2.  **建置成功**：專案能在 Go 1.25.6 環境下成功編譯 (`go build`) 並通過測試 (`go test`)。
3.  **Dry Run 成功**：在不實際 commit 的情況下，工具能產出預期的 Formula 內容變更。
4.  **CI 整合驗證**：模擬一次 Release，確認 `homebrew-tap` 收到正確的 Formula 更新 commit。
