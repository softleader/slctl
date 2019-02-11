package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// YesNoQuestion prompt out a yes-no question to Stdin
func YesNoQuestion(out io.Writer, question string) bool {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(out, fmt.Sprintf("%s [Y/n] ", question))
		answer, _ := r.ReadString('\n')
		if ans := strings.ToLower(strings.TrimSpace(answer)); ans == "y" || ans == "yes" {
			return true
		} else if ans == "n" || ans == "no" {
			return false
		}
	}
}
