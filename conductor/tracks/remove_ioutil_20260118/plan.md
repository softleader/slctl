# 實作計畫：移除 deprecated io/ioutil 套件

## Phase 1: 核心套件更新 [checkpoint: d8be100]

- [x] Task: 替換 `pkg/github/auth.go` 中的 ioutil 使用 [8d5d6cb]
    - [x] 識別並替換 `ioutil.ReadAll` → `io.ReadAll`
    - [x] 更新 import 語句
    - [x] 執行 `go vet` 確認無錯誤

- [x] Task: 替換 `pkg/config/` 目錄中的 ioutil 使用 [9b4b0c1]
    - [x] 替換 `pkg/config/config.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/config/config_test.go` 中的 ioutil 呼叫
    - [x] 更新 import 語句

- [x] Task: 替換 `pkg/environment/` 目錄中的 ioutil 使用 [dcc8e93]
    - [x] 替換 `pkg/environment/migrator.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/environment/migrator_test.go` 中的 ioutil 呼叫
    - [x] 更新 import 語句

- [x] Task: Conductor - User Manual Verification 'Phase 1: 核心套件更新' (Protocol in workflow.md)

## Phase 2: Plugin 套件更新 [checkpoint: 3f179de]

- [x] Task: 替換 `pkg/plugin/` 根目錄中的 ioutil 使用 [f90cb94]
    - [x] 替換 `pkg/plugin/create.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/plugin/loader.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/plugin/cleanup.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/plugin/repository.go` 中的 ioutil 呼叫
    - [x] 更新 import 語句

- [x] Task: 替換 `pkg/plugin/installer/` 目錄中的 ioutil 使用 [667cf9a]
    - [x] 替換 `pkg/plugin/installer/archive_installer.go` 中的 ioutil 呼叫
    - [x] 替換 `pkg/plugin/installer/archive_installer_test.go` 中的 ioutil 呼叫
    - [x] 更新 import 語句

- [x] Task: Conductor - User Manual Verification 'Phase 2: Plugin 套件更新' (Protocol in workflow.md)

## Phase 3: 最終驗證 [checkpoint: 7f31437]

- [x] Task: 執行完整驗證套件
    - [x] 執行 `make error-free`（goimports, gofmt, golint, go vet）
    - [x] 執行 `make test` 確認所有測試通過
    - [x] 執行 `grep -r "io/ioutil" .` 確認無殘留

- [x] Task: Conductor - User Manual Verification 'Phase 3: 最終驗證' (Protocol in workflow.md)
