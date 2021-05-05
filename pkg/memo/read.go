package memo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gomix/config"
	"gomix/pkg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type memos []Memo

type Memo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	CreatedAt time.Time
}

func ReadJson(w http.ResponseWriter, r *http.Request) {

	var data Data
	var name string
	// 作成済みのJSONファイルを読み込む
	if files, err := ioutil.ReadDir("doc/memo/data/json"); err == nil {
		for _, file := range files {
			data.Json = append(data.Json, file.Name())
			name = filepath.Join(pkg.Getpath(), "doc", "memo", "data", "json", file.Name())
			f, err := os.Open(name)
			if err != nil {
				log.Println(err)
			}
			defer pkg.Close(f)

			//出力ファイルをバッファリングする
			reader := bufio.NewReader(f)
			// 先頭の1byteを覗き見る
			b, _ := reader.Peek(1)
			if string(b) == "[" {
				var memoJson memos
				dec := json.NewDecoder(reader)
				if err := dec.Decode(&memoJson); err != nil {
					log.Println(err)
				} else {
					for _, v := range memoJson {
						// 構造体のインスタンス化
						memoEx := Memo{}
						memoEx.Name = v.Name
						memoEx.Text = v.Text
						memoEx.CreatedAt = time.Now()
						config.Db.Create(&memoEx)
					}
				}
			} else {
				//JSON形式でなければエラーを返す
				log.Println("このファイルはJSON形式で書かれていません")
			}
		}
	}
	url := config.Config.URL + fmt.Sprint(config.FlagPort) + "/memo"
	http.Redirect(w, r, url, http.StatusSeeOther) //キャッシュを残したくないので、303指定
}
