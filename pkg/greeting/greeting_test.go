package greeting

import (
	"testing"
)

func TestGetGreeting(t *testing.T) {
	expected := "こんにちは、太郎さん！"
	actual := GetGreeting("太郎")
	if actual != expected {
		t.Errorf("期待値: %s, 実際の値: %s", expected, actual)
	}
}
