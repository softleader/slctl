# Go 1.26.10 與依賴安全更新

- [x] 確認 Go 1.26.10 正式發布；以該工具鏈執行既有測試與漏洞基線掃描。
- [x] 更新直接與間接依賴，移除無修補版的 archiver/rardecode 舊版依賴。
- [x] 補上壓縮檔安全回歸測試，驗證正常安裝與惡意路徑、連結、覆寫的拒絕行為。
- [x] 更新建置文件、GitHub Actions 與容器基底，加入持續漏洞掃描。
- [x] 使用 Go 1.26.10 完成測試、vet、跨平台建置與漏洞掃描，記錄結果與範圍。

## 漏洞基線

以 Go 1.26.10 執行 `govulncheck@v1.9.0 -show verbose ./...`：
- GO-2025-4020：`github.com/nwaples/rardecode@v1.1.3`，無修補版。
- GO-2025-3605、GO-2024-2698：`github.com/mholt/archiver@v3.1.1+incompatible`，無修補版。
- GO-2026-5024：`golang.org/x/sys@v0.40.0`，需升級至 v0.44.0 以上。

此工作僅修改工作目錄，不自動提交或發布。漏洞結論限於掃描當時的公開資料庫與實際掃描範圍。

## 實作摘要

- Go 1.25.6 → 1.26.10；go-github v69 → v93，調整 functional options 與 release 欄位型別。
- semver v3 → v4、pb v1 → v3、YAML v2 → 維護中的 `go.yaml.in/yaml/v3`；更新可升級的直接與實際使用的間接依賴，沒有新版者保留。
- archiver → archives v0.1.5；rardecode 改用已修補的 v2.4.5。解壓縮透過 `os.Root` 限制目的路徑，拒絕連結／特殊檔案及重複寫入，保留普通檔案的執行權限。
- Docker 改用 Alpine 3.24.3 並鎖定 digest，保留原 sandbox 的 bash、curl、procps、vim、台北時區，新增 CA certificates 並更新系統套件。
- Actions 升級並鎖定 commit；GoReleaser 固定 v2.19.0。加入 push、PR、每週與 release 前的 Go／Trivy 掃描，不設定 severity 過濾、不忽略無修補版漏洞。
- 修正既有 `TestInstall_Full`：原測試建立暫存 home，卻未設定 `environment.Settings.Home` 與 release metadata，導致 `runHook → SetupEnv` 讀取開發機設定。單獨執行與乾淨 Linux 容器均重現；現在測試擁有並還原其 home、metadata 與環境變數，且檢查實際安裝結果。

## 驗證結果（2026-10-28）

所有 Go 驗證明確使用 `GOTOOLCHAIN=go1.26.10`；Linux 容器使用內建 Go 1.26.10 與 `GOTOOLCHAIN=local`。

| 驗證 | 結果 |
| --- | --- |
| `go mod verify`、`go mod tidy -diff` | 通過，模組內容驗證成功且 tidy 無差異 |
| `go test -race ./...`、`go vet ./...`（darwin/arm64） | 通過 |
| 乾淨 `golang:1.26.10-alpine3.24` 容器的 `go test ./... && go vet ./...`（linux/arm64） | 通過 |
| `CGO_ENABLED=0 go build -tags netgo` | linux/amd64、darwin/amd64、darwin/arm64、windows/amd64 全部通過 |
| `govulncheck@v1.9.0 -C cmd/slctl -scan=module` | 上述四個發行目標均為 `No vulnerabilities found` |
| `govulncheck@v1.9.0 -test ./...` | `No vulnerabilities found` |
| Trivy v0.76.1：`fs --scanners vuln --exit-code 1` | go.mod：0 漏洞 |
| Trivy v0.76.1：完整 linux/amd64 Docker image | Alpine 3.24.3 的 40 個系統套件：0；`usr/local/bin/slctl`：0 |
| `docker run --rm --platform linux/amd64 slctl:security slctl version` | 成功執行，輸出 `unreleased+unknown` |
| `actionlint@v1.7.12`、`git diff --check` | 通過 |

