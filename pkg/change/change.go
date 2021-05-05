package change

import (
	"context"
	"fmt"
	"gomix/config"
	"gomix/pkg"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"golang.org/x/text/width"
)

type Data struct {
	Numbers []Number
	Error   Err
}

type Number struct {
	S2  string
	S16 string
}

type Err struct {
	Number string
}

// 正規表現
var numReg = regexp.MustCompile(`[0-9０-９]`)

func Change(number int) (data Number) {
	s2 := fmt.Sprintf("「%v」2進数:%b", number, number)
	s16 := fmt.Sprintf("「%v」16進数:%x", number, number)
	data = Number{
		S2:  s2,
		S16: s16,
	}
	return data
}

// Index 2進数・18進数変換フォーム
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// htmlファイルの読み込み
		pkg.GenerateHTML(w, nil, "change/index")

	} else if r.Method == "POST" {

		// 処理時間の計測 start
		start := time.Now()

		var data Data
		value := r.FormValue("number")
		// 正規表現を用いて判定
		if len(numReg.FindAllString(value, -1)) != 0 {
			// 全角を半角に変換
			num := width.Narrow.String(value)
			// 数値に変換
			number, err := strconv.Atoi(num)
			if err != nil {
				// 数字以外が含まれる時
				err := Err{
					Number: "数字以外が入力されています",
				}
				// Numberのエラーメッセージをdataに格納
				data.Error = err
				pkg.GenerateHTML(w, data, "change/index")
			} else {

				workers := 5
				ch := make(chan int)
				var numbers []Number
				var mutex = &sync.Mutex{}
				var wg sync.WaitGroup
				// 数値+19の20個まで5つ並行で処理
				for i := 0; i < workers; i++ {
					wg.Add(1)
					go func() {
						for num := range ch {
							number := Change(num)
							// 処理をロックする
							mutex.Lock()
							numbers = append(numbers, number)
							mutex.Unlock()
						}
						wg.Done()
					}()
				}

				for i := 0; i < 20; i++ {
					ch <- number + i
				}

				// この時点でゴルーチンに対するチャネルの送信は終了しているため閉じる
				close(ch)

				wg.Wait()
				// 処理時間の計測 end
				end := time.Now()
				process := end.Sub(start).Seconds()
				timeout := 0.01 //タイムアウト
				if process < timeout {
					data.Numbers = numbers
					pkg.GenerateHTML(w, data, "change/index")
				} else {

					// キャンセル可能なコンテキストを作る
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel() //メモリリークに繋がるため、必ず呼ぶ

					// ctxを親にした、タイムアウトするコンテキストを作る
					ctxChild, cancelChild := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
					defer cancelChild() //メモリリークに繋がるため、必ず呼ぶ

					// どちらかのコンテキストが完了するまで待つ
					select {
					case <-ctxChild.Done():
						// 20個処理されていない
						if len(numbers) != 20 {
							err := Err{
								Number: "正しく処理されませんでした",
							}
							data.Error = err
							// 20個処理される
						} else {
							log.Printf("処理完了に%v秒かかっています", timeout)
							data.Numbers = numbers
						}
						pkg.GenerateHTML(w, data, "change/index")
					case <-ctx.Done(): //cancel()を任意で呼び出す 処理結果を画面出力しない時使用
						if len(numbers) != 20 {
							log.Printf("処理結果は%v個で20個を満たしません\n", len(numbers))
						} else {
							log.Printf("処理結果は%v個で20個を満たしています\n", len(numbers))
						}
						// キャッシュをクリアにするためリダイレクト
						url := config.Config.URL + fmt.Sprint(config.FlagPort) + r.URL.Path
						http.Redirect(w, r, url, http.StatusSeeOther) //キャッシュを残したくないので、303指定
					}
				}
			}
		} else {
			// 数値でなければerrorを返す エラーメッセージ作成
			err := Err{
				Number: "数字を入力して下さい",
			}
			// Numberのエラーメッセージをdataに格納
			data.Error = err
			pkg.GenerateHTML(w, data, "change/index")
		}

	}
}
