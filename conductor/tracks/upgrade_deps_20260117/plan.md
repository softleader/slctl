# Implementation Plan - Upgrade Dependencies

這個計畫旨在升級 `slctl` 專案的依賴與開發環境。

## Phase 1: 環境準備與依賴分析
- [x] Task: 檢查當前開發環境與工具鏈可用性
    - [x] 確認 `go` 版本與環境變數
    - [x] 執行現有測試確保基準線正確
- [x] Task: 分析 `go.mod` 中的直接依賴，識別需要優先升級的套件
    - [x] 列表當前依賴版本
    - [x] 識別有安全性更新或重大更新的套件
- [~] Task: Conductor - User Manual Verification 'Phase 1: 環境準備與依賴分析' (Protocol in workflow.md)

## Phase 2: 核心套件升級 (Go & Cobra)
- [ ] Task: 升級 Go 版本
    - [ ] 修改 `go.mod` 中的 `go` 版本
    - [ ] 更新 `Makefile` 相關版本屬性
- [ ] Task: 升級 `spf13/cobra` 與 `spf13/pflag`
    - [ ] 執行 `go get github.com/spf13/cobra@latest`
    - [ ] 驗證 CLI 入口點 `cmd/slctl` 是否正常編譯
- [ ] Task: Conductor - User Manual Verification 'Phase 2: 核心套件升級 (Go & Cobra)' (Protocol in workflow.md)

## Phase 3: 升級其他關鍵依賴套件
- [ ] Task: 升級 GitHub API 相關套件 (`google/go-github`)
    - [ ] 執行升級指令
    - [ ] 修復因 API 變更導致的編譯錯誤 (如果有的話)
- [ ] Task: 升級日誌與工具套件 (`logrus`, `archiver`, 等)
    - [ ] 批量升級輔助套件
    - [ ] 執行 `go mod tidy` 清理依賴
- [ ] Task: Conductor - User Manual Verification 'Phase 3: 升級其他關鍵依賴套件' (Protocol in workflow.md)

## Phase 4: 全面測試與清理
- [ ] Task: 執行全面自動化測試
    - [ ] 執行 `go test ./...`
    - [ ] 驗證測試覆蓋率是否達標 (>80%)
- [ ] Task: 驗證手動關鍵功能 (Init, Plugin System)
    - [ ] 測試 `slctl init`
    - [ ] 測試插件搜尋與安裝流程
- [ ] Task: Conductor - User Manual Verification 'Phase 4: 全面測試與清理' (Protocol in workflow.md)
