package c_memo

import (
	"fmt"
)

var id int
var name, text string

type AddCommand struct {
	Cmd string
}

// 簡単なコマンドの説明を記述
func (c *AddCommand) Synopsis() string {
	return "メモを作成(text,name)・検索(idあるいはnameあるいはtext)・編集(変更:idあるいはnameあるいはtext 検索:id)・削除(id)を行う(カッコ内は必要カラム)"
}

// 使い方 詳細なヘルプメッセージを返す
func (c *AddCommand) Help() string {
	return "コマンド [option...] を入力する"
}

func (c *AddCommand) Run(args []string) int {

	switch c.Cmd {
	case "add":
		// 作成
		Add(args)
	case "select":
		// 検索
		Select(args)
	case "edit":
		// 編集
		Edit(args)
	case "delete":
		// 削除
		Delete(args)
	}

	fmt.Printf("%sを実行しました\n", c.Cmd)
	return ExitCodeOK //正常終了
}
