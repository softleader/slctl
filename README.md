# slctl

Slctl is a command line interface for running commands against SoftLeader services.

## Install

You can also use [Homebrew](https://brew.sh/index_zh-tw) (on macOS):

```sh
brew install softleader/tap/slctl
```

Or install using [Chocolatey](https://chocolatey.org/) (on Windows):

```sh
TODO
```

Or manually downlaod from [releases page](https://github.com/softleader/slctl/releases).

## Getting Started 

執行 `slctl init` 透過互動式的問答自動的產生並儲存 GitHub Personal Access Token (https://github.com/settings/tokens)
也可以傳入 `--username` 或 `--password` 來整合非互動式的情境 (e.g. DevOps pipeline):

```sh
$ slctl init
$ slctl init -u <GITHUB_USERNAME> -p <GITHUB_PASSWORD>
```

執行 'scopes' 可以列出所有 slctl 需要的 Access Token 權限

```sh
$ slctl init scopes
```

使用 `--refresh` 讓 slctl 發現有重複的 Token 時, 自動的刪除既有的並產生一個全新的 Access Token
若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 `--token` 讓 slctl 驗證後直接儲存起來

```
$ slctl init --refresh
$ slctl init --token <GITHUB_TOKEN>
```

使用 `--offline` 則 slctl 不會跟 GitHub API 有任何互動, 只會配置 *$SL_HOME* 環境目錄.
同時使用 `--offline` 及 `--token` 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)

## Plugins

- [foo](#) - The foo plugin
- [bar](#) - The bar plugin

> TODO: Plugin 清單

### Writing Custom Plugins

*Slctl* 支援任何語言的 Plugin, 請參考 [Plugins Guide](https://github.com/softleader/slctl/wiki/Plugins-Guide)

