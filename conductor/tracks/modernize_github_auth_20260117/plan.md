# 實作計畫: 升級 go-github 並現代化登入流程

## Phase 1: 依賴與基礎建設 (Dependencies & Infrastructure) [checkpoint: fbd7e7f]
- [x] Task: 升級 go-github 版本 48bb7f8
    - [x] 執行 `go get github.com/google/go-github/v69` (或最新版)。
    - [x] 執行 `go mod tidy` 確保 `go.mod` 與 `go.sum` 正確更新。
    - [x] 修正因版本升級導致的編譯錯誤 (Compile Errors)。
    - [x] 執行現有測試 (Expected to fail)，確認受影響範圍。
- [x] Task: Conductor - User Manual Verification 'Dependencies & Infrastructure' (Protocol in workflow.md) fbd7e7f

## Phase 2: 重構 API Client (Refactor API Client) [checkpoint: 625957b]
- [x] Task: 重新設計 Client 建構邏輯 9ad51e2
    - [x] 建立新的測試 `pkg/github/client_test.go` 定義新版 Client 的預期行為。
    - [x] 移除 `pkg/github/client.go` 中的 `NewBasicAuthClient` 及其相關 Basic Auth 邏輯。
    - [x] 優化 `NewTokenClient` 或建立新的建構函式以支援直接傳入 Token 初始化的 Client。
    - [x] 確保現有的 `slctl` 其他指令仍能透過 Token 正確取得 Client。
- [x] Task: Conductor - User Manual Verification 'Refactor API Client' (Protocol in workflow.md) 625957b

## Phase 3: 實作裝置授權流程 (Implement Device Flow) [checkpoint: e628313]
- [x] Task: 實作 OAuth 2.0 Device Flow 邏輯 e9adcf5
    - [x] 建立 `pkg/github/auth.go` (或適當位置) 的測試檔案，定義 Device Flow 的介面與 mock。
    - [x] 實作請求 Device Code 的功能 (`POST https://github.com/login/device/code`)。
    - [x] 實作輪詢 Access Token 的功能 (`POST https://github.com/login/oauth/access_token`)，包含處理 `authorization_pending`, `slow_down` 等狀態。
    - [x] 整合 Client ID 配置 (預設或參數化)。
- [x] Task: Conductor - User Manual Verification 'Implement Device Flow' (Protocol in workflow.md) e628313

## Phase 4: 更新 CLI 初始化指令 (Update CLI Init Command) [checkpoint: 60e138c]
- [x] Task: 重構 `slctl init` 介面 b13a701
    - [x] 修改 `cmd/slctl/init.go`，移除用戶名/密碼的 Flags (`-u`, `-p`) 和相關變數。
    - [x] 更新 Help/Usage 訊息，反映新的登入選項。
- [x] Task: 整合新的登入流程 ae3c56f
    - [x] 實作邏輯：若無 `--token`，則啟動 Device Flow。
    - [x] 顯示清晰的 User Code 與 Verification URL 提示。
    - [x] 在成功取得 Token 後，呼叫現有的儲存邏輯 (Config Persistence)。
    - [x] 驗證流程：確保 Token 有效性並取得當前使用者資訊 (類似原有 `Welcome aboard %s!`)。
- [x] Task: Conductor - User Manual Verification 'Update CLI Init Command' (Protocol in workflow.md) 60e138c

## Phase 5: 驗證與清理 (Verification & Cleanup) [checkpoint: 0e02858]
- [x] Task: 全面測試與文件更新 0bcc952
    - [x] 執行所有單元測試確保沒有 Regression。
    - [x] 手動驗證流程：
        - `slctl init` (Device Flow)
        - `slctl init --token <PAT>`
        - `slctl init --offline`
    - [x] 更新使用者文件 (如果有) 說明新的登入方式。
- [x] Task: Conductor - User Manual Verification 'Verification & Cleanup' (Protocol in workflow.md) 0e02858
