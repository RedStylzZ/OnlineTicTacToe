package logging

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	TimestampFormat string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logMessage string
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message
	logMessage = fmt.Sprintf("[%s] %s - %s\n", level, timestamp, message)
	return []byte(logMessage), nil
}
