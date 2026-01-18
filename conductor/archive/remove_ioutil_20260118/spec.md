# 規格書：移除 deprecated io/ioutil 套件

## Overview
將程式碼庫中所有使用 deprecated `io/ioutil` 套件的地方替換為 Go 1.16+ 推薦的替代方案（`io` 和 `os` 套件）。此變更確保程式碼與 Go 1.24 現代化升級的精神一致。

## 背景
`io/ioutil` 套件自 Go 1.16 起已被標記為 deprecated。官方建議：
- `ioutil.ReadAll` → `io.ReadAll`
- `ioutil.ReadFile` → `os.ReadFile`
- `ioutil.WriteFile` → `os.WriteFile`
- `ioutil.TempDir` → `os.MkdirTemp`
- `ioutil.TempFile` → `os.CreateTemp`
- `ioutil.ReadDir` → `os.ReadDir`

## 影響範圍
以下 11 個檔案需要修改：
1. `pkg/github/auth.go`
2. `pkg/environment/migrator.go`
3. `pkg/environment/migrator_test.go`
4. `pkg/config/config.go`
5. `pkg/config/config_test.go`
6. `pkg/plugin/create.go`
7. `pkg/plugin/loader.go`
8. `pkg/plugin/cleanup.go`
9. `pkg/plugin/repository.go`
10. `pkg/plugin/installer/archive_installer.go`
11. `pkg/plugin/installer/archive_installer_test.go`

## Functional Requirements
1. 將所有 `ioutil.ReadAll` 呼叫替換為 `io.ReadAll`
2. 將所有 `ioutil.ReadFile` 呼叫替換為 `os.ReadFile`
3. 將所有 `ioutil.WriteFile` 呼叫替換為 `os.WriteFile`
4. 將所有 `ioutil.TempDir` 呼叫替換為 `os.MkdirTemp`
5. 將所有 `ioutil.TempFile` 呼叫替換為 `os.CreateTemp`
6. 將所有 `ioutil.ReadDir` 呼叫替換為 `os.ReadDir`
7. 更新對應的 import 語句，移除 `io/ioutil`，新增必要的 `io` 或 `os`

## Acceptance Criteria
1. 程式碼庫中不再有任何 `io/ioutil` 的 import
2. `make error-free` 通過（goimports, gofmt, golint, go vet）
3. `make test` 通過，所有現有測試正常運作
4. `git status` 在驗證後為 clean（無未預期的變更）

## Out of Scope
- 手動測試 `slctl init` 等指令
- 新增額外的單元測試
- 功能性變更
