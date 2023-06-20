package middleware

import (
	"fmt"
	"strings"

	"github.com/yulecd/pp-common/plog"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(&recoverLogger{plog.GetDefaultFieldEntry(nil)})
}

type recoverLogger struct {
	*plog.Entry
}

func (l *recoverLogger) Write(p []byte) (n int, err error) {
	str := string(p)
	str = strings.Replace(str, "\n\t", " ", -1)
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.TrimSuffix(str, "\033[0m")
	str = strings.TrimPrefix(str, "\u001b[31m")
	data := strings.Split(str, "\n")
	newData := make([]string, 0, len(data))
	for i := range data {
		line := strings.Trim(data[i], "\r")
		if len(line) > 0 {
			newData = append(newData, line)
		}
	}
	l.WithField("stack", newData).Info("recovered")
	fmt.Print(string(p))
	return len(p), err
}
