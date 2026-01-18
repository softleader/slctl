# 實作計畫：以 OAuth Flow 取代 Basic Auth

## 1. 基礎架構更新 (支援 Token 中的 OAuth)
- [ ] Task: 撰寫 `token.EnsureScopes` 支援 OAuth Token 的失敗測試
- [ ] Task: 更新 `token.EnsureScopes`，透過 GitHub API 驗證 Scope 而非直接回傳錯誤
- [ ] Task: Conductor - User Manual Verification '基礎架構更新' (Protocol in workflow.md)

## 2. 功能模組適配 (插件相關指令)
- [ ] Task: 撰寫 `slctl plugin search` 使用 OAuth Flow 的失敗測試
- [ ] Task: 在 `slctl plugin search` 中實作 OAuth Flow
- [ ] Task: 驗證 `slctl plugin install` 與新版 `EnsureScopes` 的整合
- [ ] Task: Conductor - User Manual Verification '功能模組適配' (Protocol in workflow.md)

## 3. 清理舊邏輯 (移除遺留的 Basic Auth)
- [ ] Task: 識別並移除 CLI 旗標與結構中所有遺留的 `Username`/`Password` 引用
- [ ] Task: 更新文件與幫助訊息，反映僅支援 OAuth 驗證
- [ ] Task: 最終驗證所有與 GitHub 相關的指令
- [ ] Task: Conductor - User Manual Verification '清理舊邏輯' (Protocol in workflow.md)
