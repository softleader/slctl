# Technology Stack - slctl

## Frontend / CLI
- **Language**: [Go 1.24+](https://golang.org/) - 核心開發語言。
- **Framework**: [Cobra](https://github.com/spf13/cobra) - 用於構建強大的 CLI 應用程式。
- **Argument Parsing**: [pflag](https://github.com/spf13/pflag) - 支援 POSIX 風格的 flag。

## Service Integration
- **Platform**: [GitHub API](https://developer.github.com/v3/) - 用於存取 GitHub 資源與權限驗證。
- **SDK**: [go-github v69+](https://github.com/google/go-github) - GitHub API 的 Go 語言客戶端庫。
- **Networking/OAuth2**: [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) - 處理 GitHub 認證流程。

## Infrastructure & Tooling
- **Logging**: [Logrus](https://github.com/sirupsen/logrus) - 結構化日誌記錄。
- **Build System**: [Makefile](https://www.gnu.org/software/make/manual/make.html) - 自動化構建、測試與發佈流程。
- **CI/CD**: [Travis CI](https://travis-ci.com/) - 持續整合平台。
- **Distribution**: 
  - [Homebrew](https://brew.sh/) (macOS/Linux)
  - [Chocolatey](https://chocolatey.org/) (Windows)
- **Containerization**: [Docker](https://www.docker.com/) - 提供 Sandbox 測試環境。
