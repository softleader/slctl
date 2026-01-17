# Implementation Plan - 自動開啟瀏覽器與複製驗證碼

## Phase 1: 依賴與基礎建設
- [ ] Task: 新增專案依賴 `github.com/atotto/clipboard`
    - [ ] 執行 `go get github.com/atotto/clipboard`
    - [ ] 執行 `go mod tidy` 更新 `go.mod` 與 `go.sum`

## Phase 2: 實作自動化功能 (TDD)
- [ ] Task: 重構 `cmd/slctl/init.go` 以支援測試 (Refactor for Testability)
    - [ ] 定義變數 `openBrowser` (預設呼叫 `open.Run`)
    - [ ] 定義變數 `writeToClipboard` (預設呼叫 `clipboard.WriteAll`)
    - [ ] 目的：允許在測試中替換這些函數以進行 Mock。
- [ ] Task: TDD - 實作「開啟瀏覽器」與「複製驗證碼」
    - [ ] Write Test (Red): 在 `cmd/slctl/init_test.go` 中新增測試案例，模擬 Device Flow 取得 `verification_uri` 與 `user_code` 後，驗證 `openBrowser` 與 `writeToClipboard` 是否被正確呼叫。
    - [ ] Implement (Green): 修改 `cmd/slctl/init.go` 中的 `run` 函數，整合自動化邏輯。
    - [ ] Refactor: 確保錯誤處理機制（如開啟失敗不崩潰，僅記錄 Log）。
    - [ ] Verify: 執行 `go test ./cmd/slctl/...` 確認測試通過。

## Phase 3: 驗證與交付
- [ ] Task: Conductor - User Manual Verification '功能驗證' (Protocol in workflow.md)
