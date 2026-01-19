# 實作計畫：提升測試覆蓋率至 80% 以上

本計畫遵循 TDD (測試驅動開發) 流程，旨在系統性地提升各個 Package 的覆蓋率。

## Phase 1: 初始分析與基準建立 [checkpoint: 8b2131e]
- [x] Task: 執行全專案覆蓋率掃描並產出初始報告 7b58c25
    - [ ] 執行 `go test -coverprofile=coverage.out ./...`
    - [ ] 使用 `go tool cover -func=coverage.out` 分析各 package 覆蓋率
- [x] Task: 建立低於 80% 的目標 Package 清單 88de776
- [x] Task: Conductor - User Manual Verification 'Phase 1' (Protocol in workflow.md) 8b2131e

## Phase 2: 提升核心邏輯 (pkg/) 覆蓋率
- [x] Task: 針對 `pkg/config` 補強測試並達成 >80% 覆蓋率 b6e05be
    - [ ] 撰寫測試涵蓋所有錯誤處理路徑
- [x] Task: 針對 `pkg/environment` 補強測試並達成 >80% 覆蓋率 af31cb9
- [x] Task: 針對 `pkg/github` 及其子 package 補強測試並達成 >80% 覆蓋率 f53d2ae
    - [ ] 確保 API 呼叫皆有 Mock 處理
- [ ] Task: 針對 `pkg/plugin` 及其子 package 補強測試並達成 >80% 覆蓋率
- [ ] Task: 針對其餘 `pkg/` 下的 Package (paths, formatter, etc.) 進行補強
- [ ] Task: Conductor - User Manual Verification 'Phase 2' (Protocol in workflow.md)

## Phase 3: 提升 CLI 命令 (cmd/slctl/) 覆蓋率
- [ ] Task: 針對 `cmd/slctl` 下的命令處理邏輯撰寫測試
    - [ ] 涵蓋 plugin 相關命令 (install, list, search, etc.)
    - [ ] 涵蓋 init 與 completion 命令
- [ ] Task: 確認 `cmd/slctl` 整合測試涵蓋核心流程
- [ ] Task: Conductor - User Manual Verification 'Phase 3' (Protocol in workflow.md)

## Phase 4: 最終驗證與交付
- [ ] Task: 執行最終全專案覆蓋率報告
    - [ ] 產出最終 `coverage.out` 並轉化為易讀格式
- [ ] Task: 確認所有 Package 均已超過 80% 門檻
- [ ] Task: Conductor - User Manual Verification 'Phase 4' (Protocol in workflow.md)