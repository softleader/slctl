# Track 規格說明：修復插件搜尋失敗問題

## 概述
使用者回報 `slctl plugin search` 無法找到任何結果，且 `TestFetchOnline` 測試失敗（預期至少 3 個結果，實際得到 0 個）。此問題發生在所有環境中，且已知插件應從 GitHub `softleader` 組織中具有 `slctl-plugin` Topic 的 Repository 取得。

## 功能需求
1.  **修正搜尋邏輯**：確保 `fetchOnline` 能正確調用 GitHub Search API 並取得符合條件的 Repository。
2.  **優化錯誤處理**：若搜尋結果為空，應提供更清楚的日誌或除錯資訊（例如輸出實際發送的 Query 和 API 回應狀態）。
3.  **驗證測試**：修復 `pkg/plugin/repository_test.go` 中的 `TestFetchOnline`，使其能穩定通過。

## 驗收標準
1.  執行 `go test -v ./pkg/plugin/ -run TestFetchOnline` 應通過，且回傳結果數量需符合測試要求（目前設定為 >= 3）。
2.  手動執行 `slctl plugin search` 應能列出 `softleader` 組織下標記為 `slctl-plugin` 的插件。

## 非功能需求
- **效能**：搜尋回應時間應保持在 GitHub API 正常延遲範圍內。
- **穩定性**：應正確處理 GitHub API 的 Rate Limit 或網路問題。

## 超出範圍 (Out of Scope)
- 修改插件安裝 or 執行的邏輯。
- 支援 GitHub 以外的其他插件來源。
