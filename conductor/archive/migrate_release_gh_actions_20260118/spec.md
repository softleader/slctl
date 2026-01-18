# Track Specification: Migrate Release Process to GitHub Actions

## Overview
本 Track 的目標是將專案的 Release 發布流程從 Travis CI 遷移至 GitHub Actions。
根據需求，我們將專注於使用 **GoReleaser** 來處理跨平台的編譯與發布，並設定為**手動觸發**。同時，我們將移除不再使用的 `.travis.yml` 配置。

## Functional Requirements

1.  **GitHub Actions Workflow**
    - 建立一個新的 Workflow 檔案：`.github/workflows/release.yml`。
    - 觸發條件 (Trigger)：`workflow_dispatch` (手動觸發)。
    - 執行環境：`ubuntu-latest`。

2.  **Release Automation (GoReleaser)**
    - 使用標準的 `goreleaser/goreleaser-action`。
    - 配置 GoReleaser (`.goreleaser.yaml`) 以自動處理：
        - 跨平台編譯 (根據專案需求，通常為 Linux, Darwin, Windows / AMD64, ARM64 等)。
        - 打包 (Archives)。
        - 生成 Checksums。
        - 發布 Artifacts 到 GitHub Release。
    - 確保 GoReleaser 能讀取必要的 Token (如 `GITHUB_TOKEN`) 以執行發布。

3.  **Legacy Cleanup**
    - 刪除 `.travis.yml` 檔案。

## Non-Functional Requirements
- **Maintainability**：使用標準 Action (GoReleaser)，減少維護自定義腳本的成本。
- **Security**：使用 `GITHUB_TOKEN` 進行身份驗證，避免長期 Token 暴露。

## Out of Scope
- 特定的 webhook 觸發 (原 `.travis.yml` 中的 curl 腳本)。
- Chocolatey 發布 (原 `.travis.yml` 中的 `make choco-push`)。
- 保留或依賴現有的 `Makefile` 進行發布建置 (改用 GoReleaser 原生配置)。

## Acceptance Criteria
- [ ] 能夠在 GitHub Actions 頁面手動觸發 Release Workflow。
- [ ] Workflow 執行成功，並使用 GoReleaser 完成編譯。
- [ ] 編譯產出的 Artifacts 正確上傳至 GitHub Releases。
- [ ] 專案根目錄不再包含 `.travis.yml`。
