![beta](https://img.shields.io/badge/stability-beta-darkorange.svg)
[![Build Status](https://travis-ci.com/softleader/slctl.svg?branch=master)](https://travis-ci.com/softleader/slctl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softleader/slctl/blob/master/LICENSE)
[![Build Status](https://github-basic-badges.herokuapp.com/release/softleader/slctl.svg)](https://github.com/softleader/slctl/releases)

# slctl

Slctl is a command line interface for running commands against SoftLeader services.

## Install

所有執行檔都會發佈並保留在 [Releases page](https://github.com/softleader/slctl/releases), 選擇版本及對應 OS 的執行檔, 下載後將其解壓縮並加入 PATH 中即可使用

除了直接下載外, Slctl 支援了並**優先推薦**使用以下幾種套件管理程式來安裝:

- [Homebrew](https://brew.sh) 是在 Mac 上很受歡迎的套件管理系統, 推薦 MacOS 使用者使用:

	```sh
	$ brew install softleader/tap/slctl
	```

    > Linux 使用者建議使用 [Linuxbrew](http://linuxbrew.sh/) - Homebrew 的 Linux 分支

- [Chocolatey](https://chocolatey.org/) 是 Windows 上常見的的套件管理程式, 也被許多 CI 整合 (e.g. [Travis](https://docs.travis-ci.com/user/reference/windows/#chocolatey)), 推薦 Windows 使用者使用:

	```sh
	$ choco install slctl -s http://ci.softleader.com.tw:8081/repository/choco/
	```

### Compiling from source

你可以參考 [Compiling Source](https://github.com/softleader/slctl/wiki/Compiling-Source) 章節, 從原始碼編譯後使用, Happy Hacking :cat::computer:

## Getting Started

執行 `slctl init` 透過互動式的問答產生並儲存 [GitHub Personal Access Token](https://github.com/settings/tokens), 也可以傳入 `--username` 或 `--password` 來整合非互動式的情境 (e.g. DevOps pipeline):

```sh
$ slctl init
$ slctl init -u GITHUB_USERNAME -p GITHUB_PASSWORD
```

使用 `--force` 在發現有重複的 Token 時, 會強制刪除並產生一個全新的 Access Token

```sh
$ slctl init -f
```

若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 `--token` 讓 slctl 驗證後直接儲存起來, 執行 `scopes` 可以列出所有 slctl 需要的 Access Token 權限

```sh
$ slctl init --token GITHUB_TOKEN
$ slctl init scopes
```

使用 `--offline` 則 slctl 不會跟 GitHub API 有任何互動, 只會配置 [$SL_HOME](https://github.com/softleader/slctl/wiki/Home-Path) 環境目錄.

同時使用 `--offline` 及 `--token` 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)

## Plugins

Slctl 有完整的 Plugin 系統, 你可以從收尋公司官方的 Plugin 開始:

```sh
$ slctl plugin search FILTER...
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

執行 `$ slctl plugin install SOURCE` 即可安裝

除了 GitHub Repo Source 外, Slctl 還支援了許多的 SOURCE 來源, 點擊 [Installing a Plugin](https://github.com/softleader/slctl/wiki/Plugins-Guide#installing-a-plugin) 查看更多的 *SOURCE* 說明

### Upgrading Plugins

*Slctl* 支援 GitHub Repo 的 Plugin 自動更新, 請參考 [Upgrading Plugins](https://github.com/softleader/slctl/wiki/Plugins-Guide#upgrading-plugins)

### Writing A Plugin

*Slctl* 支援任何語言的 Plugin, 請參考 [Plugins Guide](https://github.com/softleader/slctl/wiki/Plugins-Guide)
