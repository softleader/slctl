# 規格書: 升級 go-github 並現代化登入流程

## 1. 概觀
本 Track 旨在將專案相依的 `go-github` 函式庫從 v28 升級至最新穩定版本 (v69+)。鑑於 GitHub 已棄用 Basic Authentication (帳號/密碼)，且舊版 SDK 的相關支援可能已移除，本 Track 將重構 `slctl init` 的初始化流程：移除過時的帳號密碼輸入方式，轉而支援直接輸入 Personal Access Token (PAT) 以及新增 OAuth 2.0 Device Flow 以提升使用者體驗。

## 2. 功能需求

### 2.1 升級依賴
- [ ] 將 `go.mod` 中的 `github.com/google/go-github` 更新至最新主要版本 (v69+)。
- [ ] 修正因升級導致的所有編譯錯誤 (Breaking Changes)。
- [ ] 確保專案中其他與 GitHub API 互動的部分 (如 `pkg/github`) 兼容新版 SDK。

### 2.2 重構初始化流程 (Modernize Init Flow)
- [ ] **移除 Basic Auth**: 
    - 刪除 `slctl init` 中提示輸入 Username/Password 的互動邏輯。
    - 刪除 `NewBasicAuthClient` 及相關的 OTP 處理邏輯。
    - 移除 CLI flag `-u/--username` 和 `-p/--password` (或標記為 Deprecated 並報錯)。
- [ ] **支援 PAT 輸入**:
    - 保留並優化現有的 `--token` 流程。
    - 若使用者未提供 Token 且不使用 Device Flow，則提示使用者輸入 Token。
- [ ] **實作 GitHub OAuth 2.0 Device Flow**:
    - 在 `slctl init` 新增選項或預設引導使用者使用 Device Flow 登入。
    - 實作流程：
        1. CLI 使用預設的 Client ID (需配置) 向 GitHub 請求 Device Code。
        2. 顯示 User Code 與驗證網址，提示使用者在瀏覽器開啟。
        3. CLI 輪詢 (Polling) 等待使用者授權。
        4. 取得 Access Token 並儲存。
    *註：需確認或建立一個 GitHub OAuth App 以取得 Client ID。若無現成 ID，則需提供使用者自行輸入 Client ID 的選項或說明。*

## 3. 非功能需求
- **安全性**: Token 必須安全地儲存在本地配置中 (維持現有机制作法)。
- **使用者體驗**: Device Flow 的提示必須清晰，包含足夠的指引。
- **兼容性**: 盡量維持 `--offline` 模式的功能不受影響。

## 4. 驗收標準
- [ ] `go list -m github.com/google/go-github/v69` (或更新版本) 顯示正確版本。
- [ ] `slctl init` 不再詢問帳號密碼。
- [ ] `slctl init --token <VALID_TOKEN>` 可成功驗證並儲存設定。
- [ ] `slctl init` (無參數) 可觸發 Device Flow，並在瀏覽器授權後成功取得 Token 並儲存。
- [ ] 所有現有單元測試通過。

## 5. 排除範圍 (Out of Scope)
- 實作 Web Flow (Authorization Code Grant)。
- 對 `slctl` 其他非 GitHub 相關功能的重大重構。