Trivy 掃描的映像檔 manifest：`sha256:4045e128c4652d525288be88893de6f6abb2751e97df73e339f4f988468fd015`。

本機未執行 GitHub Actions，也未實際執行 Windows 二進位檔或完整 Windows 測試。CI 對 Windows 執行建置、vet、漏洞掃描；既有 POSIX 路徑斷言的測試仍限 Linux/macOS。外部 plugin、Homebrew tapper、開發／發佈工具本身不在 Go 執行檔與 sandbox 套件的漏洞結論內。

## 測試與 mutation 對照

以下為程式路徑檢查的 mutation 對照，不宣稱執行過自動 mutation testing。安全測試在舊 archiver 實作上已觀察到紅燈（兄弟目錄逃逸、絕對／Windows 路徑、symlink、FIFO、既有 symlink 逃逸），換成新實作後轉綠。

| 可能的錯誤變更 | 應轉紅的測試與可觀察結果 |
| --- | --- |
| 不檢查反斜線、磁碟代號或 ADS | `TestExtractRejectsUnsafeEntries`：Unix 上會錯誤接受這些名稱；單獨移除 `IsLocal` 仍有 Root 防線，不視為失去所有越界保護 |
| 將 Root 操作改為一般路徑拼接寫檔 | `TestExtractCannotFollowExistingSymlinkOutsideRoot`：既有 symlink 指向的外部目錄會出現檔案；不依賴字串路徑檢查 |
| 不拒絕 symlink／hardlink／特殊檔案 | `TestExtractRejectsLinksAndSpecialFiles`：不再返回錯誤，或建立了不應存在的項目 |
| 移除 `O_EXCL` | `TestExtractDoesNotOverwriteExistingFiles`：重複項目被接受，或第一份內容遭覆寫；新實作沒有重試覆寫的 fallback |
| 不複製資料、忽略目錄、把執行檔 mode 固定成 0644 | `TestExtractPreservesFilesAndExecutableMode`：內容、空目錄或 Unix 執行權限不符 |
| 只靠副檔名判斷壓縮格式 | `TestExtractIdentifiesArchiveContents`：無副檔名 ZIP 無法還原預期內容 |
| 將無效壓縮檔當作成功 | `TestExtractRejectsInvalidArchive`：應有錯誤卻回傳 nil |
| NewTokenClient 丟失 OAuth HTTP transport 或使用錯誤 token | `TestNewTokenClientAuthenticatesRequests`：本機 HTTP 端點拒絕缺少／錯誤的 Bearer header，使用者資料無法取得 |
| 安裝回傳成功但未建立可載入的 plugin | `TestInstall_Full`：安裝目錄不存在或 metadata 不符 |
| 新 SDK 的 API 路徑／scopes／release 解析錯誤 | 既有 member、token、repository、GitHub installer 與 update 測試繼續透過本機 HTTP server 或明確 release fixture 檢查行為 |

穩定性驗證：

```sh
GOTOOLCHAIN=go1.26.10 go test ./pkg/plugin/installer \
  -run '^(TestExtract|TestArchiveInstaller_Install)' \
  -race -shuffle=20261028 -count=20
GOTOOLCHAIN=go1.26.10 go test ./cmd/slctl -run '^TestInstall_Full$' -race -count=20
GOTOOLCHAIN=go1.26.10 go test ./pkg/github -run '^TestNewTokenClientAuthenticatesRequests$' -race -count=20
```

上述三組均 20/20 通過。新解壓縮函式的 statement coverage 為 83.8%（不是整個 installer 套件的 coverage）。新安全測試使用獨立標準函式庫 fixture 與各自的 `t.TempDir()`，沒有共享 home、sleep 或網路等待；HTTP 驗證使用各測試自己的 httptest server 並清理。RAR 與其他非 ZIP/TAR/GZIP 格式尚未新增獨立 fixture；仍透過維護中的格式實作處理，這部分不宣稱已做逐格式行為驗證。
