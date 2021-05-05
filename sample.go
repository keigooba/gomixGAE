package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"sync"
// 	"syscall"
// 	"time"

// 	"github.com/mattn/go-shellwords"
// )

// // func tr(src io.Reader, dst io.Writer, errDst io.Writer) error {
// // 	cmd := exec.Command("tr", "a-z", "A-Z")
// // 	// 実行するコマンド tr a-z A-Z
// // 	stdin, _ := cmd.StdinPipe()
// // 	stdout, _ := cmd.StdoutPipe()
// // 	stderr, _ := cmd.StderrPipe()
// // 	err := cmd.Start() //コマンドの実行を開始する
// // 	if err != nil {
// // 		return err
// // 	}
// // 	var wg sync.WaitGroup
// // 	wg.Add(3)
// // 	go func() {
// // 		// コマンドの標準入力にsrcからコピーする
// // 		_, err := io.Copy(stdin, src)
// // 		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE {
// // 			// ignore EPIPE
// // 		} else if err != nil {
// // 			log.Println("failed to write to STDIN", err)
// // 		}
// // 		stdin.Close()
// // 		wg.Done()
// // 	}()
// // 	go func() {
// // 		io.Copy(dst, stdout)
// // 		stdout.Close()
// // 		wg.Done()
// // 	}()
// // 	go func() {
// // 		io.Copy(errDst, stderr)
// // 		fmt.Println("stderr")
// // 		stderr.Close()
// // 		wg.Done()
// // 	}()
// // 	wg.Wait()
// // 	//標準入出力のI/Oを行うgoroutineが全て終わるまで待つ
// // 	return cmd.Wait()
// // 	// コマンドの終了を待つ
// // }

// // パーミッション確認
// func execLs() {
// 	out, _ := exec.Command("ls", "-l", "config").Output()
// 	fmt.Println(string(out))
// }

// // go-shellwords 一つずつ出力を確認しながら実行したいとき使用
// func shell() {
// 	aregs, err := shellwords.Parse("ls -l config")
// 	// argsは["ls", "-l", "config"]となる
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	out, err := exec.Command(aregs[0], aregs[1:]...).Output()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	fmt.Println(string(out))
// }

// func getHTTP(url string, dst io.Writer) error {
// 	// 10秒でタイムアウトするContextを作る
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client := &http.Client{}
// 	// client := &http.Client{
// 	// 	// 10秒でタイムアウトする
// 	// 	Timeout: 10 * time.Second,
// 	// }
// 	req, _ := http.NewRequest("GET", url, nil)
// 	// contextを与えたリクエストを使って実行
// 	resp, err := client.Do(req.WithContext(ctx)) // 呼び出す度に10秒待つ
// 	// resp, err := client.Do(req)
// 	if err != nil {
// 		// レスポンスヘッダーの取得までにエラー
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	_, err = io.Copy(dst, resp.Body)
// 	// ボディ取得完了までにエラー
// 	return err
// }

// // ゴルーチンの処理
// var wg sync.WaitGroup

// func goroutine() {
// 	queue := make(chan string)

// 	for i := 0; i < 2; i++ {
// 		wg.Add(1)
// 		go fetchURL(queue)
// 	}

// 	queue <- "https://www.example.com"
// 	queue <- "https://www.example.net"
// 	queue <- "https://www.example.net/foo"
// 	queue <- "https://www.example.net/bar"

// 	close(queue)
// 	wg.Wait()
// }

// func fetchURL(queue chan string) {
// 	for url := range queue {
// 		fmt.Println("fetching", url)
// 	}
// 	fmt.Println("worker exit")
// 	wg.Done()
// 	return
// }

// // シグナルによるゴルーチンの脱出
// func signalCalls() {
// 	defer log.Println("受信を終了します")
// 	// シグナルを決める
// 	trapSignals := []os.Signal{
// 		syscall.SIGHUP,
// 		syscall.SIGINT, //Ctrl+Cのシグナル
// 		syscall.SIGTERM,
// 		syscall.SIGQUIT,
// 	}
// 	// 受信するチャンネル
// 	sigCh := make(chan os.Signal, 1)
// 	// 受信する
// 	signal.Notify(sigCh, trapSignals...)

