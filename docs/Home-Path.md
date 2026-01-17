Slctl 預設使用 `${user.home}/.config/slctl` 作為 home path, 用來存放 slctl 使用到的檔案, 如 configs, plugins 等:

- *Unix/MacOS* - `~/.config/slctl`
- *Windows 10* - `C:/Users/<username>/.config/slctl`

> *3.6.x* 之前的版本預設的 home 是 `~/.sl`, 升級 *3.7.x* 會自動更新目錄到 `~/.config/slctl`

可以使用 `slctl home` 查看當前的 home path:

```sh
$ slctl home
/Users/matt/.config/slctl
```

**home path 不允許包含任何空白** (windows 7 的使用者請特別注意), 可以參考下一章節來做調整: [Change Home Path](#change-home-path)

## Change Home Path

設定系統變數 `SL_HOME` 將改變 slctl 的 global home path, 依照不同的作業系統有不同的設定, 請參考 [Set System Variables](#set-system-variables) 設置

你也可以在執行任何 slctl command 時, 手動傳入 `--home` 指定***當次*** command 的 home path (Overrides *$SL_HOME*) 

```sl
$ slctl home --home /some/where/else
/some/where/else
```

## Set System Variables

### Windows

1. Open system properties
1. Select the "*Advanced*" tab, and the "*Environment Variables*" button
1. Click the first "*New...*" button.
1. Set *SL_HOME* to the location of your slctl home. e.g. `C:\.config\slctl`

### Unix-based Operating System (Linux, Solaris and Mac OS)

1. Find out if you're using *Bash* or *ZSH* by running the command `echo $SHELL` in your Terminal.
1. Exporting *SL_HOME* in your shell config (`~/.bashrc` or `~/.zshrc`), e.g. 

```sh
export SL_HOME=/.config/slctl
```