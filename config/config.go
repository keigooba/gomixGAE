// +build darwin,amd64 windows linux,!android
// +build go1.1

package config

import (
	"encoding/json"
	"fmt"
	"gomix/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/markbates/pkger"
	"google.golang.org/appengine"
)

type ConfigList struct {
	Port    int    `json:"port"`
	LogFile string `json:"log_file"`
	Static  string `json:"static"`
	URL     string `json:"up_url"` //本番時up_urlに変更
}

// Config Configの定義
var Config ConfigList

var Db *gorm.DB

// ポート変更のためここで定義
var FlagPort int

func init() {
	// Configの設定の読み込み
	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 現在の日付
	nowDate := time.Now().Format("200601")

	// ログファイルの設定
	if !appengine.IsAppEngine() { //GAEなら実行できないため、実行しない

		output := &utils.Output{OutStream: os.Stdout, ErrStream: os.Stderr}
		if Pwd != "" {
			// 絶対パスを指定
			err = utils.LoggingSettings(Pwd+"/"+Config.LogFile+nowDate+".log", output)
		} else {
			err = utils.LoggingSettings(Config.LogFile+nowDate+".log", output)
		}
		if err != nil {
			log.Println(err)
		}
		// コマンドの実行
		err = utils.Command()
		if err != nil {
			log.Println(err) //本番のroot権限ではコマンド実行できないため、出力のみ
		}
	}

	// DB接続
	Db = utils.GormConnect()

	log.Println("DBクリア")
}

// LoadConfig Configの設定
func LoadConfig() error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	} else if utils.Pwd != "" {
		cwd = utils.Pwd
	}

	fname := filepath.Join(cwd, "config", "config.json")
	f, err := pkger.Open(fname)
	if err != nil {
		return err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	//Configにconfig.jsonを読み込む
	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		return err
	}

	// 環境変数の値の判定
	format := "Port: %d\nLogFile: %s\nStatic: %s\nURL: %s\n"
	_, err = fmt.Printf(format, Config.Port, Config.LogFile, Config.Static, Config.URL)
	if err != nil {
		return err
	}
	return nil //自明であればnilにする
}
