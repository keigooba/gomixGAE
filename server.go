package main

//go:generate statik -src=doc

import (
	"fmt"
	"gomix/config"
	"gomix/pkg"
	"gomix/pkg/change"
	"gomix/pkg/memo"
	"gomix/pkg/reflect"
	"net/http"
	"os"

	stats_api "github.com/fukata/golang-stats-api-handler"
	"github.com/markbates/pkger"
)

func StartMainServer() error {

	// doc以下のファイル読み込み
	dir := pkger.Dir(config.Config.Static) //バイナリファイルに静的ファイルを埋め込める
	files := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "config/info.json") //ファイルにアクセス
	})

	http.HandleFunc("/", pkg.Index)
	http.HandleFunc("/change", change.Index)
	http.HandleFunc("/memo", memo.Index)
	http.HandleFunc("/read", memo.ReadJson)
	http.HandleFunc("/data/", memo.Open)
	http.HandleFunc("/reflect", reflect.Index)
	http.HandleFunc("/stats", stats_api.Handler)
	port := os.Getenv("PORT")
	if port == "" {
		return http.ListenAndServe(":"+fmt.Sprint(config.FlagPort), nil)
	}
	return http.ListenAndServe(":"+port, nil)
}
