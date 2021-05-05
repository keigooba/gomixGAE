package cli

import (
	"flag"
	"fmt"
	"gomix/config"
)

var ManageMemo bool

// Gitリポジトリのバージョン start.shで最新のバージョンに更新される
var Version = "1.0.0"

func CmdFlag() {

	// サードパーティーkingpinでの実装の場合
	// port := kingpin.Flag("port", "ポート設定が可能").Default("8888").Short('p').Int()
	// kingpin.Parse()

	// -v -versionが指定された場合にshowVerionが真になるよう定義
	flag.BoolVar(&ManageMemo, "memo", false, "メモの管理")
	flag.BoolVar(&ManageMemo, "m", false, "メモの管理(short)")

	// ポート設定のオプション
	// envPort, _ := strconv.Atoi(os.Getenv("PORT")) //環境変数でも指定できる
	flag.IntVar(&config.FlagPort, "port", config.Config.Port, "ポート設定が可能")
	flag.IntVar(&config.FlagPort, "p", config.Config.Port, "ポート設定が可能(short)")

	// Gitリポジトリのバージョン確認
	var showVersion bool

	// -v -versionが指定された場合にshowVerionが真になるよう定義
	flag.BoolVar(&showVersion, "version", false, "バージョン確認")
	flag.BoolVar(&showVersion, "v", false, "バージョン確認(short)")
	flag.Parse() //引数からオプションをパースする

	//ポート確認
	fmt.Println("port", config.FlagPort)
	if showVersion {
		// バージョン番号を表示する
		fmt.Println("version", Version)
	}
}
