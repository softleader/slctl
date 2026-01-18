Slctl 是使用 [Go](https://golang.org/) 寫的, 需要事先安裝好:

- [Go 1.25.6](https://golang.org/doc/install)
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

## Sandbox

Slctl 會使用 Docker 開啟 [Alpine](https://hub.docker.com/_/alpine/) 環境, 並將執行檔編譯好放入 container 中, 你可以在其中測試並不用擔會干擾當前的環境

```sh
$ make sanbox
```