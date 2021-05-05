package c_memo

import (
	"flag"
	"fmt"
	cliCmd "gomix/cli"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

const (
	ExitCodeOK        int = iota // 0 iota 連番になる
	ExitCodeError                // 1
	ExitCodeFileError            // 2
)

func flagSet(args []string) (intvar int, strvar string) {

	// selectのオプションフラグの設定
	flags := flag.NewFlagSet(os.Args[2], flag.ContinueOnError)
	flags.IntVar(&intvar, "id", 0, "ID")
	flags.StringVar(&strvar, "name", "", "名前")
	flags.StringVar(&strvar, "text", "", "テキスト")
	if err := flags.Parse(args); err != nil {
		log.Printf("メモの%sフラグをパースできないため、値を確認して下さい: %s", os.Args[2], err)
	}

	return intvar, strvar
}

func CliMemo() int {
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

	return exitCode
}
