package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type Output struct {
	OutStream, ErrStream io.Writer
}

// ログの出力を変数にも保持
var LogBuffer bytes.Buffer

// LoggingSettings ログファイルの出力
func LoggingSettings(logFile string, output *Output) error {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return fmt.Errorf("ログファイルの作成に失敗: %s", err)
	}
	// 標準出力or標準エラー出力指定
	multiLogFile := io.MultiWriter(output.OutStream, logfile, &LogBuffer)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)

	// logrusを用いたエラーメッセージの標準出力
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetLevel(logrus.WarnLevel)

	return nil
}
