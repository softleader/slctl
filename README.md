![beta](https://img.shields.io/badge/stability-beta-darkorange.svg)
[![Build Status](https://travis-ci.com/softleader/slctl.svg?token=4jYjzyvNx4sjHcYtGC5V&branch=master)](https://travis-ci.com/softleader/slctl)

# slctl

Slctl is a command line interface for running commands against SoftLeader services.

## Install

Slctl 建議透過套件管理來安裝:

- MacOS 使用者建議使用 [Homebrew](https://brew.sh)

	```sh
	brew install softleader/tap/slctl
	```

- Windows 或 Linux 使用者建議使用 [GoFish](https://gofi.sh/):

	```sh
	gofish add https://github.com/softleader/fish-food
	gofish install slctl
	```

你也可以直接從 [Releases page](https://github.com/softleader/slctl/releases) 下載執行檔, 將其解壓縮後加入 PATH 中即可使用

當然, 你也可以參考 [Builing Soruce](https://github.com/softleader/slctl/wiki/Building-Source) 從原始碼編譯後使用, Happing Hacking :cat::computer:

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

使用 `--force` 在發現有重複的 Token 時, 會強制刪除既有的並產生一個全新的 Access Token, 若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 `--token` 讓 slctl 驗證後直接儲存起來

```sh
$ slctl init -f
$ slctl init --token GITHUB_TOKEN
```

使用 `--offline` 則 slctl 不會跟 GitHub API 有任何互動, 只會配置 [$SL_HOME](https://github.com/softleader/slctl/wiki/Home-Path) 環境目錄.

同時使用 `--offline` 及 `--token` 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)

## Plugins

Slctl 有完整的 Plugin 系統, 你可以從收尋公司官方的 Plugin 開始:

```sh
$ slctl plugin search FILTER..
```

使用空白分隔傳入多個 FILTER, 會以 Or 且模糊條件來過濾 SOURCE; 反之列出全部

```sh
$ slctl plugin search
$ slctl plugin search whereis contacts
```

傳入 `--installed` 只列出已安裝的 Plugin

```sh
$ slctl plugin search -i
```

查詢的結果將會被 cache 並留存一天, 傳入 `--force` 可以強制更新 cache

```sh
$ slctl plugin search -f
```

### Installing Plugins

執行 `$ slctl plugin install SOURCE` 即可安裝, 點擊 [Installing a Plugin](https://github.com/softleader/slctl/wiki/Plugins-Guide#installing-a-plugin) 查看多種 *SOURCE* 的安裝方式

### Upgrading Plugins

*Slctl* 支援 GitHub Repo 的 Plugin 自動更新, 請參考 [Upgrading Plugins](https://github.com/softleader/slctl/wiki/Plugins-Guide#upgrading-plugins)

### Writing Custom Plugins

*Slctl* 支援任何語言的 Plugin, 請參考 [Plugins Guide](https://github.com/softleader/slctl/wiki/Plugins-Guide)

