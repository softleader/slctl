# 實作計畫：修復插件搜尋與 TestFetchOnline 失敗

## 第一階段：診斷與環境準備
- [x] Task: 建立失敗測試案例並確認錯誤訊息 (Red Phase) [6286cf0]
    - [ ] 執行 `go test -v ./pkg/plugin/ -run TestFetchOnline`
    - [ ] 記錄並分析目前的 API 回傳結果
- [ ] Task: 驗證 GitHub Search API 呼叫參數
    - [ ] 在 `fetchOnline` 中增加 Debug Log，輸出 Query 字串
    - [ ] 確保 `org:softleader topic:slctl-plugin` 符合 GitHub API 規範

## 第二階段：修正搜尋邏輯 (Green Phase)
- [ ] Task: 修正 `fetchOnline` 搜尋邏輯
    - [ ] 調整 `query` 字串格式（例如將 `+` 改為空格，讓 SDK 處理編碼）
    - [ ] 確保 API Client 獲得正確的 Token 授權
- [ ] Task: 優化搜尋結果處理
    - [ ] 增加對空結果的進一步診斷（例如檢查回應的狀態碼與內容）
- [ ] Task: 執行測試並驗證 (Green Phase)
    - [ ] 確保 `TestFetchOnline` 成功通過且結果數量符合預期
- [ ] Task: Conductor - User Manual Verification '修正搜尋邏輯' (Protocol in workflow.md)

## 第三階段：品質保證與清理
- [ ] Task: 程式碼清理與型別檢查
    - [ ] 執行 `go fmt` 與 `go vet`
- [ ] Task: 最終驗證
    - [ ] 手動執行 `slctl plugin search` 確保功能正常
- [ ] Task: Conductor - User Manual Verification '品質保證與清理' (Protocol in workflow.md)
