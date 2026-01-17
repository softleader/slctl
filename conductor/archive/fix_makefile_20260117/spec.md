# Specification - Fix Makefile Toolchain

## Overview
修正 `Makefile` 中用於安裝開發工具（如 `goimports`, `golint`）的過時指令。

## Goals
- 將 `Makefile` 中的 `go get` 替換為現代的 `go install ...@latest`。
- 確保工具執行時能被正確找到。
- 驗證 `make link` 與 `make error-free` 能順利執行。

## Acceptance Criteria
- 執行 `make link` 不會出現 "No such file or directory" 錯誤。
- 成功構建並安裝 `slctl`。
