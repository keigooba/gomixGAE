package change

import (
	"testing"
)

func TestChange(t *testing.T) {

	number := 255

	// テストケースの検証
	data := Change(number)
	if data.S2 != "「255」2進数:11111111" {
		t.Errorf("2進数は%sで失敗しています。", data.S2)
	}
	if data.S16 != "「255」16進数:ff" {
		t.Errorf("16進数は%sで失敗しています。", data.S16)
	}

	// t.Log, t.Logf でログを出すと `go test -v` と実行したときのみ表示される
	t.Logf("result is %s, %s", data.S2, data.S16)
}
