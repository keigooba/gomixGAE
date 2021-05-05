package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// サーバー停止のログ通知
func signalCall() {
	defer log.Println("プロセスを終了します")
	// シグナルを決める
	trapSignals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT, //Ctrl+Cのシグナル
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	// 受信するチャンネル
	sigCh := make(chan os.Signal, 1)
	// シグナルを受信する
	signal.Notify(sigCh, trapSignals...)

	for {
		sig := <-sigCh // signal.Notifyが受信されたら受け取る
		log.Println("シグナルを受信しました", sig)
		//終了させるためキャンセルを実行
		log.Panic("サーバーが強制終了しました")
	}
}
