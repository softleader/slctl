package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

const (
	completionDesc = `
產生 bash 或 zsh 的 auto completion script

	$ slctl completion bash

可以增加在 '.bashrc' 或是 '.zshrc' 中:

	$ echo "source <(slctl completion bash)" >> ~/.bashrc  # for bash users
	$ echo "source <(slctl completion zsh)" >> ~/.zshrc    # for zsh users
`

	zshInitialization = `#compdef slctl
__slctl_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__slctl_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__slctl_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__slctl_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__slctl_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__slctl_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__slctl_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__slctl_filedir() {
	local RET OLD_IFS w qw
	__slctl_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __slctl_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__slctl_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__slctl_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__slctl_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__slctl_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__slctl_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__slctl_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__slctl_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__slctl_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__slctl_type/g" \
	<<'BASH_COMPLETION_EOF'
`

	zshTail = `
BASH_COMPLETION_EOF
}
__slctl_bash_source <(__slctl_convert_bash_to_zsh)
_complete slctl 2>/dev/null
`
)

var (
	completionShells = map[string]func(cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func newCompletionCmd() *cobra.Command {
	var shells []string
	for s := range completionShells {
		shells = append(shells, s)
	}
	cmd := &cobra.Command{
		Use:       "completion SHELL",
		Short:     "generate completion script for the specified shell (bash or zsh)",
		Long:      completionDesc,
		Args:      cobra.ExactArgs(1),
		ValidArgs: shells,
		RunE: func(cmd *cobra.Command, args []string) error {
			run := completionShells[strings.ToLower(args[0])]
			return run(cmd)
		},
	}
	return cmd
}

func runCompletionBash(cmd *cobra.Command) error {
	return cmd.Root().GenBashCompletion(logrus.StandardLogger().Out)
}

// runCompletionZsh cobra 的 GenZshCompletion 壞了, 在還沒修好前先自己寫..
func runCompletionZsh(cmd *cobra.Command) error {
	io.WriteString(logrus.StandardLogger().Out, zshInitialization)
	buf := new(bytes.Buffer)
	if err := cmd.Root().GenBashCompletion(buf); err != nil {
		return err
	}
	logrus.StandardLogger().Out.Write(buf.Bytes())
	io.WriteString(logrus.StandardLogger().Out, zshTail)
	return nil
}
