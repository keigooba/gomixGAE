package utils

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

// 出力の先のテスト
func TestLoggingSettings(t *testing.T) {
	// 現在の日付
	nowDate := time.Now().Format("200601")
	logFile := Pwd + "/" + "log/system_" + nowDate + ".log"

	// 実行したメッセージをメモリから取得できる
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	output := &Output{OutStream: outStream, ErrStream: errStream}

	// 標準出力or標準エラー出力指定
	err := LoggingSettings(logFile, output)
	if err != nil {
		t.Errorf("%sの作成に失敗しました", logFile)
	}

	// これが標準出力に出力されているか
	test := "テスト"
	log.Printf("ログに%sは出力されている", test)

	expected := fmt.Sprintf("%sは出力されている", test)
	if !strings.Contains(outStream.String(), expected) {
		t.Errorf("%qには%qが含まれていない", outStream.String(), expected)
	}

}
