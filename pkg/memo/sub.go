package memo

import (
	"gomix/pkg"
	"io/ioutil"
	"log"
	"os"
)

//以下機能は主要機能に導入なし

// 一時ディレクトリの作成・削除
func Dosomething() error {
	err := os.MkdirAll("newdir", 0755)
	if err != nil {
		log.Println(err)
	}
	//  (2)ディレクトリ削除
	defer func() {
		err := os.RemoveAll("newdir")
		if err != nil {
			log.Println(err)
		}
	}()

	f, err := os.Create("newdir/newfile")
	if err != nil {
		log.Println(err)
	}
	// (1)ファイルハンドルが閉じられる
	defer pkg.Close(f)

	return nil
}

// ファイルの作成・名前変更・deferの操作
func MytemFile() (*os.File, error) {
	file, err := ioutil.TempFile("", "temp") //適当なディレクトリ/tempランダム文字列 ファイルの作成
	if err != nil {
		return nil, err
	}
	// defer file.Close() //Closeが遅い deferはfunc()の呼び出し形式を取る 引数にはdeferを呼び出した時点の値が入る
	pkg.Close(file) //Renameを実行するため、すぐ閉じる

	// defer file.Close()するとwindowsではファイルが開かれていると認識され、Renameできない
	if err = os.Rename(file.Name(), file.Name()+".go"); err != nil {
		return nil, err
	}
	return file, nil
}
