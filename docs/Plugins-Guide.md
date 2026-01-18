Slctl 有完整的 plugin 系統以延伸開發更多的功能, 包含許多優點:

1. 支援多種[安裝來源](#installing-a-plugin)
2. 不限制撰寫的語言
3. 取得 slctl 本身的[使用者設定](#environment-variables), 如 GitHub Access Token

## Stateless Plugins

Slctl 建議將 Plugin 設計為 Stateless, 也就本身無狀態的. 若需要儲存資料或狀態, 儲存的位置請參考 [Mount Volume](#mount-volume) 章節, 但必須考慮到日後升級 Plugin 時的資料轉移.

## Installing a Plugin

Plugin 可以透過指令 `$ slctl plugin install SOURCE` 來安裝, *SOURCE* 支援以下幾種:

### Local Reference

Plugin 可以是本機上的任何目錄, 透過給予絕對或相對路徑來安裝

```sh
$ slctl plugin install /path/to/plugin-dir/
```

### Archive

Plugin 也可以是來自於網路上或在本機中壓縮檔, 透過給予網址或路徑來安裝, 格式支援: *.zip*, *.tar*, *.tar.gz*, *.tgz*, *.tar.bz2*, *.tbz2*, *.tar.xz*, *.txz*, *.tar.lz4*, *.tlz4*, *.tar.sz*, *.tsz*, *.rar*, *.bz2*, *.gz*, *.lz4*, *.sz*, *.xz*

```sh
$ slctl plugin install /path/to/plugin-archive.zip
$ slctl plugin install http://host/plugin-archive.zip
```

### GitHub Repo

Plugin 也可以是一個 GitHub repo, 傳入 `github.com/OWNER/REPO`, slctl 會自動收尋最新一版的 release, 並從該 release 的所有下載檔中, 嘗試找出含有當前 OS 名稱的壓縮檔來安裝, 當找不到時會改下載第一個壓縮檔來安裝

```sh
$ slctl plugin install github.com/softleader/slctl-whereis
```

傳入 `--tag` 可以指定 release 版本

```sh
$ slctl plugin install github.com/softleader/slctl-whereis --tag 1.0.0
```

傳入 `--tag` 及 `--asset` 可以指定 release 版本以及要下載第幾個 asset 檔案 (從 0 開始) 來安裝

```sh
$ slctl plugin install github.com/softleader/slctl-whereis --tag 1.0.0 --asset 2
```

傳入 `--force` 在 install 時自動刪除已存在的 plugin

```sh
$ slctl plugin install github.com/softleader/slctl-whereis -f
```

傳入 `--dry-run` 可以模擬真實的 install, 但不會真的影響當前的配置 

```sh
$ slctl plugin install github.com/softleader/slctl-whereis -f --dry-run
```

> Read more about [SoftLeader Official Plugins](./Official-Plugins.md)

## Upgrading Plugins

Upgrade plugin which installed from GitHub Repo

```sh
$ slctl plugin upgrade NAME...
```

*NAME* 可傳入指定要更新的 Plugin 完整名稱 (一或多個, 以空白區隔); 反之更新全部

```sh
$ slctl plugin upgrade
$ slctl plugin upgrade whereis
```

傳入 `--tag` 可以指定要更新的 release 版本

```sh
$ slctl plugin upgrade whereis --tag 1.0.0
```

傳入 `--tag` 及 `--asset` 可以指定要更新的 release 版本以及要下載第幾個 asset 檔案 (從 0 開始)

```sh
$ slctl plugin upgrade whereis --tag 1.0.0 --asset 2
```

傳入 `--dry-run` 可以模擬真實的 upgrade, 但不會真的影響當前的配置, 通常可以用來檢查 plugin 是否有新版的再決定是否要更新

```sh
$ slctl plugin upgrade --dry-run
```

## Building Plugins

Plugin 的根目錄下必須有一份 `metadata.yaml` 檔案來描述該 Plugin 的相關資訊, 包含:

```
name: foo
version: 0.1.0
usage: foo
description: The foo plugin
exec:
  command: $SL_BIN
  platform:
  - os: darwin
    arch: ""
    command: $SL_BIN
  - os: windows
    arch: ""
    command: $SL_BIN
hook:
  command: echo hello foo
  platform: []
ignoreGlobalFlags: false
github:
  scopes: []
```

- *name* - 從 slctl 執行的 subcommand, 如 `foo` 即代表將使用 `slctl foo` 來執行, *name* 不可跟 slctl 第一層 command 或任何其他 plugin name 重複, *name* 也必須符正規式驗證: `^[\w\d_-]+$`
- *version* - 必須是符合 [Semantic Versioning 2](https://semver.org/) 規格的版本號
- *usage* - 顯示在 `slctl --help` 的說明, 應以一句話簡短說明
- *description* - 顯示在 `slctl plugin list` 的說明
- *exec* - 執行 plugin 時的指令, 如果有定義 platform 且符合 runtime 環境, 會優先選擇 platform 的 command 執行
- *hook* - 安裝好 plugin 後要執行的指令, 通常用來配置環境使用, 如果有定義 platform 且符合 runtime 環境, 會優先選擇 platform 的 command 執行
- *ignoreGlobalFlags* - 執行 plugin 時忽略 [Global Flags](#global-flags)
- *github.scopes* - plugin 需要的 token 權限, slctl 會自動的補齊不足的權限, 執行 `slctl init scopes` 可查詢 slctl 預設要求的權限清單

> Os 及 Arch 可參考: [Golang environment variables](https://golang.org/doc/install/source#environment)

### Plugin Templates

Plugin 本身沒有撰寫的語言限制, slctl 推薦並預設產生 golang 的範本, 選擇不同撰寫語言時, 需注意該語言本身的限制: 如執行 java plugin 的 runtime 必須有 JRE

slctl 已內含了幾種語言的範本, 使用 `--lang` 來指定產生語言範本

```sh
$ slctl plugin create foo --lang java
```

使用 `plugin create langs` 列出所有內含的範本語言

```sh
$ slctl plugin create langs
golang
java
nodejs
...
```

Slctl 預設會在當前目錄下, 建立一個名為 Plugin 名稱的目錄, 並將範本產生在該目錄中, 可以傳入 `--output` 來指定 Plugin 的產生目錄

```sh
$ slctl plugin create foo -o /path/to/plugin-dir
```

## Environment Variables

Slctl 在執行 plugin 時, 會將以下環境變數設置到系統變數中, 讓不同語言的 plugin 都可以從中取出, 如:

- *golang* - `os.LookupEnv("SL_TOKEN")` 
- *java* - `System.getenv("SL_TOKEN")`
- *nodejs* - `process.env.SL_TOKEN`

這些變數有個共同點都是 *SL_* 開頭, 可以透過 `plugin evns` 查看變數清單:

```sh
$ slctl plugin envs
SL_PLUGIN_DIR=~/.sl/plugins/foo
SL_PLUGIN_MOUNT=~/.sl/mounts/foo
SL_PLUGIN=~/.sl/plugins
SL_VERBOSE=false
SL_CLI=slctl
SL_VERSION=<version.of.slctl>
SL_HOME=~/.sl/
SL_OFFLINE=false
SL_TOKEN=<github.token>
SL_PLUGIN_NAME=foo
SL_BIN=slctl
...
```

## Mount Volume

Slctl 會依照當前 [Home Path](./Home-Path.md) 的設定, 分派一個 Mount Volume 給該 Plugin 使用, Slctl 會保證:

- Mount Volume 在執行的當下一定會存在
- Mount Volume 不會因為移除 Plugin 而遭到刪除

Slctl 會將 Mount Volume 的路徑以 `$SL_PLUGIN_MOUNT` 變數存放, 以 golang 語言為例:

```go
mount, found := os.LookupEnv("SL_PLUGIN_MOUNT")
if !found {
	return errors.New("$SL_PLUGIN_MOUNT not found")
}
```

Mount Volume 的路徑只會以 Plugin Name 作區分, 也就是同個 Plugin Name 在不同版本之間將會使用相同的 Mount Volume, 因此 Plugin 版本更新時一定要考慮到資料的更新

### Umount

Slctl 提供了指令可以移除 Plugin 的 Mount Volume, 你可以視為將 Plugin 還原至初始狀態:

```sh
$ slctl plugin umount PLUGIN_NAME
```

使用前請注意: `umount` 指令不可逆

## Global Flags

執行 slctl 時, 有些 flag 會優先被 slctl 使用, 我們稱這些為 global flags, 執行 `plugin flags` 會列出 global flags 清單:

```sh
$ slctl plugin flags
--home
--offline
--verbose
-v
...
```

執行 plugin 時, global flags 預設**不會**傳給 plugin, 若 plugin 將 `metadata.yaml` 中的 `ignoreGlobalFlags` 設為 true, 則 slctl 會將所有 flags 完整的傳入 plugin 中