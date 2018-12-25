package formatter

import (
	"bytes"
	"github.com/sirupsen/logrus"
)

// 什麼都不 format 的 formatter
type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(entry.Message)
	return buf.Bytes(), nil
}
