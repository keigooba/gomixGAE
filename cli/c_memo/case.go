package c_memo

import (
	"fmt"
	"gomix/config"
	"gomix/pkg/memo"
	"strconv"
	"time"
)

func Add(args []string) {
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		text = args[1]
	}
	memoEx := memo.Memo{}
	memoEx.Name = name
	memoEx.Text = text
	memoEx.CreatedAt = time.Now()
	config.Db.Create(&memoEx)
}

func Select(args []string) {
	var search string
	var data interface{}
	// オプションフラグを設定し、値をセット
	intvar, strvar := flagSet(args)

	if len(args) > 0 {
		switch args[0] {
		case "-id":
			data = intvar
			search = "id=?"
		case "-name":
			data = strvar
			search = "name=?"
		case "-text":
			data = strvar
			search = "text=?"
		}
	}

	memoEx := []memo.Memo{}
	config.Db.Find(&memoEx, search, data)
	fmt.Println("検索結果\n", memoEx)
}

func Edit(args []string) {
	if len(args) > 2 {
		// 数値に変換
		id, _ = strconv.Atoi(args[2])
	}

	// オプションフラグを設定し、値をセット
	_, strvar := flagSet(args)
	memoExBefore := memo.Memo{}
	memoExBefore.ID = id
	memoExAfter := memoExBefore
	config.Db.First(&memoExAfter)

	if len(args) > 2 {
		switch args[0] {
		case "-name":
			memoExAfter.Name = strvar
		case "-text":
			memoExAfter.Text = strvar
		}
	}
	config.Db.Save(&memoExAfter)
}

func Delete(args []string) {
	if len(args) > 0 {
		// 数値に変換
		id, _ = strconv.Atoi(args[0])
	}
	memoEx := memo.Memo{}
	memoEx.ID = id
	config.Db.First(&memoEx)
	config.Db.Delete(&memoEx)
}
