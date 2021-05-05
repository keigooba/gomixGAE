package main

import (
	"fmt"
	cliCmd "gomix/cli"
	"gomix/cli/c_memo"
	"gomix/config"
	"gomix/pkg/memo"
	"os"
)

func main() {

	// マイグレーション
	config.Db.AutoMigrate(&memo.Memo{})

	// オプションコマンドの設定
	cliCmd.CmdFlag()

	// コマンド入力の有無
	if len(os.Args) > 1 {

		// サーバー停止の通知設定
		go signalCall()

		// -memo,-mが入力された時
		if cliCmd.ManageMemo {
			// サブコマンドの設定
			exitCode := c_memo.CliMemo()
			fmt.Printf("終了ステータスコードは%dです", exitCode)
			// サブコマンド使用時にはexitする
			os.Exit(exitCode)
		}
	}

	// エントリーポイントの設定・サーバー起動
	err := StartMainServer()
	if err != nil {
		panic(err)
	}
}
