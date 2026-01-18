# Track Specification: Replace Basic Auth with OAuth Flow

## 1. 概述 (Overview)
由於 GitHub 已經棄用 Basic Auth，本專案需要全面改用 OAuth 2.0 Device Flow 或 Personal Access Token (PAT)。目前 `slctl init` 已實作 Device Flow，但 `slctl plugin install` 與 `slctl plugin search` 等其他功能仍在使用舊的 Basic Auth 邏輯或呼叫已標記為棄用的函式，導致功能失效。

## 2. 功能需求 (Functional Requirements)
- **全面支援 OAuth Flow**：
    - 更新 `pkg/github/token/token.go` 中的 `EnsureScopes`，使其不再回傳「Basic Auth is deprecated」錯誤，而是改為檢查當前的 OAuth Token/PAT 權限（如果可行）或僅驗證 Token 格式。
    - 確保 `slctl plugin install` 在安裝完成後能正確驗證權限。
    - 更新 `slctl plugin search`，使其在向 GitHub 請求資料時使用新的 OAuth Flow 驗證。
- **移除/調整舊的 Auth 邏輯**：
    - 搜尋程式碼庫中所有涉及 `Username`、`Password` 的 Basic Auth 結構或旗標。
    - 移除或將其替換為對 `Token` 的支援。
    - 如果某些 CLI 參數仍保留（為了相容性），應在執行時提示已棄用並導向 `slctl init`。

## 3. 非功能需求 (Non-Functional Requirements)
- **錯誤訊息友善化**：當 Token 失效或權限不足時，提示使用者重新執行 `slctl init`。
- **向下相容性**：如果使用者環境中已有 `$SL_TOKEN` 或 `configs.yaml` 中的 `token`，應優先使用之。

## 4. 驗收標準 (Acceptance Criteria)
- [ ] `slctl plugin install github.com/softleader/slctl-contacts` 能夠成功執行，不再出現 `basic Auth is deprecated` 錯誤。
- [ ] `slctl plugin search` 能成功列出 GitHub 上的 Plugins。
- [ ] 全域搜尋專案，不再有任何主動觸發 Basic Auth (username/password) 的邏輯。
- [ ] 所有現有測試通過，並針對新的 OAuth 驗證邏輯新增（或更新）單元測試。

## 5. 範圍外 (Out of Scope)
- 修正 `slctl contacts` 本身的 404 Bug（該問題已知且屬於該 Plugin 範疇）。
