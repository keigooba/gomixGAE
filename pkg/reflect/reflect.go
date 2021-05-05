package reflect

import (
	"fmt"
	"gomix/pkg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	chs := make(chan []byte)
	files, err := ioutil.ReadDir("doc/memo/data/txt")
	if err != nil {
		log.Println(err)
	}
	var wg sync.WaitGroup
	for _, file := range files {
		name := filepath.Join(pkg.Getpath(), "doc", "memo", "data", "txt", file.Name())
		f, err := os.Open(name)
		if err != nil {
			return
		}
		defer pkg.Close(f)
		// ファイルから読み込んだバイト列を逐次chに送る
		wg.Add(1)
		go func() {
			buf := make([]byte, 4096)
			n, err := f.Read(buf)
			if err != nil {
				return
			}
			wg.Done()
			chs <- buf[:n]
		}()
	}
	wg.Wait()

	vec := []byte("\n")
	// 順次チャンネルからバイト列を読み込みそれを表示する
	go func() {
		for in := range chs {
			rv := reflect.ValueOf(in)
			i := rv.Kind()
			if fmt.Sprint(i) == "slice" {
				_, err := io.WriteString(w, string(in))
				if err != nil {
					log.Println(err)
				}
				_, err = io.WriteString(w, string(vec))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

}
