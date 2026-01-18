# Specification - Upgrade Dependencies

## Overview
升級專案的 Go 版本、核心框架 (Cobra) 以及其他第三方依賴套件，確保專案符合現代開發標準，並解決可能存在的舊版本漏洞。

## Goals
- 將 Go 版本升級至最新穩定版 (e.g. 1.23+)。
- 更新 `go.mod` 中的核心套件（Cobra, Go-GitHub 等）。
- 確保現有功能（CLI 指令、插件系統）在升級後運作正常。
- 清理不必要的舊代碼或過時的依賴。

## Scope
- `go.mod`, `go.sum`
- `Makefile` (更新相關編譯參數)
- `cmd/`, `pkg/` 中的核心引用的 API 適配

## Constraints
- 必須維持與現有軟體功能的相容性（Breaking changes 需要特別處理）。
- 插件系統 (Plugin System) 的載入與執行機制不得受損。

## Acceptance Criteria
- `go test ./...` 完整通過。
- `slctl init` 指令運作正常。
- `slctl plugin search` 指令運作正常。
- 專案成功構建（`make build`）。