// 	test := make(chan int, 1)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	// 別goroutineでシグナルを待ち受ける
// 	go func() {
// 		sig := <-sigCh
// 		fmt.Println("シグナルが入力されました", sig)
// 		//終了させるためキャンセルを実行
// 		cancel()
// 	}()
// 	go func() {
// 		test <- 100 //こういった形でゴルーチンをつくることで受信処理する。
// 	}()
// 	fmt.Println(test)
// 	Receive(ctx, test)
// 	log.Panic("プロセスが強制終了しました")
// }

// func Receive(ctx context.Context, test chan int) {
// 	defer log.Println("Receiveの受信を閉じます")
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return //無限ループを抜ける
// 		case v := <-test:
// 			fmt.Println(v)
// 		}
// 		//何らかの受信処理 プロセスが終了すれば受信できなくなる
// 	}
// }

// /*シグナルの用途 このサーバーのプロセスが強制終了したあとの挙動

// ・外部からの新規リクエストの受信停止
// ・その場でプロセスが終了するとdeferの処理完了を待たないため、panicで待つ
// ・bufio.Flush()してからファイルを閉じる

// 考えられる用途
// 別のサーバーと連携していてこのサーバーが停止した時、ログを残すようにして停止のレスポンスを返す。
// あるいは別のサーバーをhttpで10秒以内に受信できなくなった時、停止済みのレスポンスを返す。

// 利用法
// 主に終了通知（受信先の停止等)に使い、ゴルーチンで並列処理としてブロッキングして利用する
// */

// /*シグナル活用法
// ・外部からコマンドを受け付けて停止処理を行う
// ・特定の条件を満たしたら終了する
// */

// // 独自シグナルによるタイムアウト処理
// type MySignal struct {
// 	message string
// }

// func (s MySignal) String() string {
// 	return s.message
// }

// func (s MySignal) Signal() {}

// func mysignal() {
// 	log.Println("[info] Start")
// 	// シグナルを決める
// 	trapSignals := []os.Signal{
// 		syscall.SIGHUP,
// 		syscall.SIGINT, //Ctrl+Cのシグナル
// 		syscall.SIGTERM,
// 		syscall.SIGQUIT,
// 	}
// 	// 受信するチャンネル
// 	sigCh := make(chan os.Signal, 1)

// 	// 10秒後にsigChにMySignalの値を送信
// 	time.AfterFunc(10*time.Second, func() {
// 		sigCh <- MySignal{"timed out"}
// 	})

// 	signal.Notify(sigCh, trapSignals...)

// 	// 受信するまで待ち受ける
// 	sig := <-sigCh
// 	switch s := sig.(type) { //型アサーションで判別
// 	case syscall.Signal:
// 		// osからのシグナルの場合
// 		log.Panicf("[info] Got signal: %s(%d)", s, s)
// 	case MySignal:
// 		// アプリケーション独自のシグナルの場合
// 		log.Panicf("[info] %s", s) //.String()が入る
// 	}
// }

// func readFromFile(ch chan []byte, f *os.File) {
// 	defer close(ch)
// 	defer f.Close()

// 	buf := make([]byte, 4096)
// 	for {
// 		if n, err := f.Read(buf); err == nil {
// 			ch <- buf[:n]
// 		}
// 	}
// }

// func makeChannelForFiles(files []fs.FileInfo) ([]reflect.Value, error) {
// 	cs := make([]reflect.Value, len(files))

// 	for i, file := range files {
// 		// データ送信のチャネル
// 		ch := make(chan []byte)

// 		// ファイルをオープン
// 		f, err := os.Open(file.Name())
// 		if err != nil {
// 			return nil, err
// 		}
// 		go readFromFile(ch, f)

// 		cs[i] = reflect.ValueOf(ch)
// 	}
// 	return cs, nil
// }

// func makeSelectCases(cs ...reflect.Value) ([]reflect.SelectCase, error) {
// 	cases := make([]reflect.SelectCase, len(cs))
// 	for i, ch := range cs {
// 		if ch.Kind() != reflect.Chan {
// 			return nil, errors.New("チャンネルが必要です")
// 		}
// 		cases[i] = reflect.SelectCase{
// 			Chan: ch,
// 			Dir:  reflect.SelectRecv,
// 		}
// 	}
// 	return cases, nil
// }

// func Index(w http.ResponseWriter, r *http.Request) {
// 	files, err := ioutil.ReadDir("doc/memo/data/json")
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	makeChannelForFiles(files)
// }
