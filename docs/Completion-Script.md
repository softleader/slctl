Slctl 可以產生 bash 或 zsh 的自動補齊指令, 讓你輸入指令時按下 tab 就可以獲得提示:

```sh
$ slctl completion bash  # for bash users
$ slctl completion zsh   # for zsh users
```

可以增加在 `.bashrc` 或是 `.zshrc` 中, 每次建立 shell session 時自動的 source 最新的 completion script:

```sh
$ echo "source <(slctl completion bash)" >> ~/.bashrc  # for bash users
$ echo "source <(slctl completion zsh)" >> ~/.zshrc    # for zsh users
```