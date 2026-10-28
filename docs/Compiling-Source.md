Slctl 是使用 [Go](https://golang.org/) 寫的, 需要事先安裝好:

- [Go 1.26.10](https://golang.org/doc/install)
- [Git](https://git-scm.com/) - 程式碼管理
- make - For easy build (Optional)
- [Docker](https://www.docker.com/) - For Sandbox testing (Optional)

## 環境建置

下載 Source:

```sh
$ git clone git@github.com:softleader/slctl.git
$ cd slctl
```

取得依賴的 Dependencies:

```sh
$ make bootstrap
```

執行專案測試, 確保一切正常:

```sh
$ make test
```

## 漏洞檢查

使用 `go.mod` 指定的 Go 工具鏈執行：

```sh
$ GOTOOLCHAIN=go1.26.10 make vuln
```

此指令使用固定版本的 `govulncheck` 與最新 Go 漏洞資料庫，檢查 CLI 所用模組（不只可達函式），以及專案與測試程式碼。GitHub Actions 會在 push、PR、每週排程及 release 前執行檢查，並以 Trivy 掃描 Go 依賴與完整 sandbox 映像檔；任何嚴重程度的已知漏洞都會讓掃描失敗，不忽略尚無修補版的漏洞。

CI 在 Linux／macOS 執行測試；Windows 執行建置、vet 與漏洞掃描。既有測試包含 POSIX 路徑斷言，尚未全面移植至 Windows。

掃描通過表示掃描當時未發現資料庫已知漏洞，不代表能保證沒有未揭露或未來新增的 CVE。外部下載的 plugin、Homebrew tapper、套件管理器本身不包含在 slctl 的依賴掃描範圍內。

## Sandbox

Slctl 會使用 Docker 開啟 [Alpine](https://hub.docker.com/_/alpine/) 環境, 並將執行檔編譯好放入 container 中, 你可以在其中測試並不用擔會干擾當前的環境

```sh
$ make sanbox
```