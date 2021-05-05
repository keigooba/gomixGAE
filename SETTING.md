# gomix

# 初期設定

1. Go module 導入 go mod init
2. git 導入
3. 各種メトリクス取得 API導入 golang-stats-api-handler
4. Makefile の作成
5. ログファイルの作成 logrusの導入
6. header,footer の切り分け
7. test ファイルの作成
8. config ファイルの作成
9. golang-lintの導入
10. 自動化 sh ファイルの作成・exec.commandで自動実行
11. 静的ファイルをバイナリファイルに埋め込む pkgerの導入
12. バイナリファイルにGitのバージョン埋め込み設定・最新リリースバージョンチェック
13. サーバー停止時のシグナルの作成
14. githubによるバイナリファイルリリースの設定
15. goxによるMac・Linux・Windowsに対応したそれぞれのバイナリファイルの作成
16. ポートの変更のコマンド作成

# 注意事項

1. runtime.GOOS or //+build を用いて windows でも同様の動作環境で動くようにすること
2. 各種メトリクス取得 API を利用して Munin や Zabbix 等のエージェント経由でメモリや GC の状況をモニタリングすること
3. ログレベルを分けて出力を見たい場合、logrusでログ出力する
4. リリース更新時、upVersion・info.jsonの更新を行うこと
5. golangci-lintでコードを検査し、可読性を担保する
