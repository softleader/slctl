![beta](https://img.shields.io/badge/stability-beta-darkorange.svg)
[![Build Status](https://travis-ci.com/softleader/slctl.svg?token=4jYjzyvNx4sjHcYtGC5V&branch=master)](https://travis-ci.com/softleader/slctl)

# slctl

Slctl is a command line interface for running commands against SoftLeader services.

## Install

MacOS 使用者可以透過 [Homebrew](https://brew.sh/index_zh-tw) 來安裝:

```sh
brew install softleader/tap/slctl
```

Windows 使用者可以透過 [Chocolatey](https://chocolatey.org/) 來安裝:

```sh
TODO
```

你也可以參考 [Builing Soruce](https://github.com/softleader/slctl/wiki/Building-Source) 來 hacking slctl (:cat::computer:) 或是從 [releases page](https://github.com/softleader/slctl/releases) 直接下載執行檔.

## Getting Started 

執行 `slctl init` 透過互動式的問答產生並儲存 [GitHub Personal Access Token](https://github.com/settings/tokens), 也可以傳入 `--username` 或 `--password` 來整合非互動式的情境 (e.g. DevOps pipeline):

```sh
$ slctl init
$ slctl init -u GITHUB_USERNAME -p GITHUB_PASSWORD
```

執行 `scopes` 可以列出所有 slctl 需要的 Access Token 權限

```sh
$ slctl init scopes
```

使用 `--force` 讓 slctl 發現有重複的 Token 時, 強制刪除既有的並產生一個全新的 Access Token, 若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 `--token` 讓 slctl 驗證後直接儲存起來

```sh
$ slctl init -f
$ slctl init --token GITHUB_TOKEN
```

使用 `--offline` 則 slctl 不會跟 GitHub API 有任何互動, 只會配置 [$SL_HOME](https://github.com/softleader/slctl/wiki/Home-Path) 環境目錄. 同時使用 `--offline` 及 `--token` 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)

## Plugins

Searching SoftLeader official plugin

```sh
$ slctl plugin search NAME
```

*NAME* 可傳入指定的 Plugin 名稱, 會視為模糊條件來過濾; 反之列出全部

```sh
$ slctl plugin search
$ slctl plugin search whereis
```

查詢的結果將會被 cache 留存最多一天, 傳入 `--force` 可以強制更新 cache

```sh
$ slctl plugin search -f
```

執行 `$ slctl plugin install SOURCE` 即可安裝, 點擊 [Installing a Plugin](https://github.com/softleader/slctl/wiki/Plugins-Guide#installing-a-plugin) 查看多種 *SOURCE* 的安裝方式

以下列出所有 Plugin 清單 (包含官方或個人貢獻的 Plugin)

- [softleader/slctl-whereis](https://github.com/softleader/slctl-whereis) - 快速查看同事現在在哪兒
- [softleader/slctl-make](https://github.com/softleader/slctl-make) - 在不同作業系統間都可以使用 GUN Make
- [softleader/slctl-contacts](https://github.com/softleader/slctl-contacts) - 查看公司通訊錄

### Upgrading Plugins

*Slctl* 支援 GitHub Repo 的 Plugin 自動更新, 請參考 [Upgrading Plugins](https://github.com/softleader/slctl/wiki/Plugins-Guide#upgrading-plugins)

### Writing Custom Plugins

*Slctl* 支援任何語言的 Plugin, 請參考 [Plugins Guide](https://github.com/softleader/slctl/wiki/Plugins-Guide)

