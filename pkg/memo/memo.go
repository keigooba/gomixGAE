package memo

import (
	"bufio"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gomix/config"
	"gomix/pkg"

	"github.com/dustin/go-humanize"
	"github.com/mattn/go-isatty"
)

type Data struct {
	Json []string
	Txt  []string
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		var data Data
		// 作成済みファイルを読み込む
		slice := []string{"json", "txt"}
		for _, extension := range slice {
			if files, err := ioutil.ReadDir("doc/memo/data/" + extension); err == nil {
				for _, file := range files {
					if extension == "json" {
						data.Json = append(data.Json, file.Name())
					} else if extension == "txt" {
						data.Txt = append(data.Txt, file.Name())
					}
				}
			}
		}

		// ファイルの読み込み
		pkg.GenerateHTML(w, data, "memo/index")

	} else if r.Method == "POST" {

		memo := r.FormValue("memo")
		var extension string
		if string(memo[0]) == "[" {
			extension = "json"
		} else {
			extension = "txt"
		}

		// "/"でパス文字列を結合しない 物理パスを操作する場合filepathを使う
		dir := filepath.Join(pkg.Getpath(), "doc", "memo", "data", extension)
		// dirのディレクトリを作成する MkdirAll 必要な親ディレクトリ全てを作成する
		err := os.MkdirAll(dir, 0777) // os.ModePerm Unixパーミッションビット、0o777
		if err != nil {
			log.Println(err)
		}

		// 乱数シード生成
		var rTime int64
		if err := binary.Read(crand.Reader, binary.LittleEndian, &rTime); err != nil {
			rTime = time.Now().Unix()
		}
		rand.Seed(rTime)

		// 10進数の文字列に変換
		stringTime := strconv.FormatInt(rTime, 10)

		// 以降はパスを物理パスとして扱うのでfilepathパッケージを使う filepath.Base ファイル名を得る
		docPath := "/doc/memo/data/" + "memo_" + stringTime + "." + extension
		name := filepath.Join(pkg.Getpath(), "doc", "memo", "data", extension, filepath.Base(docPath))

		// メモに記載された文字のバイト数をログに出力
		log.Println(humanize.Bytes(uint64(len(memo))))

		defaultUmask := syscall.Umask(0) //umask値を変更
		err = ioutil.WriteFile(name, []byte(memo), 0777) // os.ModePerm Unixパーミッションビット、0o777
		syscall.Umask(defaultUmask)
		if err != nil {
			log.Println(err)
		}
		url := config.Config.URL + fmt.Sprint(config.FlagPort) + r.URL.Path
		http.Redirect(w, r, url, http.StatusSeeOther) //キャッシュを残したくないので、303指定
	}
}

func Open(w http.ResponseWriter, r *http.Request) {
	// httpリクエストは論理パスなのでpathを使う
	var extension string
	if strings.Contains(r.URL.Path, "json") {
		extension = "json"
	} else if strings.Contains(r.URL.Path, "txt") {
		extension = "txt"
	}

	if ok, err := path.Match("/data/"+extension+"/memo_*."+extension, r.URL.Path); err != nil || !ok {
		http.NotFound(w, r)
		return
	}

	// 指定したファイルを開く
	name := filepath.Join(pkg.Getpath(), "doc", "memo", "data", extension, filepath.Base(r.URL.Path))

	f, err := os.Open(name)
	if err != nil {
		log.Println(err)
	}
	defer pkg.Close(f)

	// ハッシュ作成
	hash := sha256.New()

	// ResponseWriterとハッシュ値に書き出す
	Mw := io.MultiWriter(w, hash)

	var b *bufio.Writer
	// var written int64
	if isatty.IsTerminal(f.Fd()) { //.Fd()で端末か判定
		// 出力先が端末
		_, err = io.Copy(Mw, f) //Writerにファイルを書き出す
	} else {
		//出力先がファイルやパイプ
		b = bufio.NewWriter(Mw) //バッファリングする
		_, err = io.Copy(b, f)  //Writerにファイルを書き出す
	}

	if err != nil {
		log.Println("ファイルの書き出しに失敗しました。")
	}

	if b != nil {
		// バッファリングされていたときFlush()する
		err := b.Flush()
		if err != nil {
			log.Println(err)
		}
	}

	//  書き込んだバイト数・ファイル名・ハッシュ値を取得する
	// fmt.Printf("Wrote %d, %s, %x", written, f.Name(), hash.Sum(nil))
}
