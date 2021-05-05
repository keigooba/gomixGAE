package c_memo

import (
	"fmt"
	cliCmd "gomix/cli"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCliMemo(t *testing.T) {

	os.Args = strings.Split("./gomix -memo add テスト テスト2", " ")
	os.Args = strings.Split("./gomix -memo select -id 2", " ")
	os.Args = strings.Split("./gomix -memo edit -name ササエ 2", " ")
	os.Args = strings.Split("./gomix -memo delete 3", " ")


	// CLI structを生成する
	// 以下ではこのstructに書く設定を追加していく
	c := cli.NewCLI("-memo", cliCmd.Version)

	var cmd string
	if len(os.Args) > 2 {
		// ユーザの引数を登録する
		c.Args = os.Args[2:]
		cmd = os.Args[2]
	}

	// サブコマンドを登録する
	// cli.CommandFactoryという関数である
	c.Commands = map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			return &AddCommand{
				Cmd: cmd,
			}, nil
		},
		"select": func() (cli.Command, error) {
			return &AddCommand{
				Cmd: cmd,
			}, nil
		},
		"edit": func() (cli.Command, error) {
			return &AddCommand{
				Cmd: cmd,
			}, nil
		},
		"delete": func() (cli.Command, error) {
			return &AddCommand{
				Cmd: cmd,
			}, nil
		},
	}

	// コマンドを実行する
	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("cliコマンドの実行に失敗しました: %s\n", err)
	}

	if status := CliMemo(); status != exitCode || status != ExitCodeOK {
		t.Errorf("ステータスコードは%dではなく、%dになっています", ExitCodeOK, status)
	}

	// t.Log, t.Logf でログを出すと `go test -v` と実行したときのみ表示される
	t.Logf("result is %d", exitCode)
}
